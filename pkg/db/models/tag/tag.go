package tag

import (
	"errors"

	"github.com/ryantking/marina/pkg/db"

	udb "upper.io/db.v3"
)

const (
	colName = "tag"
)

var (
	// ErrManifestNotFound is thrown when a manifest cannot be found
	ErrManifestNotFound = errors.New("no manifest for the given repo, org, and tag could be found")
)

// GetManifest returns the manifest with its type for a given tag
func GetManifest(name, repoName, orgName string) ([]byte, string, error) {
	col, err := db.GetCollection(colName)
	if err != nil {
		return nil, "", err
	}

	t := Model{}
	err = col.Find("name", name).And("repo_name", repoName).And("org_name", orgName).One(&t)
	if err == udb.ErrNoMoreRows {
		return nil, "", ErrManifestNotFound
	}

	return t.Manifest, t.ManifestType, nil
}

// UpdateManifest updates the manifest for a given tag, creating it if it does not exist
func UpdateManifest(name, repoName, orgName string, manifest []byte, manifestType string) error {
	col, err := db.GetCollection(colName)
	if err != nil {
		return err
	}

	t := Model{
		Name:         name,
		RepoName:     repoName,
		OrgName:      orgName,
		Manifest:     manifest,
		ManifestType: manifestType,
	}
	res := col.Find("name", name).And("repo_name", repoName).And("org_name", orgName)
	exists, err := res.Exists()
	if err != nil {
		return err
	}
	if exists {
		return res.Update(&t)
	}

	_, err = col.Insert(&t)
	return err
}
