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

var (
	// ErrNoChunksForUpload is thrown when an upload does not have any chunks uploaded
	ErrNoChunksForUpload = errors.New("no upload chunks found")
)

// New creates a new upload chunk
func New(uuid string, start, end int64) error {
	col, err := db.GetCollection(CollectionName)
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
	col, err := db.GetCollection(CollectionName)
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

// GetLastRange returns the last chunk received for an upload
func GetLastRange(uuid string) (string, error) {
	col, err := db.GetCollection(CollectionName)
	if err != nil {
		return "", errors.Wrap(err, "error retrieving collection")
	}

	chunk := new(Model)
	err = col.Find("uuid", uuid).OrderBy("-range_start").Limit(1).One(chunk)
	if err == udb.ErrNoMoreRows {
		return "", ErrNoChunksForUpload
	}
	if err != nil {
		return "", errors.Wrap(err, "error retrieving upload chunks")
	}

	return fmt.Sprintf("%d-%d", chunk.RangeStart, chunk.RangeEnd), nil
}

// DeleteForUUID deletes all chunks for a given UUID
func DeleteForUUID(uuid string) error {
	col, err := db.GetCollection(CollectionName)
	if err != nil {
		return errors.Wrap(err, "error retrieving collection")
	}

	return col.Find("uuid", uuid).Delete()
}
