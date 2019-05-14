package repo

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db"
)

// Model represents an entry in the database
type Model struct {
	Name    string `db:"name"`
	OrgName string `db:"org_name"`
}

const (
	// CollectionName is the name of the table in the database
	CollectionName = "repository"
)

var (
	getCollection = db.GetCollection
)

// Exists checks whether or not a given organization exists
func Exists(name, orgName string) (bool, error) {
	col, err := getCollection(CollectionName)
	if err != nil {
		return false, err
	}

	exists, err := col.Find("name", name).And("org_name", orgName).Exists()
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetNames returns the names of all available repos
func GetNames() ([]string, error) {
	col, err := getCollection(CollectionName)
	if err != nil {
		return nil, err
	}

	repos := make([]Model, 0)
	err = col.Find().OrderBy("org_name").OrderBy("repo_name").All(&repos)
	if err != nil {
		return []string{}, errors.Wrap(err, "error retrieving repos")
	}

	return repoNames(repos), nil
}

// GetNamesPaginated
func GetNamesPaginated(pageSize uint, last string) ([]string, string, error) {
	col, err := getCollection(CollectionName)
	if err != nil {
		return nil, "", err
	}

	repos := make([]Model, 0, pageSize)
	res := col.Find().OrderBy("org_name").OrderBy("repo_name").Paginate(pageSize).Cursor("name")
	err = res.NextPage(last).All(repos)
	if err != nil {
		return []string{}, "", errors.Wrap(err, "error retrieving repos")
	}
	nextLast := repos[len(repos)-1].Name
	exists, err := res.NextPage(nextLast).Exists()
	if err != nil {
		return []string{}, "", errors.Wrap(err, "error checking if next page exists")
	}
	if !exists {
		nextLast = ""
	}

	return repoNames(repos), nextLast, nil

}

func repoNames(repos []Model) []string {
	names := make([]string, len(repos))
	for i, repo := range repos {
		names[i] = fmt.Sprintf("%s/%s", repo.OrgName, repo.Name)
	}

	return names
}
