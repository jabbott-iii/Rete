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
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// SeedFeaturesFromCSV reads the canonical feature catalog CSV and upserts all
// categories and features into the database. It returns the count of features
// processed without error.
//
// Expected CSV columns (header row required):
//
//	category, feature, cobra_command, short_description
func SeedFeaturesFromCSV(db *Database, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	return SeedFeaturesFromReader(db, f)
}

// SeedFeaturesFromReader parses the CSV from r and upserts features, allowing
// callers to pass any io.Reader (useful in tests).
func SeedFeaturesFromReader(db *Database, r io.Reader) (int, error) {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true

	header, err := cr.Read()
	if err != nil {
		return 0, fmt.Errorf("read CSV header: %w", err)
	}

	// Build column index map for resilience against column reordering.
	colIdx := make(map[string]int, len(header))
	for i, h := range header {
		colIdx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	required := []string{"category", "feature", "cobra_command", "short_description"}
	for _, col := range required {
		if _, ok := colIdx[col]; !ok {
			return 0, fmt.Errorf("CSV missing required column %q", col)
		}
	}

	count := 0
	catCache := make(map[string]*FeatureCategory)

	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, fmt.Errorf("read CSV row: %w", err)
		}

		catName := strings.TrimSpace(row[colIdx["category"]])
		featName := strings.TrimSpace(row[colIdx["feature"]])
		cobraCmd := strings.TrimSpace(row[colIdx["cobra_command"]])
		shortDesc := strings.TrimSpace(row[colIdx["short_description"]])

		if catName == "" || featName == "" || cobraCmd == "" {
			continue // skip blank rows
		}

		// Upsert category (cached to avoid redundant DB calls).
		cat, ok := catCache[catName]
		if !ok {
			cat, err = db.UpsertCategory(catName)
			if err != nil {
				return count, fmt.Errorf("upsert category %q: %w", catName, err)
			}
			catCache[catName] = cat
		}

		feature := &Feature{
			CategoryID:       cat.ID,
			Name:             featName,
			CobraCommand:     cobraCmd,
			ShortDescription: shortDesc,
		}
		if err := db.UpsertFeature(feature); err != nil {
			return count, fmt.Errorf("upsert feature %q: %w", featName, err)
		}
		count++
	}

	return count, nil
}
