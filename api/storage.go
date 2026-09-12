package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/AliAlbhrani/StudentsArchive/engine"
	"github.com/AliAlbhrani/StudentsArchive/models"
	"github.com/danielgtaylor/huma/v2"
)

var (
	invalidContentType = huma.Error400BadRequest("invalid content type")

	maxFileSize int64 = 5 * 1024 * 1024 // 5MB
)

func UploadToStorage(ctx context.Context, i *models.StorageUpload) (*Response[string], error) {
	file := i.File.File

	if i.File.ContentType != "image/jpeg" && i.File.ContentType != "image/png" {
		slog.ErrorContext(ctx, "invalid content type", "contentType", i.File.ContentType)
		return nil, invalidContentType
	}
	if i.File.Size > maxFileSize {
		slog.ErrorContext(ctx, "file size too large", "size", i.File.Size)
		return nil, huma.Error413RequestEntityTooLarge("file size too large")
	}

	filename, err := engine.CS.Upload(file, i.File.Filename)
	if err != nil {
		slog.ErrorContext(ctx, "failed to upload to storage", "error", err.Error())
		return nil, somthingWentWrong
	}
	return SuccessResponse(filename, http.StatusOK), nil
}

func ReadFileFromStorage(ctx context.Context, i *models.PathParamID) (*huma.StreamResponse, error) {
	file, err := engine.CS.Download(i.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read file from storage", "error", err.Error())
		return nil, somthingWentWrong
	}

	return &huma.StreamResponse{
		Body: func(hctx huma.Context) {
			defer file.Body.Close()

			hctx.SetHeader("Content-Type", file.ContentType)
			hctx.SetHeader("Content-Disposition", `inline; filename="`+i.ID+`"`)

			writer := hctx.BodyWriter()
			if _, err := io.Copy(writer, file.Body); err != nil {
				slog.ErrorContext(ctx, "failed to stream file", "error", err.Error())
			}
		},
	}, nil
}

func initStorageRouter() {
	group := huma.NewGroup(API, "/storage")
	AuthHandle(group, "POST", "/upload", UploadToStorage)
	AuthHandle(group, "GET", "/upload", ReadFileFromStorage)
}
