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
	"strings"
	"testing"
)

const sampleCSV = `category,feature,cobra_command,short_description
Reconnaissance & Discovery,Ping Sweep,recon ping-sweep,Discover live hosts by sending ICMP echo requests
Reconnaissance & Discovery,Port Scan,recon port-scan,Scan TCP/UDP ports on a target host
Network Diagnostics,Traceroute,diag traceroute,Trace the path packets take to a destination
Payload Delivery,Bind Shell,payload bind-shell,Set up a bind shell listener on a port
`

func TestSeedFeaturesFromReader(t *testing.T) {
	db := newTestDB(t)

	count, err := SeedFeaturesFromReader(db, strings.NewReader(sampleCSV))
	if err != nil {
		t.Fatalf("SeedFeaturesFromReader: %v", err)
	}
	if count != 4 {
		t.Errorf("expected 4 features seeded, got %d", count)
	}
}

func TestSeedFeaturesIdempotent(t *testing.T) {
	db := newTestDB(t)

	count1, err := SeedFeaturesFromReader(db, strings.NewReader(sampleCSV))
	if err != nil {
		t.Fatalf("first seed: %v", err)
	}

	count2, err := SeedFeaturesFromReader(db, strings.NewReader(sampleCSV))
	if err != nil {
		t.Fatalf("second seed: %v", err)
	}

	// Second seed should still return the same row count (idempotent upsert).
	if count1 != count2 {
		t.Errorf("expected idempotent seed counts (%d == %d)", count1, count2)
	}

	// But DB should still have only 4 features (no duplicates).
	features, err := db.ListFeatures("")
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}
	if len(features) != 4 {
		t.Errorf("expected 4 features after double seed, got %d", len(features))
	}
}

func TestSeedFeaturesCategories(t *testing.T) {
	db := newTestDB(t)

	_, err := SeedFeaturesFromReader(db, strings.NewReader(sampleCSV))
	if err != nil {
		t.Fatalf("SeedFeaturesFromReader: %v", err)
	}

	cats, err := db.ListCategories()
	if err != nil {
		t.Fatalf("ListCategories: %v", err)
	}
	if len(cats) != 3 {
		t.Errorf("expected 3 categories, got %d", len(cats))
	}
}

func TestSeedFeaturesMissingColumn(t *testing.T) {
	db := newTestDB(t)

	bad := "feature,cobra_command\nPing Sweep,recon ping-sweep\n"
	_, err := SeedFeaturesFromReader(db, strings.NewReader(bad))
	if err == nil {
		t.Fatal("expected error for CSV missing 'category' column")
	}
}

func TestSeedFeaturesFromCSVFile(t *testing.T) {
	db := newTestDB(t)

	// Use the real table.csv in the repo root (two directories up from internal/).
	count, err := SeedFeaturesFromCSV(db, "../table.csv")
	if err != nil {
		t.Fatalf("SeedFeaturesFromCSV: %v", err)
	}
	if count == 0 {
		t.Error("expected at least one feature seeded from table.csv")
	}

	features, err := db.ListFeatures("")
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}
	if len(features) != count {
		t.Errorf("DB feature count %d does not match seed count %d", len(features), count)
	}
}
