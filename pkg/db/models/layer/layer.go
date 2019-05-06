package layer

import "github.com/ryantking/marina/pkg/db"

const (
	colName = "layer"
)

// New returns a new layer
func New(digest, repoName, orgName string) (*Model, error) {
	col, err := db.GetCollection(colName)
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

// Exists returns whether or not a layer exists
func Exists(digest, repoName, orgName string) (bool, error) {
	col, err := db.GetCollection(colName)
	if err != nil {
		return false, err
	}
	exists, err := col.Find("digest", digest).And("repo_name", repoName).And("org_name", orgName).Exists()
	if err != nil {
		return false, err
	}

	return exists, nil
}
