package upload

import "github.com/ryantking/marina/pkg/db"

const (
	colName = "upload"
)

// New returns a new upload
func New() (*Model, error) {
	upl := Model{}
	col, err := db.GetCollection(colName)
	if err != nil {
		return nil, err
	}
	err = col.InsertReturning(&upl)
	if err != nil {
		return nil, err
	}

	return &upl, nil
}

// Save saves the upload to the database
func (upl *Model) Save() error {
	col, err := db.GetCollection(colName)
	if err != nil {
		return err
	}
	return col.Find("uuid", upl.UUID).Update(upl)
}
