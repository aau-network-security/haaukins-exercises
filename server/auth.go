package server

import (
	"context"
	"errors"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
)

const (
	AUTH_KEY = "au"
)

var (
	ErrInvalidAuthKey     = errors.New("invalid Authentication Key")
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrMissingKey         = errors.New("no Authentication Key provided")
)

type Authenticator interface {
	AuthenticateContext(context.Context) error
}

type auth struct {
	sKey string // Signin Key
	aKey string // Auth Key
}

func NewAuthenticator(Skey, AKey string) Authenticator {
	return &auth{sKey: Skey, aKey: AKey}
}

func (a *auth) AuthenticateContext(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ErrMissingKey
	}

	if len(md["token"]) == 0 {
		return ErrMissingKey
	}

	token := md["token"][0]
	if token == "" {
		return ErrMissingKey
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return ctx, ErrInvalidTokenFormat
		}

		return []byte(a.sKey), nil
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to parse token")
		return ErrInvalidTokenFormat
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return ErrInvalidTokenFormat
	}

	authKey, ok := claims[AUTH_KEY].(string)
	if !ok {
		return ErrInvalidTokenFormat
	}

	if authKey != a.aKey {
		return ErrInvalidAuthKey
	}

	return nil
}
