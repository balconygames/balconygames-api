package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"

	pkghttp "gitlab.com/balconygames/analytics/pkg/http"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

// ContextKey naming for context key type
type ContextKey string

// ContextUserKey should have parsed user by JWT token.
const ContextUserKey ContextKey = "user"
const ContextSignerKey ContextKey = "signer"

func NewJWTHttpVerifierMiddleware(secret string) func(handler http.Handler) http.Handler {
	auth := jwtauth.New(jwt.SigningMethodHS256.Name, []byte(secret), nil)
	return jwtauth.Verifier(auth)
}

// NewJWTUserMiddleware should extract user from JWT token and pass it via context
func NewJWTUserMiddleware(secret string) func(next http.Handler) http.Handler {
	tokenSigner := NewSigner(secret)

	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := DecodeJWT(r, tokenSigner)
			if err != nil {
				pkghttp.Error(w, errors.Wrap(err, "invalid jwt token"))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return fn
}

func NewTokenSignerMiddleware(secret string) func(next http.Handler) http.Handler {
	tokenSigner := NewSigner(secret)

	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextSignerKey, tokenSigner)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return fn
}

// DecodeJWT returns parsed user information using JWT token from
// Authorization header.
func DecodeJWT(req *http.Request, tokenSigner *Signer) (*UserInfo, error) {
	token := req.Header.Get("Authorization")

	if token == "" {
		return nil, errors.New("missing authorization header")
	}

	if len(token) > 7 && strings.ToUpper(token[0:6]) == "BEARER" {
		token = token[7:]
	}

	claims, err := tokenSigner.Decode(token)
	if err != nil {
		return nil, errors.Wrap(err, "invalid jwt token")
	}
	user := claims.UserInfo

	return &user, nil
}

// GetUser retrieve user information from context of request.
func GetUser(req *http.Request) *UserInfo {
	user, _ := req.Context().Value(ContextUserKey).(*UserInfo)
	return user
}

func GetScope(req *http.Request) *sharedmodels.Scope {
	user := GetUser(req)

	return &sharedmodels.Scope{
		GameID: user.GameID,
		AppID:  user.AppID,
		UserID: user.UserID,
	}
}

// GetJWTUser retrieve user information from context of request.
func GetSigner(req *http.Request) *Signer {
	user, _ := req.Context().Value(ContextSignerKey).(*Signer)
	return user
}
