package upload

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ryantking/marina/pkg/db"
	udb "upper.io/db.v3"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "upload"
)

var col udb.Collection

// Collection returns the collection for the upload type
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

// New creates a new upload and returns its UUID
func New() (string, error) {
	col, err := Collection()
	if err != nil {
		return "", err
	}
	uuid := uuid.New()
	upl := Model{UUID: uuid.String()}
	_, err = col.Insert(&upl)
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// Save saves the upload to the database
func (upl *Model) Save() error {
	col, err := Collection()
	if err != nil {
		return err
	}
	return col.Find("uuid", upl.UUID).Update(upl)
}
