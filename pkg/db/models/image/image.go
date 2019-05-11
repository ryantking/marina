package image

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/db"
	"github.com/ryantking/marina/pkg/db/models/tag"
	"github.com/ryantking/marina/pkg/docker"

	udb "upper.io/db.v3"
)

const (
	// CollectionName is the name of the table in the database
	CollectionName = "image"
)

var (
	// ErrManifestNotFound is thrown when a manifest cannot be found
	ErrManifestNotFound = errors.New("no manifest for the given repo, org, and reference could be found")

	// ErrDeleteOnTag is thrown when a delete is called on a tag not a digest
	ErrDeleteOnTag = errors.New("cannot delete by tag")
)

// GetManifest returns the manifest with its type for a given reference
func GetManifest(ref, repoName, orgName string) (docker.Manifest, string, error) {
	col, err := db.GetCollection(CollectionName)
	if err != nil {
		return nil, "", errors.Wrap(err, "error retrieving collection")
	}

	digest := ref
	if !strings.HasPrefix(ref, "sha256:") || len(ref) != 71 {
		digest, err = tag.GetDigest(ref, repoName, orgName)
		if err != nil {
			return nil, "", errors.Wrap(err, "error retireving digest")
		}
	}

	image := new(Model)
	err = col.Find("digest", digest).And("repo_name", repoName).And("org_name", orgName).One(&image)
	if err == udb.ErrNoMoreRows {
		return nil, "", ErrManifestNotFound
	}
	if err != nil {
		return nil, "", errors.Wrap(err, "error checking if digest exists")
	}
	man, err := docker.ParseManifest(image.ManifestType, image.Manifest)
	if err != nil {
		return nil, "", errors.Wrap(err, "error parsing manifest")
	}

	return man, image.ManifestType, nil
}

// UpdateManifest updates the manifest for a given tag, creating it if it does not exist
func UpdateManifest(ref, repoName, orgName string, manifest docker.Manifest, manifestType string) error {
	col, err := db.GetCollection(CollectionName)
	if err != nil {
		return errors.Wrap(err, "error retrieving collection")
	}

	digest := manifest.Digest()
	b, err := json.Marshal(manifest)
	if err != nil {
		return errors.Wrap(err, "error marshalling manifest")
	}
	image := &Model{
		Digest:       digest,
		RepoName:     repoName,
		OrgName:      orgName,
		Manifest:     b,
		ManifestType: manifestType,
	}
	res := col.Find("digest", digest)
	exists, err := res.Exists()
	if err != nil {
		return errors.Wrap(err, "error checking if manifest exists")
	}
	if exists {
		err := res.Update(image)
		if err != nil {
			return errors.Wrap(err, "error updating manifest")
		}
	} else {
		_, err := col.Insert(image)
		if err != nil {
			return errors.Wrap(err, "error creating manifest")
		}
	}

	if ref != digest {
		err := tag.Set(digest, ref, repoName, orgName)
		if err != nil {
			return errors.Wrap(err, "error updating tag")
		}
	}

	return nil
}

// Delete deletes a manifest from the table along with all tags
func Delete(digest string) error {
	col, err := db.GetCollection(CollectionName)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(digest, "sha256:") || len(digest) != 72 {
		return ErrDeleteOnTag
	}
	err = tag.DeleteDigest(digest)
	if err != nil {
		return errors.Wrap(err, "error deleting tags")
	}

	return col.Find("digest", digest).Delete()
}
