package manifest

import (
	"context"
	"encoding/json"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/prisma"
)

func updateImage(c echo.Context, manifest docker.Manifest, repoID string) (*prisma.Image, error) {
	manifestType := c.Request().Header.Get(echo.HeaderContentType)
	b, err := json.Marshal(manifest)
	if err != nil {
		return nil, err
	}
	return client.UpsertImage(prisma.ImageUpsertParams{
		Where: prisma.ImageWhereUniqueInput{Digest: prisma.Str(manifest.Digest())},
		Update: prisma.ImageUpdateInput{
			Manifest:     prisma.Str(string(b)),
			ManifestType: &manifestType,
		},
		Create: prisma.ImageCreateInput{
			Digest:       manifest.Digest(),
			Manifest:     string(b),
			ManifestType: manifestType,
			Repo: prisma.RepositoryCreateOneWithoutImagesInput{
				Connect: &prisma.RepositoryWhereUniqueInput{ID: &repoID},
			},
		},
	}).Exec(c.Request().Context())
}

func replaceTag(ctx context.Context, ref, imageID string) error {
	err := deleteTag(ctx, ref, imageID)
	if err != nil {
		return err
	}
	_, err = client.CreateTag(prisma.TagCreateInput{
		Ref:   ref,
		Image: prisma.ImageCreateOneWithoutTagsInput{Connect: &prisma.ImageWhereUniqueInput{ID: &imageID}},
	}).Exec(ctx)
	return err
}

func deleteTag(ctx context.Context, ref, imageID string) error {
	_, err := client.DeleteManyTags(&prisma.TagWhereInput{
		Ref:   &ref,
		Image: &prisma.ImageWhereInput{ID: &imageID},
	}).Exec(ctx)
	return err
}
