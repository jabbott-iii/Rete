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
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// subcommandSet returns the set of Use values for direct subcommands of cmd.
func subcommandSet(cmd *cobra.Command) map[string]bool {
	m := make(map[string]bool, len(cmd.Commands()))
	for _, c := range cmd.Commands() {
		m[c.Use] = true
	}
	return m
}

func TestRootCmdHasExpectedSubcommands(t *testing.T) {
	db := newTestDB(t)
	root := NewRootCmd(db)

	expected := []string{"recon", "diag", "sec", "payload", "features", "jobs"}
	got := subcommandSet(root)
	for _, want := range expected {
		if !got[want] {
			t.Errorf("root command missing subcommand %q", want)
		}
	}
}

func TestReconCmdHasAllSubcommands(t *testing.T) {
	db := newTestDB(t)
	recon := NewReconCmd(db)

	expected := []string{"ping-sweep", "port-scan", "os-finger", "service-enum", "dns", "whois", "subdomain"}
	got := subcommandSet(recon)
	for _, want := range expected {
		if !got[want] {
			t.Errorf("recon missing subcommand %q", want)
		}
	}
}

func TestDiagCmdHasAllSubcommands(t *testing.T) {
	db := newTestDB(t)
	diag := NewDiagCmd(db)

	expected := []string{"traceroute", "bandwidth", "latency", "packet-loss", "mtu", "arp-scan"}
	got := subcommandSet(diag)
	for _, want := range expected {
		if !got[want] {
			t.Errorf("diag missing subcommand %q", want)
		}
	}
}

func TestSecCmdHasAllSubcommands(t *testing.T) {
	db := newTestDB(t)
	sec := NewSecCmd(db)

	expected := []string{"sniff", "forge", "arp-spoof", "vuln-scan", "brute", "tls-audit"}
	got := subcommandSet(sec)
	for _, want := range expected {
		if !got[want] {
			t.Errorf("sec missing subcommand %q", want)
		}
	}
}

func TestPayloadCmdHasAllSubcommands(t *testing.T) {
	db := newTestDB(t)
	payload := NewPayloadCmd(db)

	expected := []string{"rev-shell", "file-xfer", "bind-shell", "encode"}
	got := subcommandSet(payload)
	for _, want := range expected {
		if !got[want] {
			t.Errorf("payload missing subcommand %q", want)
		}
	}
}

func TestReconPortScanRecordsJob(t *testing.T) {
	db := newTestDB(t)
	root := NewRootCmd(db)

	var out bytes.Buffer
	root.SetOut(&out)

	root.SetArgs([]string{"recon", "port-scan", "--target", "10.0.0.1", "--ports", "22,80"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	jobs, err := db.ListJobs(10)
	if err != nil {
		t.Fatalf("ListJobs: %v", err)
	}
	if len(jobs) == 0 {
		t.Fatal("expected at least one job to be recorded")
	}
	if jobs[0].Command != "recon port-scan" {
		t.Errorf("unexpected command %q", jobs[0].Command)
	}
	if jobs[0].Target != "10.0.0.1" {
		t.Errorf("unexpected target %q", jobs[0].Target)
	}
}

func TestFeaturesListNoFeatures(t *testing.T) {
	db := newTestDB(t)
	root := NewRootCmd(db)

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"features", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if !strings.Contains(out.String(), "No features found") {
		t.Errorf("expected 'No features found' message, got: %q", out.String())
	}
}

func TestJobsListNoJobs(t *testing.T) {
	db := newTestDB(t)
	root := NewRootCmd(db)

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"jobs", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if !strings.Contains(out.String(), "No jobs") {
		t.Errorf("expected 'No jobs' message, got: %q", out.String())
	}
}
