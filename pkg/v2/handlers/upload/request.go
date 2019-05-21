package upload

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/prisma"
	"github.com/ryantking/marina/pkg/web"
)

var orderBy = prisma.ChunkOrderByInputRangeEndDesc

func parsePath(c echo.Context) (string, string, string, error) {
	repo, org, err := web.ParsePath(c)
	if err != nil {
		return "", "", "", err
	}

	uuid := c.Param("uuid")
	upload, err := client.Upload(prisma.UploadWhereUniqueInput{Uuid: &uuid}).Exec(c.Request().Context())
	if err == prisma.ErrNoResult {
		c.Set("docker_err_code", docker.CodeBlobUploadUnknown)
		return "", "", "", echo.NewHTTPError(http.StatusNotFound, "blob upload unknown to registry")
	}
	if err != nil {
		return "", "", "", err
	}

	return upload.Uuid, repo, org, nil
}

func getLastRange(ctx context.Context, uuid string) (string, error) {
	chunks, err := client.Chunks(&prisma.ChunksParams{
		First:   prisma.Int32(1),
		OrderBy: &orderBy,
		Where: &prisma.ChunkWhereInput{
			Upload: &prisma.UploadWhereInput{
				Uuid: &uuid,
			},
		},
	}).Exec(ctx)
	if len(chunks) == 0 {
		return "0-0", nil
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d-%d", chunks[0].RangeStart, chunks[0].RangeEnd), nil
}

func parseLength(c echo.Context) (int64, error) {
	s := c.Request().Header.Get(echo.HeaderContentLength)
	if s == "" {
		return -1, nil
	}

	sz, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return sz, nil
}

func parseRange(c echo.Context) (int64, error) {
	s := c.Request().Header.Get("Content-Range")
	if s == "" {
		return 0, nil
	}
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid Content-Range header")
	}
	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		c.Set("docker_err_code", "BLOB_UPLOAD_INVALID")
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return start, nil
}
