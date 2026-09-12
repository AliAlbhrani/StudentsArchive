package models

import "github.com/danielgtaylor/huma/v2"

type StorageUpload struct {
	File huma.FormFile `form:"file"`
}

type PathParamID struct {
	ID string `path:"id"`
}
