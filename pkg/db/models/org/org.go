package org

import (
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db"
)

// Model represents a single entry in the database
type Model struct {
	Name string `db:"name"`
}

const (
	// CollectionName is the name of the table in the database
	CollectionName = "organization"
)

var (
	getCollection = db.GetCollection
)

// Exists checks whether or not a given organization exists
func Exists(name string) (bool, error) {
	col, err := getCollection(CollectionName)
	if err != nil {
		return false, errors.Wrap(err, "error retrieving collection")
	}

	exists, err := col.Find("name", name).Exists()
	if err != nil {
		return false, errors.Wrap(err, "error checking if org exists")
	}

	return exists, nil
}
