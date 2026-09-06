package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AliAlbhrani/StudentsArchive/engine"
	"github.com/AliAlbhrani/StudentsArchive/sqlc"
	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
)

type StageEnum int64

const (
	StageEnumFirst  StageEnum = 1
	StageEnumSecond StageEnum = 2
	StageEnumThird  StageEnum = 3
	StageEnumFourth StageEnum = 4
	StageEnumFifth  StageEnum = 5
)

type User struct {
	FullName  string `json:"full_name" validate:"required" min:"3" max:"255"`
	Username  string `json:"username" validate:"required" min:"3" max:"255"`
	Password  string `json:"password" validate:"required" min:"6" max:"255"`
	Stage     int64  `json:"stage" validate:"required" min:"1" max:"5"`
	CreatedAt string `json:"created_at"`
}

type CreateUserDto struct {
	Body struct {
		User
	}
}

// Register handles user registration
func register(ctx context.Context, createUser *CreateUserDto) (*Success[TokensResponse], error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUser.Body.User.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.ErrorContext(ctx, "failed to hash password", "error", err)
		return nil, huma.Error500InternalServerError("failed to hash password")
	}
	id, err := engine.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		FullName: createUser.Body.User.FullName,
		Username: createUser.Body.User.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create user", "error", err)
		return nil, huma.Error400BadRequest("failed to create user")
	}
	tokens, err := GenerateToken(int64(id))
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate tokens", "error", err)
		return nil, huma.Error500InternalServerError("failed to generate tokens")
	}
	return SuccessResponse(*tokens, 200), nil
}

func InitUsersRoutes() {
	usersRoutes := huma.NewGroup(API, "/users")
	Handle(usersRoutes, http.MethodPost, "/register", register)
}
