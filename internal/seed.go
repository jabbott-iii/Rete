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

import "fmt"

// defaultCategory names the FeatureCategory rows seeded on startup.
type defaultCategory struct {
	Name        string
	Description string
}

// defaultFeature names the Feature rows seeded on startup.
type defaultFeature struct {
	Category    string
	Key         string
	Name        string
	Description string
}

var defaultCategories = []defaultCategory{
	{Name: "System Resources", Description: "Disk, memory, and CPU headroom on the local host."},
	{Name: "Container Runtime", Description: "Docker daemon availability and container health."},
	{Name: "Orchestration", Description: "Kubernetes cluster reachability and status."},
	{Name: "Services", Description: "Uptime of managed services or the host process."},
	{Name: "Configuration", Description: "Common environment misconfigurations."},
}

var defaultFeatures = []defaultFeature{
	{Category: "System Resources", Key: keyDiskSpace, Name: "Disk Space", Description: "Verifies free disk space on the configured mount point."},
	{Category: "System Resources", Key: keyMemory, Name: "Memory", Description: "Verifies available memory and swap usage."},
	{Category: "System Resources", Key: keyCPULoad, Name: "CPU Load", Description: "Verifies the system load average relative to available CPUs."},
	{Category: "Container Runtime", Key: keyDocker, Name: "Docker Status", Description: "Verifies the Docker daemon is reachable."},
	{Category: "Orchestration", Key: keyKubernetes, Name: "Kubernetes Status", Description: "Verifies the configured Kubernetes cluster is reachable."},
	{Category: "Services", Key: keyServiceUptime, Name: "Service Uptime", Description: "Verifies a systemd service (or the host) is up and running."},
	{Category: "Configuration", Key: keyMisconfig, Name: "Misconfiguration", Description: "Scans for common environment misconfigurations."},
}

// EnsureDefaultFeatures seeds the built-in feature catalog if it is not already present.
func EnsureDefaultFeatures(db *Database) error {
	conn := db.Conn()

	categoryIDs := make(map[string]uint, len(defaultCategories))
	for _, cat := range defaultCategories {
		record := FeatureCategory{Name: cat.Name, Description: cat.Description}
		if err := conn.Where(FeatureCategory{Name: cat.Name}).FirstOrCreate(&record).Error; err != nil {
			return fmt.Errorf("seed category %q: %w", cat.Name, err)
		}
		categoryIDs[cat.Name] = record.ID
	}

	for _, feat := range defaultFeatures {
		categoryID, ok := categoryIDs[feat.Category]
		if !ok {
			return fmt.Errorf("seed feature %q: unknown category %q", feat.Key, feat.Category)
		}

		record := Feature{
			CategoryID:  categoryID,
			Key:         feat.Key,
			Name:        feat.Name,
			Description: feat.Description,
		}
		if err := conn.Where(Feature{Key: feat.Key}).FirstOrCreate(&record).Error; err != nil {
			return fmt.Errorf("seed feature %q: %w", feat.Key, err)
		}
	}

	return nil
}
