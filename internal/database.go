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
	"fmt"
	"errors"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//--------------------------------------------------core-------------------------------------------------------------------------------------------------//

// Database owns the gorm connection for internal data access.
type Database struct {
	conn *gorm.DB
}

// NewDatabase opens (or creates) the sqlite file and runs schema migrations.
func NewDatabase(path string) (*Database, error) {
	if path == "" {
		path = "rete.db"
	}

	conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := conn.AutoMigrate(
		&FeatureCategory{},
		&Feature{},
		&ScanJob{},
		&ScanResult{},
	); err != nil {
		return nil, fmt.Errorf("auto-migrate schema: %w", err)
	}

	return &Database{conn: conn}, nil
}

// Conn exposes the raw gorm handle for advanced queries/transactions.
func (d *Database) Conn() *gorm.DB {
	return d.conn
}

//-----------------------------------------------------------models and types------------------------------------------------------------------------------------------------//

// FeatureCategory groups related capabilities (e.g. "Reconnaissance & Discovery").
type FeatureCategory struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:255;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Features  []Feature `gorm:"foreignKey:CategoryID"`
}

// Feature is a single capability entry sourced from table.csv.
type Feature struct {
	ID               uint            `gorm:"primaryKey;autoIncrement"`
	CategoryID       uint            `gorm:"not null;index"`
	Category         FeatureCategory `gorm:"foreignKey:CategoryID"`
	Name             string          `gorm:"size:255;not null"`
	CobraCommand     string          `gorm:"size:255;not null"`
	ShortDescription string          `gorm:"type:text"`
	CreatedAt        time.Time       `gorm:"autoCreateTime"`
}

// JobStatus represents the lifecycle state of a ScanJob.
type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusRunning  JobStatus = "running"
	JobStatusDone     JobStatus = "done"
	JobStatusFailed   JobStatus = "failed"
	JobStatusCanceled JobStatus = "canceled"
)

// ScanJob records an invocation of a diagnostic or security command.
type ScanJob struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	FeatureID *uint      `gorm:"index"`
	Feature   *Feature   `gorm:"foreignKey:FeatureID"`
	Target    string     `gorm:"size:512"`
	Command   string     `gorm:"size:512;not null"`
	Args      string     `gorm:"type:text"`
	Status    JobStatus  `gorm:"size:32;default:'pending';not null"`
	StartedAt *time.Time `gorm:"column:started_at"`
	EndedAt   *time.Time `gorm:"column:ended_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	Results   []ScanResult `gorm:"foreignKey:JobID"`
}

// ScanResult stores the output lines produced by a ScanJob.
type ScanResult struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	JobID     uint      `gorm:"not null;index"`
	Job       ScanJob   `gorm:"foreignKey:JobID"`
	Output    string    `gorm:"type:text"`
	IsError   bool      `gorm:"default:false;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

//-----------------------------------------------------------feature CRUD-------------------------------------------//

// UpsertCategory creates a category if it does not exist, returning the persisted record.
func (d *Database) UpsertCategory(name string) (*FeatureCategory, error) {
	var cat FeatureCategory
	res := d.conn.Where(FeatureCategory{Name: name}).FirstOrCreate(&cat, FeatureCategory{Name: name})
	if res.Error != nil {
		return nil, res.Error
	}
	return &cat, nil
}

// UpsertFeature creates a feature if one with the same CobraCommand does not exist.
func (d *Database) UpsertFeature(f *Feature) error {
	if f == nil {
		return errors.New("feature is nil")
	}
	var existing Feature
	res := d.conn.Where("cobra_command = ?", f.CobraCommand).First(&existing)
	if res.Error == nil {
		return nil // already seeded
	}
	return d.conn.Create(f).Error
}

// ListCategories returns all feature categories.
func (d *Database) ListCategories() ([]*FeatureCategory, error) {
	var cats []*FeatureCategory
	if err := d.conn.Preload("Features").Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

// ListFeatures returns all features, optionally filtered by category name.
func (d *Database) ListFeatures(category string) ([]*Feature, error) {
	q := d.conn.Preload("Category")
	if category != "" {
		q = q.Joins("JOIN feature_categories ON feature_categories.id = features.category_id").
			Where("feature_categories.name = ?", category)
	}
	var features []*Feature
	if err := q.Find(&features).Error; err != nil {
		return nil, err
	}
	return features, nil
}

//-----------------------------------------------------------job CRUD------------------------------------------------//

// CreateJob persists a new ScanJob.
func (d *Database) CreateJob(job *ScanJob) error {
	if job == nil {
		return errors.New("job is nil")
	}
	return d.conn.Create(job).Error
}

// UpdateJob saves changes to an existing ScanJob.
func (d *Database) UpdateJob(job *ScanJob) error {
	if job == nil {
		return errors.New("job is nil")
	}
	if job.ID == 0 {
		return errors.New("job id is required")
	}
	return d.conn.Save(job).Error
}

// ListJobs returns recent scan jobs, newest first.
func (d *Database) ListJobs(limit int) ([]*ScanJob, error) {
	if limit <= 0 {
		limit = 20
	}
	var jobs []*ScanJob
	if err := d.conn.Preload("Feature").Order("id DESC").Limit(limit).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetJobByID fetches a job by primary key, preloading its results.
func (d *Database) GetJobByID(id uint) (*ScanJob, error) {
	var job ScanJob
	if err := d.conn.Preload("Results").Preload("Feature").First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// AddResult appends an output line to a scan job.
func (d *Database) AddResult(result *ScanResult) error {
	if result == nil {
		return errors.New("result is nil")
	}
	return d.conn.Create(result).Error
}