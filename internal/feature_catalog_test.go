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
	"testing"
)

func TestSeedDefaultFeatures(t *testing.T) {
	db := newTestDB(t)

	if err := SeedDefaultFeatures(db); err != nil {
		t.Fatalf("SeedDefaultFeatures: %v", err)
	}

	features, err := db.ListFeatures("")
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}
	if len(features) != len(DefaultFeatureCatalog) {
		t.Errorf("expected %d features, got %d", len(DefaultFeatureCatalog), len(features))
	}
}

func TestSeedDefaultFeaturesIdempotent(t *testing.T) {
	db := newTestDB(t)

	if err := SeedDefaultFeatures(db); err != nil {
		t.Fatalf("first SeedDefaultFeatures: %v", err)
	}
	if err := SeedDefaultFeatures(db); err != nil {
		t.Fatalf("second SeedDefaultFeatures: %v", err)
	}

	features, err := db.ListFeatures("")
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}
	if len(features) != len(DefaultFeatureCatalog) {
		t.Errorf("expected %d features after double seed, got %d", len(DefaultFeatureCatalog), len(features))
	}
}

func TestEnsureDefaultFeaturesSeedsWhenEmpty(t *testing.T) {
	db := newTestDB(t)

	if err := EnsureDefaultFeatures(db); err != nil {
		t.Fatalf("EnsureDefaultFeatures: %v", err)
	}

	features, err := db.ListFeatures("")
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}
	if len(features) != len(DefaultFeatureCatalog) {
		t.Errorf("expected %d features after EnsureDefaultFeatures, got %d", len(DefaultFeatureCatalog), len(features))
	}
}

func TestEnsureDefaultFeaturesSkipsWhenPopulated(t *testing.T) {
	db := newTestDB(t)

	// Seed once so catalog is not empty.
	if err := SeedDefaultFeatures(db); err != nil {
		t.Fatalf("SeedDefaultFeatures: %v", err)
	}

	// EnsureDefaultFeatures should be a no-op; no duplicates should appear.
	if err := EnsureDefaultFeatures(db); err != nil {
		t.Fatalf("EnsureDefaultFeatures: %v", err)
	}

	features, err := db.ListFeatures("")
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}
	if len(features) != len(DefaultFeatureCatalog) {
		t.Errorf("expected %d features, got %d", len(DefaultFeatureCatalog), len(features))
	}
}

func TestDefaultCatalogCategories(t *testing.T) {
	db := newTestDB(t)

	if err := SeedDefaultFeatures(db); err != nil {
		t.Fatalf("SeedDefaultFeatures: %v", err)
	}

	cats, err := db.ListCategories()
	if err != nil {
		t.Fatalf("ListCategories: %v", err)
	}
	// Verify at least 4 categories are present.
	if len(cats) < 4 {
		t.Errorf("expected at least 4 categories, got %d", len(cats))
	}
}
