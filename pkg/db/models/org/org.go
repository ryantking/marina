package org

import (
	"fmt"

	"github.com/ryantking/marina/pkg/db"
	udb "upper.io/db.v3"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "organization"
)

var col udb.Collection

// Collection returns the collection for the organization type
func Collection() (udb.Collection, error) {
	if col != nil {
		return col, nil
	}

	db, err := db.Get()
	if err != nil {
		return nil, err
	}
	col := db.Collection(CollectionName)
	if !col.Exists() {
		panic(fmt.Sprintf("collection '%s' does not exist", CollectionName))
	}

	return col, nil
}

// Exists checks whether or not a given organization exists
func Exists(name string) (bool, error) {
	col, err := Collection()
	if err != nil {
		return false, err
	}

	exists, err := col.Find("name", name).Exists()
	if err != nil {
		return false, err
	}

	return exists, nil
}
