package main

import (
	"context"

	"github.com/ryantking/marina/pkg/prisma"
)

func main() {
	client := prisma.New(nil)
	ctx := context.Background()

	_, err := client.CreateOrganization(prisma.OrganizationCreateInput{
		Name: "library",
		Repos: &prisma.RepositoryCreateManyWithoutOrgInput{
			Create: []prisma.RepositoryCreateWithoutOrgInput{
				prisma.RepositoryCreateWithoutOrgInput{
					Name: "alpine",
					Blobs: &prisma.BlobCreateManyWithoutRepoInput{
						Create: []prisma.BlobCreateWithoutRepoInput{
							prisma.BlobCreateWithoutRepoInput{
								Digest: "sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
							},
						},
					},
				},
				prisma.RepositoryCreateWithoutOrgInput{
					Name: "redis",
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
				prisma.RepositoryCreateWithoutOrgInput{Name: "nginx"},
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
				prisma.RepositoryCreateWithoutOrgInput{Name: "mysql"},
				prisma.RepositoryCreateWithoutOrgInput{Name: "mysql-client"},
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
}
