package chunk

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db"
	udb "upper.io/db.v3"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "upload_chunk"
)

var col udb.Collection

// Collection returns the collection for the chunk type
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

// New creates a new upload chunk
func New(uuid string, start, end int64) error {
	col, err := Collection()
	if err != nil {
		return errors.Wrap(err, "error retrieving collection")
	}
	chunk := Model{UUID: uuid, RangeStart: start, RangeEnd: end}
	_, err = col.Insert(&chunk)
	if err != nil {
		return errors.Wrap(err, "error inserting upload chunk")
	}
	return nil
}

// GetAll returns all upload chunks for a UUID
func GetAll(uuid string) ([]*Model, error) {
	col, err := Collection()
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving collection")
	}
	chunks := make([]*Model, 0)
	err = col.Find("uuid", uuid).OrderBy("range_start").All(&chunks)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving upload chunks")
	}

	return chunks, nil
}
