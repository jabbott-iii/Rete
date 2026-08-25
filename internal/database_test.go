/*
Copyright 2026 Joseph Anthony Abbott III

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
	"path/filepath"
	"testing"
)

func newTestDB(t *testing.T) *Database {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rete_test.db")
	db, err := NewDatabase(path)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	return db
}

func TestNewDatabaseCreatesSchema(t *testing.T) {
	db := newTestDB(t)
	if db == nil {
		t.Fatal("expected non-nil Database")
	}
	if db.Conn() == nil {
		t.Fatal("expected non-nil gorm connection")
	}
}

func TestUpsertCategoryIdempotent(t *testing.T) {
	db := newTestDB(t)

	const name = "Reconnaissance & Discovery"

	cat1, err := db.UpsertCategory(name)
	if err != nil {
		t.Fatalf("first UpsertCategory: %v", err)
	}
	cat2, err := db.UpsertCategory(name)
	if err != nil {
		t.Fatalf("second UpsertCategory: %v", err)
	}
	if cat1.ID != cat2.ID {
		t.Errorf("expected same ID on second upsert; got %d and %d", cat1.ID, cat2.ID)
	}
}

func TestUpsertFeatureIdempotent(t *testing.T) {
	db := newTestDB(t)

	cat, err := db.UpsertCategory("Network Diagnostics")
	if err != nil {
		t.Fatalf("UpsertCategory: %v", err)
	}

	f := &Feature{
		CategoryID:       cat.ID,
		Name:             "Traceroute",
		CobraCommand:     "diag traceroute",
		ShortDescription: "Trace the path packets take to a destination",
	}

	if err := db.UpsertFeature(f); err != nil {
		t.Fatalf("first UpsertFeature: %v", err)
	}

	// second upsert must not error and must not create a duplicate
	if err := db.UpsertFeature(f); err != nil {
		t.Fatalf("second UpsertFeature: %v", err)
	}

	features, err := db.ListFeatures("")
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}
	count := 0
	for _, feat := range features {
		if feat.CobraCommand == "diag traceroute" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 traceroute feature, got %d", count)
	}
}

func TestListFeaturesFilterByCategory(t *testing.T) {
	db := newTestDB(t)

	cat1, _ := db.UpsertCategory("Reconnaissance & Discovery")
	cat2, _ := db.UpsertCategory("Payload Delivery")

	_ = db.UpsertFeature(&Feature{CategoryID: cat1.ID, Name: "DNS Lookup", CobraCommand: "recon dns", ShortDescription: "DNS"})
	_ = db.UpsertFeature(&Feature{CategoryID: cat2.ID, Name: "Bind Shell", CobraCommand: "payload bind-shell", ShortDescription: "Bind shell"})

	features, err := db.ListFeatures("Reconnaissance & Discovery")
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}
	if len(features) != 1 {
		t.Errorf("expected 1 feature, got %d", len(features))
	}
	if features[0].CobraCommand != "recon dns" {
		t.Errorf("unexpected command %q", features[0].CobraCommand)
	}
}

func TestJobCRUD(t *testing.T) {
	db := newTestDB(t)

	job := &ScanJob{
		Command: "recon port-scan",
		Target:  "10.0.0.1",
		Args:    "ports=1-1024 proto=tcp",
		Status:  JobStatusPending,
	}
	if err := db.CreateJob(job); err != nil {
		t.Fatalf("CreateJob: %v", err)
	}
	if job.ID == 0 {
		t.Fatal("expected non-zero job ID after create")
	}

	fetched, err := db.GetJobByID(job.ID)
	if err != nil {
		t.Fatalf("GetJobByID: %v", err)
	}
	if fetched.Command != job.Command {
		t.Errorf("command mismatch: got %q want %q", fetched.Command, job.Command)
	}

	fetched.Status = JobStatusDone
	if err := db.UpdateJob(fetched); err != nil {
		t.Fatalf("UpdateJob: %v", err)
	}

	updated, err := db.GetJobByID(job.ID)
	if err != nil {
		t.Fatalf("GetJobByID after update: %v", err)
	}
	if updated.Status != JobStatusDone {
		t.Errorf("expected status %q, got %q", JobStatusDone, updated.Status)
	}
}

func TestListJobsLimit(t *testing.T) {
	db := newTestDB(t)

	for i := 0; i < 5; i++ {
		_ = db.CreateJob(&ScanJob{Command: "diag latency", Target: "8.8.8.8", Status: JobStatusPending})
	}

	jobs, err := db.ListJobs(3)
	if err != nil {
		t.Fatalf("ListJobs: %v", err)
	}
	if len(jobs) != 3 {
		t.Errorf("expected 3 jobs, got %d", len(jobs))
	}
}

func TestAddResult(t *testing.T) {
	db := newTestDB(t)

	job := &ScanJob{Command: "sec sniff", Target: "eth0", Status: JobStatusRunning}
	_ = db.CreateJob(job)

	result := &ScanResult{JobID: job.ID, Output: "captured 10 packets", IsError: false}
	if err := db.AddResult(result); err != nil {
		t.Fatalf("AddResult: %v", err)
	}

	fetched, err := db.GetJobByID(job.ID)
	if err != nil {
		t.Fatalf("GetJobByID: %v", err)
	}
	if len(fetched.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(fetched.Results))
	}
	if fetched.Results[0].Output != result.Output {
		t.Errorf("output mismatch: got %q want %q", fetched.Results[0].Output, result.Output)
	}
}
