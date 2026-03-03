package security

import (
	"authentication-project-exam/internal/core/port/outbound"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrorInvalidToken = errors.New("invalid token")

type JWTTokenIssuer struct {
	secret               string
	expireMinutes        int
	refreshExpireMinutes int
}

type jwtClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTTokenIssuer(secret string, expireMinutes int, refreshExpireMinutes int) outbound.TokenManager {
	return &JWTTokenIssuer{
		secret:               secret,
		expireMinutes:        expireMinutes,
		refreshExpireMinutes: refreshExpireMinutes,
	}
}

func (j *JWTTokenIssuer) Issue(userId, username string) (string, string, error) {
	now := time.Now().UTC()

	accessClaims := jwtClaims{
		UserID:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.expireMinutes) * time.Minute)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}
	refreshClaims := jwtClaims{
		UserID:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.refreshExpireMinutes) * time.Minute)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}
	return accessTokenStr, refreshTokenStr, nil
}

func (j *JWTTokenIssuer) Verify(tokenStr string) (*outbound.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrorInvalidToken
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, ErrorInvalidToken
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, ErrorInvalidToken
	}

	return &outbound.TokenClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
	}, nil
}
