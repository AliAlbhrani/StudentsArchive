package api

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

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

type CreateUserDto struct {
	Body struct {
		FullName string `json:"full_name" validate:"required" min:"3" max:"255"`
		Username string `json:"username" validate:"required" min:"3" max:"255"`
		Password string `json:"password" validate:"required" min:"6" max:"255"`
		Stage    int64  `json:"stage" validate:"required" min:"1" max:"5"`
	}
}

// Register handles user registration
func register(ctx context.Context, createUser *CreateUserDto) (*Success[TokensResponse], error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUser.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.ErrorContext(ctx, "failed to hash password", "error", err)
		return nil, huma.Error500InternalServerError("failed to hash password")
	}
	id, err := engine.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		FullName: createUser.Body.FullName,
		Username: createUser.Body.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create user", "error", err)
		return nil, huma.Error400BadRequest("failed to create user")
	}
	tokens, err := GenerateToken(int(id))
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate tokens", "error", err)
		return nil, huma.Error500InternalServerError("failed to generate tokens")
	}
	return SuccessResponse(*tokens, 200), nil
}

type CreateUserProfileDto struct {
	UserID int
	Body   struct {
		PhotoURL *string `json:"photo_url" validate:"required" min:"3" max:"255"`
		Bio      *string `json:"bio" validate:"required" min:"3" max:"255"`
	}
}

// CreateUserProfile handles creating a user profile
func CreateUserProfile(ctx context.Context, createUserProfile *CreateUserProfileDto) (*Success[string], error) {
	var ok bool
	createUserProfile.UserID, ok = ctx.Value(userIDKey).(int)
	if !ok {
		slog.ErrorContext(ctx, "user id not found")
		return nil, huma.Error400BadRequest("user id not found")
	}
	err := engine.Queries.CreateUserProfile(ctx, sqlc.CreateUserProfileParams{
		UserID:   &createUserProfile.UserID,
		PhotoUrl: createUserProfile.Body.PhotoURL,
		Bio:      createUserProfile.Body.Bio,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create user profile", "error", err)
		return nil, huma.Error500InternalServerError("failed to create user profile")
	}
	return SuccessResponse("user profile created", 200), nil
}

type GetUserProfileDto struct {
	UserID    int
	FullName  string    `json:"full_name"`
	Username  string    `json:"username"`
	Bio       *string   `json:"bio"`
	Stage     int       `json:"stage"`
	PhotoURL  *string   `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}

// GetUserProfile returns the user profile by user id
func GetUserProfile(ctx context.Context, _ *struct{}) (*Success[GetUserProfileDto], error) {
	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		slog.ErrorContext(ctx, "user id not found")
		return nil, huma.Error400BadRequest("user id not found")
	}

	profile, err := engine.Queries.GetUserProfile(ctx, int32(userID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to get user profile", "error", err)
		return nil, huma.Error500InternalServerError("failed to get user profile")
	}
	res := GetUserProfileDto{
		UserID:    userID,
		FullName:  profile.FullName,
		Username:  profile.Username,
		Bio:       profile.Bio,
		Stage:     profile.Stage,
		PhotoURL:  profile.PhotoUrl,
		CreatedAt: *profile.CreatedAt,
	}
	return SuccessResponse(res, 200), nil
}

type LoginDto struct {
	Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
}

func login(ctx context.Context, dto *LoginDto) (*Success[*TokensResponse], error) {
	user, err := engine.Queries.Login(ctx, dto.Body.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error401Unauthorized("invalid username")
		}
		slog.ErrorContext(ctx, "failed to login", "error", err)
		return nil, huma.Error500InternalServerError("failed to login")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Body.Password)); err != nil {
		slog.ErrorContext(ctx, "invalid password")
		return nil, huma.Error401Unauthorized("invalid password")
	}
	tokens, err := GenerateToken(int(user.ID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate token", "error", err)
		return nil, huma.Error500InternalServerError("failed to generate token")
	}
	return SuccessResponse(tokens, 200), nil
}

func InitUsersRoutes() {
	usersRoutes := huma.NewGroup(API, "/users")
	Handle(usersRoutes, http.MethodPost, "/register", register)
	Handle(usersRoutes, http.MethodPost, "/login", login)
	AuthHandle(usersRoutes, http.MethodPost, "/profile", CreateUserProfile)
	AuthHandle(usersRoutes, http.MethodGet, "/profile", GetUserProfile)
}
