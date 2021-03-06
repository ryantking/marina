package testutil

import (
	"context"

	"github.com/ryantking/marina/pkg/prisma"

	// MySQL driver needed for tests
	_ "github.com/go-sql-driver/mysql"
)

var (
	client = prisma.New(nil)
)

func Seed(ctx context.Context) {
	_, err := client.CreateOrganization(prisma.OrganizationCreateInput{
		Name: "library",
		Repos: &prisma.RepositoryCreateManyWithoutOrgInput{
			Create: []prisma.RepositoryCreateWithoutOrgInput{
				prisma.RepositoryCreateWithoutOrgInput{
					Name:     "alpine",
					FullName: "library/alpine",
					Blobs: &prisma.BlobCreateManyWithoutRepoInput{
						Create: []prisma.BlobCreateWithoutRepoInput{
							prisma.BlobCreateWithoutRepoInput{
								Digest: "sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
							},
						},
					},
					Images: &prisma.ImageCreateManyWithoutRepoInput{
						Create: []prisma.ImageCreateWithoutRepoInput{
							prisma.ImageCreateWithoutRepoInput{
								Digest:       "sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb",
								Manifest:     `{"config": {"digest": "sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb"}}`,
								ManifestType: "application/vnd.docker.distribution.manifest.v2+json",
								Tags: &prisma.TagCreateManyWithoutImageInput{
									Create: []prisma.TagCreateWithoutImageInput{
										prisma.TagCreateWithoutImageInput{Ref: "3.9"},
									},
								},
							},
						},
					},
				},
				prisma.RepositoryCreateWithoutOrgInput{
					Name:     "redis",
					FullName: "library/redis",
					Blobs: &prisma.BlobCreateManyWithoutRepoInput{
						Create: []prisma.BlobCreateWithoutRepoInput{
							prisma.BlobCreateWithoutRepoInput{
								Digest: "sha256:6270adb5794c6987109e54af00ab456977c5d5cc6f1bc52c1ce58d32ec0f15f4",
							},
							prisma.BlobCreateWithoutRepoInput{
								Digest: "sha256:f99f83132c0a553c6444af96c9f8b905bf1e12835f4f1fa1aa67e9b39b461b1b",
							},
							prisma.BlobCreateWithoutRepoInput{
								Digest: "sha256:03eafa7928765518eb3d5f2bd2bd97dac655833dd088db9e7590eadbb721aa72",
							},
							prisma.BlobCreateWithoutRepoInput{
								Digest: "sha256:dc5ef0c854b450e002c49b8d3dad8a70b7d70a695f59104d6dea0bb87582baa2",
							},
							prisma.BlobCreateWithoutRepoInput{
								Digest: "sha256:4c8e44b308d0a525ccda8da176d6daaa6c117c87148ae85396a8171e6996f9f4",
							},
							prisma.BlobCreateWithoutRepoInput{
								Digest: "sha256:76c5d4259e40d4282ba496be18926dcef62063f99eb463ed2d802dd9a0e303a0",
							},
						},
					},
				},
				prisma.RepositoryCreateWithoutOrgInput{Name: "nginx", FullName: "library/nginx"},
			},
		},
	}).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}

	_, err = client.CreateOrganization(prisma.OrganizationCreateInput{
		Name: "mysql",
		Repos: &prisma.RepositoryCreateManyWithoutOrgInput{
			Create: []prisma.RepositoryCreateWithoutOrgInput{
				prisma.RepositoryCreateWithoutOrgInput{Name: "mysql", FullName: "mysql/mysql"},
				prisma.RepositoryCreateWithoutOrgInput{Name: "mysql-client", FullName: "mysql/mysql-client"},
			},
		},
	}).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}

	_, err = client.CreateUpload(prisma.UploadCreateInput{
		Uuid: "6b3c9a93-af5d-473f-a4ce-9710022185cd",
		Chunks: &prisma.ChunkCreateManyWithoutUploadInput{
			Create: []prisma.ChunkCreateWithoutUploadInput{
				prisma.ChunkCreateWithoutUploadInput{RangeStart: 0, RangeEnd: 1023},
				prisma.ChunkCreateWithoutUploadInput{RangeStart: 1024, RangeEnd: 2047},
			},
		},
	}).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}

	_, err = client.CreateUpload(prisma.UploadCreateInput{Uuid: "3f497dc6-9458-4c2d-8368-2e71d35c77e5"}).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}
}

func Clear(ctx context.Context) {
	_, err := client.DeleteManyTags(nil).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}
	_, err = client.DeleteManyImages(nil).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}
	_, err = client.DeleteManyChunks(nil).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}
	_, err = client.DeleteManyUploads(nil).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}
	_, err = client.DeleteManyBlobs(nil).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}
	_, err = client.DeleteManyRepositories(nil).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}
	_, err = client.DeleteManyOrganizations(nil).Exec(ctx)
	if err != nil {
		panic(err.Error())
	}
}
