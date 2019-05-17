package org

import (
	"github.com/jinzhu/gorm"
	"github.com/ryantking/marina/pkg/db"
)

// Org represents a single entry in the database
type Org struct {
	gorm.Model
	Name  string `gorm:"column:name"`
	Repos []Repo `gorm:"foreignkey:org_id"`
}

// TableName returns the organization table name
func (Org) TableName() string {
	return "organizations"
}

// Repo represents an entry in the database
type Repo struct {
	gorm.Model
	Name  string `gorm:"column:name"`
	OrgID uint
}

// TableName returns the organization table name
func (Repo) TableName() string {
	return "repositories"
}

// OrgExists checks whether or not a given organization exists
func OrgExists(name string) (bool, error) {
	res := db.Get().Table("organizations").First(&Org{}, "name = ?", name)
	if res.RecordNotFound() {
		return false, nil
	}
	if res.Error != nil {
		return false, res.Error
	}

	return true, nil
}

// Exists checks whether or not a given repository exists
func Exists(repoName, orgName string) (bool, error) {
	res := db.Get().Table("repositories").Joins("JOIN organizations ON repositories.org_id = organizations.id").
		First(&Repo{}, "repositories.name = ? AND organizations.name = ?", repoName, orgName)
	if res.RecordNotFound() {
		return false, nil
	}
	if res.Error != nil {
		return false, res.Error
	}

	return true, nil
}
