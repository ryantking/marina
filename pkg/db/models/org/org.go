package org

import "github.com/ryantking/marina/pkg/db"

const (
	colName    = "organization"
	defaultOrg = "library"
)

// CreateIfNotExists creates an organization if not currently present
func CreateIfNotExists(name string) error {
	db, err := db.Get()
	if err != nil {
		return err
	}
	col := db.Collection(colName)
	exists, err := col.Find("name", name).Exists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	org := Model{
		Name: name,
	}
	_, err = col.Insert(org)
	return err
}
