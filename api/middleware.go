package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/AliAlbhrani/StudentsArchive/env"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID int    `json:"user_id"`
	Scope  string `json:"scope"`
	*jwt.RegisteredClaims
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

const (
	ScopeAccessToken  = "ACCESS_TOKEN"
	ScopeRefreshToken = "REFRESH_TOKEN"
	userIDKey         = "userID"
)

// GenerateToken generates a JWT token for the given user ID.
func GenerateToken(userID int) (*TokensResponse, error) {
	claims := TokenClaims{
		UserID: userID,
		Scope:  ScopeAccessToken,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 10000 * time.Hour)),
			Issuer:    env.ISSUER,
			Audience:  jwt.ClaimStrings{env.ISSUER},
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(env.SECRET_KEY))
	if err != nil {
		return nil, err
	}
	refClaims := TokenClaims{
		UserID: userID,
		Scope:  ScopeRefreshToken,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    env.ISSUER,
			Audience:  jwt.ClaimStrings{env.ISSUER},
		},
	}
	refToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims).SignedString([]byte(env.SECRET_KEY))
	if err != nil {
		return nil, err
	}
	return &TokensResponse{
		AccessToken:  token,
		RefreshToken: refToken,
	}, nil
}

func AuthMiddleware(ctx huma.Context, next func(huma.Context)) {
	tokenStr := ctx.Header("Authorization")
	if tokenStr == "" {
		huma.WriteErr(API, ctx, http.StatusUnauthorized, "missing token", nil)
		return
	}
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}
	claims := TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		return []byte(env.SECRET_KEY), nil
	})
	if err != nil {
		slog.Error("invalid token", "error", err)
		huma.WriteErr(API, ctx, http.StatusUnauthorized, "invalid token", nil)
		return
	}
	if !token.Valid {
		huma.WriteErr(API, ctx, http.StatusUnauthorized, "invalid token", nil)
		return
	}
	if claims.Scope != ScopeAccessToken {
		huma.WriteErr(API, ctx, http.StatusUnauthorized, "invalid token", nil)
		return
	}
	if claims.UserID == 0 {
		huma.WriteErr(API, ctx, http.StatusUnauthorized, "invalid token", nil)
		return
	}
	ctx = huma.WithValue(ctx, userIDKey, claims.UserID)
	next(ctx)
}

func GetUserID(ctx context.Context) int {
	if userID, ok := ctx.Value(userIDKey).(int); ok {
		return userID
	}
	return 0
}
