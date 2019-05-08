package layer

import (
	"fmt"

	"github.com/ryantking/marina/pkg/db"
	udb "upper.io/db.v3"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "layer"
)

var col udb.Collection

// Collection returns the collection for the repository type
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

// New returns a new layer
func New(digest, repoName, orgName string) (*Model, error) {
	col, err := Collection()
	if err != nil {
		return nil, err
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
func Exists(digest, repoName, orgName string) (bool, error) {
	col, err := Collection()
	if err != nil {
		return false, err
	}

	exists, err := col.Find("digest", digest).And("repo_name", repoName).And("org_name", orgName).Exists()
	if err != nil {
		return false, err
	}

	return exists, nil
}
