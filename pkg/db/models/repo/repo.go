package repo

import (
	"github.com/ryantking/marina/pkg/db"
	"github.com/ryantking/marina/pkg/db/models/org"
)

const (
	colName = "repository"
)

// CreateIfNotExists creates a repository if not currently present
func CreateIfNotExists(repoName, orgName string) error {
	err := org.CreateIfNotExists(orgName)
	if err != nil {
		return err
	}

	col, err := db.GetCollection(colName)
	if err != nil {
		return err
	}
	exists, err := col.Find("name", repoName).And("org_name", orgName).Exists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	repo := Model{
		Name:    repoName,
		OrgName: orgName,
	}
	_, err = col.Insert(repo)
	return err
}
