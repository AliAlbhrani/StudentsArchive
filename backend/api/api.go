package api

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/AliAlbhrani/StudentsArchive/env"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

var (
	somthingWentWrong error = huma.Error500InternalServerError("something went wrong")
	unauthorized      error = huma.Error401Unauthorized("unauthorized")
	notFound          error = huma.Error404NotFound("not found")
)

var API huma.API

func Init() {

	cfg := huma.DefaultConfig("students-archive-api", "1.0.0")
	cfg.CreateHooks = nil
	cfg.DocsRenderer = huma.DocsRendererScalar
	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"Authorization": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	router := http.NewServeMux()
	API = humago.New(router, cfg)

	initRoutes()

	log.Println("server running on port: ", env.PORT)
	http.ListenAndServe(fmt.Sprint(":", env.PORT), router)
}

// register a route with authentication middleware
func AuthHandle[I, O any](group *huma.Group, method string, pattern string, handler func(context.Context, *I) (*O, error)) {
	op := huma.Operation{
		Method:      method,
		Path:        pattern,
		Middlewares: huma.Middlewares{AuthMiddleware},
	}
	huma.Register(group, op, handler)
	slog.Debug("registered route", "method", method, "path", pattern, "auth", true)
}

func Handle[I, O any](group *huma.Group, method string, pattern string, handler func(context.Context, *I) (*O, error)) {
	op := huma.Operation{
		Method: method,
		Path:   pattern,
	}
	huma.Register(group, op, handler)
	slog.Debug("registered route", "method", method, "path", pattern, "auth", false)
}

// Success is a generic success response type
type Response[O any] struct {
	Status int `json:"status"`
	Body   O   `json:"data"`
}

// PaginatedResponse is a generic pagination response type
type PaginatedResponse[O any] struct {
	Status int   `json:"status"`
	Data   []O   `json:"data"`
	Total  int64 `json:"total"`
}

func SuccessResponse[O any](data O, status int) *Response[O] {
	return &Response[O]{
		Status: status,
		Body:   data,
	}
}

func PaginationResponse[O any](data []O, total int64, status int) *PaginatedResponse[O] {
	return &PaginatedResponse[O]{
		Status: status,
		Data:   data,
		Total:  total,
	}
}

type PaginationRequest struct {
	Page   int    `json:"page" default:"1" min:"1"`
	Limit  int    `json:"limit" default:"10" min:"1" max:"10"`
	Search string `json:"search"`
}

func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}

func initRoutes() {
	initUsersRoutes()
	initPostsRouter()
	initStorageRouter()
}
