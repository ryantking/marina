package org

import (
	"github.com/jinzhu/gorm"
	"github.com/ryantking/marina/pkg/db"
)

// Model represents a single entry in the database
type Model struct {
	gorm.Model
	Name string `gorm:"column:name"`
}

// TableName returns the organization table name
func (Model) TableName() string {
	return "organizations"
}

// Exists checks whether or not a given organization exists
func Exists(name string) (bool, error) {
	db := db.Get()
	res := db.Where("name = ?", name).First(&Model{})
	if res.RecordNotFound() {
		return false, nil
	}
	if res.Error != nil {
		return false, res.Error
	}

	return true, nil
}
