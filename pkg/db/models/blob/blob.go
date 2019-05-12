package blob

import (
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "blob"
)

var (
	getCollection = db.GetCollection
)

// Model represents an entry in the database
type Model struct {
	Digest   string `db:"digest"`
	RepoName string `db:"repo_name"`
	OrgName  string `db:"org_name"`
}

// Create creates a new blob in the database
func Create(digest, repoName, orgName string) error {
	col, err := getCollection(CollectionName)
	if err != nil {
		return errors.Wrap(err, "error retrieving database connection")
	}
	l := Model{
		Digest:   digest,
		RepoName: repoName,
		OrgName:  orgName,
	}
	_, err = col.Insert(&l)
	if err != nil {
		return errors.Wrap(err, "error inserting model into the database")
	}

	return nil
}

// Exists checks whether or not a given organization exists
func Exists(digest string) (bool, error) {
	col, err := getCollection(CollectionName)
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
	col, err := getCollection(CollectionName)
	if err != nil {
		return errors.Wrap(err, "error retrieving collection")
	}

	return col.Find("digest", digest).Delete()
}
