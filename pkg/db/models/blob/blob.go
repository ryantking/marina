package blob

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db"
	udb "upper.io/db.v3"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "blob"
)

var col udb.Collection

// Collection returns the collection for the repository type
func Collection() (udb.Collection, error) {
	if col != nil {
		return col, nil
	}

	db, err := db.Get()
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving database connection")
	}
	col := db.Collection(CollectionName)
	if !col.Exists() {
		panic(fmt.Sprintf("collection '%s' does not exist", CollectionName))
	}

	return col, nil
}

// New returns a new blob
func New(digest, repoName, orgName string) (*Model, error) {
	col, err := Collection()
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving database connection")
	}
	l := Model{
		Digest:   digest,
		RepoName: repoName,
		OrgName:  orgName,
	}
	_, err = col.Insert(&l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// Exists checks whether or not a given organization exists
func Exists(digest string) (bool, error) {
	col, err := Collection()
	if err != nil {
		return false, errors.Wrap(err, "error retrieving collection")
	}

	exists, err := col.Find("digest", digest).Exists()
	if err != nil {
		return false, errors.Wrap(err, "error checking if row exists")
	}

	return exists, nil
}

// Delete deletes a blob
func Delete(digest string) error {
	col, err := Collection()
	if err != nil {
		return errors.Wrap(err, "error retrieving collection")
	}

	return col.Find("digest", digest).Delete()
}
