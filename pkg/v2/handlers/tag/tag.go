package tag

import (
	"context"

	"github.com/ryantking/marina/pkg/prisma"
)

var (
	client  = prisma.New(nil)
	orderBy = prisma.TagOrderByInputUpdatedAtDesc
)

func getTags(ctx context.Context, repoID string, n int32, last string) ([]prisma.Tag, error) {
	tagsParams := &prisma.TagsParams{
		OrderBy: &orderBy,
		Where: &prisma.TagWhereInput{
			Image: &prisma.ImageWhereInput{
				Repo: &prisma.RepositoryWhereInput{ID: &repoID},
			},
		},
	}
	if n != 0 {
		tagsParams.First = &n
	}
	if last != "" {
		tags, err := client.Tags(&prisma.TagsParams{
			Where: &prisma.TagWhereInput{
				Ref: &last,
				Image: &prisma.ImageWhereInput{
					Repo: &prisma.RepositoryWhereInput{ID: &repoID},
				},
			},
		}).Exec(ctx)
		if err != nil {
			return nil, err
		}
		tagsParams.After = &tags[0].ID
	}

	return client.Tags(tagsParams).Exec(ctx)
}
