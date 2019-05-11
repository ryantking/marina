package org

import (
	"github.com/ryantking/marina/pkg/db"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "organization"
)

// Exists checks whether or not a given organization exists
func Exists(name string) (bool, error) {
	col, err := db.GetCollection(CollectionName)
	if err != nil {
		return false, err
	}

	exists, err := col.Find("name", name).Exists()
	if err != nil {
		return false, err
	}

	return exists, nil
}
