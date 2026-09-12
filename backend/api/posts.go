// Package api handles the HTTP API for the students archive.
package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AliAlbhrani/StudentsArchive/engine"
	"github.com/AliAlbhrani/StudentsArchive/models"
	"github.com/AliAlbhrani/StudentsArchive/sqlc"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
)

func GetAllPostsHandler(ctx context.Context, i *PaginationRequest) (*PaginatedResponse[sqlc.Post], error) {
	posts, err := engine.Queries.GetPosts(ctx, sqlc.GetPostsParams{
		Limit:  i.Limit,
		Offset: i.Offset(),
		Search: i.Search,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		slog.ErrorContext(ctx, "failed to get posts", "error", err.Error())
		return nil, somthingWentWrong
	}
	total, err := engine.Queries.GetPostsCount(ctx, i.Search)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get posts count", "error", err.Error())
		return nil, somthingWentWrong
	}
	return PaginationResponse(posts, total, http.StatusOK), nil
}

func CreatePostHandler(ctx context.Context, req *models.CreatePostRequest) (*Response[string], error) {
	req.UserID = GetUserID(ctx)
	if req.UserID == 0 {
		return nil, unauthorized
	}
	if len(req.Body.Images) > 3 {
		return nil, huma.Error400BadRequest("images must be less than or equal to 3")
	}
	err := engine.Queries.CreatePost(ctx, sqlc.CreatePostParams{
		UserID:  &req.UserID,
		Title:   req.Body.Title,
		Content: &req.Body.Content,
		Images:  req.Body.Images,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create post", "error", err.Error())
		return nil, somthingWentWrong
	}
	return SuccessResponse("post created successfully", http.StatusOK), nil
}

func initPostsRouter() {
	group := huma.NewGroup(API, "/posts")
	AuthHandle(group, http.MethodGet, "", GetAllPostsHandler)
	AuthHandle(group, http.MethodPost, "", CreatePostHandler)
}
