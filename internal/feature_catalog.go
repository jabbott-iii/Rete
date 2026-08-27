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

// CatalogEntry represents a single feature in the compiled-in default catalog.
type CatalogEntry struct {
	Category         string
	Feature          string
	CobraCommand     string
	ShortDescription string
}

// DefaultFeatureCatalog is the canonical feature catalog compiled into the binary.
// It mirrors the data previously distributed as table.csv.
var DefaultFeatureCatalog = []CatalogEntry{
	// Reconnaissance & Discovery
	{"Reconnaissance & Discovery", "Ping Sweep", "recon ping-sweep", "Discover live hosts by sending ICMP echo requests"},
	{"Reconnaissance & Discovery", "Port Scan", "recon port-scan", "Scan TCP/UDP ports on a target host"},
	{"Reconnaissance & Discovery", "OS Fingerprinting", "recon os-finger", "Identify the remote host operating system"},
	{"Reconnaissance & Discovery", "Service Version Enumeration", "recon service-enum", "Detect running services and their versions"},
	{"Reconnaissance & Discovery", "DNS Lookup", "recon dns", "Perform forward and reverse DNS resolution"},
	{"Reconnaissance & Discovery", "WHOIS Lookup", "recon whois", "Query WHOIS data for a domain or IP address"},
	{"Reconnaissance & Discovery", "Subdomain Enumeration", "recon subdomain", "Enumerate subdomains via DNS brute-force"},
	// Network Diagnostics
	{"Network Diagnostics", "Traceroute", "diag traceroute", "Trace the path packets take to a destination"},
	{"Network Diagnostics", "Bandwidth Test", "diag bandwidth", "Measure available network bandwidth to a host"},
	{"Network Diagnostics", "Latency Monitor", "diag latency", "Monitor round-trip latency over time"},
	{"Network Diagnostics", "Packet Loss Check", "diag packet-loss", "Measure packet loss percentage to a remote host"},
	{"Network Diagnostics", "MTU Discovery", "diag mtu", "Discover the maximum transmission unit on a path"},
	{"Network Diagnostics", "ARP Scan", "diag arp-scan", "Discover hosts on the local network via ARP"},
	// Security & Penetration Testing
	{"Security & Penetration Testing", "Packet Sniffer", "sec sniff", "Capture and inspect packets on a network interface"},
	{"Security & Penetration Testing", "Packet Forger", "sec forge", "Craft and inject custom network packets"},
	{"Security & Penetration Testing", "ARP Spoofer", "sec arp-spoof", "Poison ARP caches to perform MITM interception"},
	{"Security & Penetration Testing", "Vulnerability Scanner", "sec vuln-scan", "Scan for common vulnerabilities on a target"},
	{"Security & Penetration Testing", "Brute Force", "sec brute", "Attempt credential brute-force against a service"},
	{"Security & Penetration Testing", "SSL/TLS Audit", "sec tls-audit", "Audit SSL/TLS configuration on a remote host"},
	// Payload Delivery
	{"Payload Delivery", "Reverse Shell Helper", "payload rev-shell", "Generate reverse shell payload templates"},
	{"Payload Delivery", "File Transfer", "payload file-xfer", "Transfer files over raw TCP or HTTP channels"},
	{"Payload Delivery", "Bind Shell", "payload bind-shell", "Set up a bind shell listener on a port"},
	{"Payload Delivery", "Encoded Payload", "payload encode", "Encode a payload to evade simple pattern filters"},
}

// SeedDefaultFeatures upserts every entry from DefaultFeatureCatalog into the
// database. It is safe to call multiple times; existing records are not
// duplicated.
func SeedDefaultFeatures(db *Database) error {
	catCache := make(map[string]*FeatureCategory)

	for _, entry := range DefaultFeatureCatalog {
		cat, ok := catCache[entry.Category]
		if !ok {
			var err error
			cat, err = db.UpsertCategory(entry.Category)
			if err != nil {
				return err
			}
			catCache[entry.Category] = cat
		}

		f := &Feature{
			CategoryID:       cat.ID,
			Name:             entry.Feature,
			CobraCommand:     entry.CobraCommand,
			ShortDescription: entry.ShortDescription,
		}
		if err := db.UpsertFeature(f); err != nil {
			return err
		}
	}
	return nil
}

// EnsureDefaultFeatures seeds the default catalog if no features exist yet.
// This is intended to be called once at startup.
func EnsureDefaultFeatures(db *Database) error {
	features, err := db.ListFeatures("")
	if err != nil {
		return err
	}
	if len(features) > 0 {
		return nil
	}
	return SeedDefaultFeatures(db)
}
