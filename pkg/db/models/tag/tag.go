package tag

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ryantking/marina/pkg/db"
	"github.com/ryantking/marina/pkg/docker"

	udb "upper.io/db.v3"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "tag"
)

var (
	col udb.Collection

	// ErrManifestNotFound is thrown when a manifest cannot be found
	ErrManifestNotFound = errors.New("no manifest for the given repo, org, and tag could be found")
)

// Collection returns the collection for the organization type
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

// GetManifest returns the manifest with its type for a given tag
func GetManifest(name, repoName, orgName string) (docker.Manifest, string, error) {
	col, err := Collection()
	if err != nil {
		return nil, "", err
	}

	var t Model
	err = col.Find("name", name).And("repo_name", repoName).And("org_name", orgName).One(&t)
	if err == udb.ErrNoMoreRows {
		return nil, "", ErrManifestNotFound
	}
	if err != nil {
		return nil, "", err
	}
	man, err := docker.ParseManifest(t.ManifestType, t.Manifest)
	if err != nil {
		return nil, "", err
	}

	return man, t.ManifestType, nil
}

// UpdateManifest updates the manifest for a given tag, creating it if it does not exist
func UpdateManifest(name, repoName, orgName string, manifest docker.Manifest, manifestType string) error {
	col, err := Collection()
	if err != nil {
		return err
	}

	b, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	t := Model{
		Name:         name,
		RepoName:     repoName,
		OrgName:      orgName,
		Manifest:     b,
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

// List lists all tags for a given repo
func List(repoName, orgName string, pageSize uint, last string) ([]string, string, error) {
	col, err := Collection()
	if err != nil {
		return []string{}, "", err
	}

	tags := []Model{}
	res := col.Find("repo_name", repoName).And("org_name", orgName).OrderBy("-_last_updated")
	resCur := res
	if pageSize > 0 {
		res = res.Paginate(pageSize).Cursor("name")
		resCur = res
		if last != "" {
			res = res.NextPage(last)
		}
	}
	err = res.All(&tags)
	if err != nil {
		return []string{}, "", err
	}
	nextLast := ""
	if pageSize > 0 && uint(len(tags)) == pageSize {
		s := tags[len(tags)-1].Name
		exists, err := resCur.NextPage(s).Exists()
		if err != nil {
			return []string{}, "", err
		}
		if exists {
			nextLast = s
		}
	}
	tagNames := make([]string, len(tags), len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	return tagNames, nextLast, nil
}
