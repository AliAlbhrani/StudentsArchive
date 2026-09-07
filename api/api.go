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

var API huma.API

func Init() {

	cfg := huma.DefaultConfig("students-archive-api", "1.0.0")
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

	InitRoutes()

	log.Println("server running on port: ", env.PORT)
	http.ListenAndServe(fmt.Sprint(":", env.PORT), router)
}

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

type Success[O any] struct {
	Status int `json:"status"`
	Body   O   `json:"data"`
}

func SuccessResponse[O any](data O, status int) *Success[O] {
	return &Success[O]{
		Status: status,
		Body:   data,
	}
}

func InitRoutes() {
	InitUsersRoutes()
}
