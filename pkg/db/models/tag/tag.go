package tag

import (
	"errors"
	"fmt"

	"github.com/ryantking/marina/pkg/db"

	udb "upper.io/db.v3"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "tag"
)

var (
	col udb.Collection

	// ErrTagNotFound is thrown when a tag doesn't exist for a given repo and org
	ErrTagNotFound = errors.New("no tag with name for the given repo and org")
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

// GetDigest gets the associated digest for a tag
func GetDigest(name, repoName, orgName string) (string, error) {
	col, err := Collection()
	if err != nil {
		return "", err
	}

	tag := new(Model)
	err = col.Find("name", name).And("repo_name", repoName).And("org_name", orgName).One(tag)
	if err == udb.ErrNoMoreRows {
		return "", ErrTagNotFound
	}
	if err != nil {
		return "", err
	}

	return tag.ImageDigest, nil
}

// Set assocites a given tag with a digest
func Set(digest, name, repoName, orgName string) error {
	col, err := Collection()
	if err != nil {
		return err
	}

	res := col.Find("name", name).And("repo_name", repoName).And("org_name", orgName)
	exists, err := res.Exists()
	if err != nil {
		return err
	}
	if exists {
		return res.Update(map[string]string{"image_digest": digest})
	}

	tag := &Model{
		Name:        name,
		RepoName:    repoName,
		OrgName:     orgName,
		ImageDigest: digest,
	}
	_, err = col.Insert(tag)
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
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	return tagNames, nextLast, nil
}

// DeleteDigest deletes all tags for a given digest
func DeleteDigest(digest string) error {
	col, err := Collection()
	if err != nil {
		return err
	}

	return col.Find("digest", digest).Delete()
}
