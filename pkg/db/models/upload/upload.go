package upload

import (
	"github.com/google/uuid"
	"github.com/ryantking/marina/pkg/db"
	"github.com/ryantking/marina/pkg/db/models/upload/chunk"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "upload"
)

// New creates a new upload and returns its UUID
func New() (string, error) {
	col, err := db.GetCollection(CollectionName)
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

// Finish sets an upload's status to done
func Finish(uuid string) error {
	col, err := db.GetCollection(CollectionName)
	if err != nil {
		return err
	}

	return col.Find("uuid", uuid).Update(map[string]bool{"done": true})
}

// Delete deletes the upload and all chunks
func Delete(uuid string) error {
	col, err := db.GetCollection(CollectionName)
	if err != nil {
		return err
	}

	err = chunk.DeleteForUUID(uuid)
	if err != nil {
		return err
	}

	return col.Find("uuid", uuid).Delete()
}
