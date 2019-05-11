package repo

import (
	"fmt"

	"github.com/ryantking/marina/pkg/db"
	udb "upper.io/db.v3"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "repository"
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

// Exists checks whether or not a given organization exists
func Exists(name, orgName string) (bool, error) {
	col, err := Collection()
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
func GetNames(pageSize uint, last string) ([]string, string, error) {
	col, err := Collection()
	if err != nil {
		return nil, "", err
	}

	repos := make([]Model, 0)
	res := col.Find().OrderBy("org_name")
	resCur := res
	if pageSize > 0 {
		res = res.Paginate(pageSize).Cursor("name")
		resCur = res
		if last != "" {
			res = res.NextPage("")
		}
	}
	err = res.All(&repos)
	if err != nil {
		return nil, "", err
	}
	nextLast := ""
	if pageSize > 0 && uint(len(repos)) == pageSize {
		s := repos[len(repos)-1].Name
		exists, err := resCur.NextPage(s).Exists()
		if err != nil {
			return nil, "", err
		}
		if exists {
			nextLast = s
		}
	}

	names := make([]string, len(repos))
	for i, repo := range repos {
		names[i] = fmt.Sprintf("%s/%s", repo.OrgName, repo.Name)
	}

	return names, nextLast, nil
}
