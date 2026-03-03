package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTTokenIssuer_IssueAndVerify(t *testing.T) {
	issuer := NewJWTTokenIssuer("my-secret-key", 30, 60)

	token, refreshToken, err := issuer.Issue("u-1", "john")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty token")
	}

	if refreshToken == "" {
		t.Fatal("expected non-empty refresh token")
	}

	claims, err := issuer.Verify(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if claims == nil {
		t.Fatal("expected claims, got nil")
	}

	if claims.UserID != "u-1" {
		t.Fatalf("expected user id u-1, got %s", claims.UserID)
	}

	if claims.Username != "john" {
		t.Fatalf("expected username john, got %s", claims.Username)
	}
}

func TestJWTTokenIssuer_Verify_InvalidToken(t *testing.T) {
	issuer := NewJWTTokenIssuer("my-secret-key", 30, 60)

	_, err := issuer.Verify("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	if err != ErrorInvalidToken {
		t.Fatalf("expected ErrorInvalidToken, got %v", err)
	}
}

func TestJWTTokenIssuer_Verify_WrongSigningMethod(t *testing.T) {
	claims := jwt.MapClaims{"user_id": "u-1", "username": "john"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenStr, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	issuer := NewJWTTokenIssuer("secret", 30, 60)
	_, err = issuer.Verify(tokenStr)
	if err == nil {
		t.Fatal("expected error for wrong signing method")
	}
	if err != ErrorInvalidToken {
		t.Fatalf("expected ErrorInvalidToken, got %v", err)
	}
}

func TestJWTTokenIssuer_Verify_ExpiredToken(t *testing.T) {
	now := time.Now().UTC()
	expiredClaims := struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		jwt.RegisteredClaims
	}{
		UserID:   "u-1",
		Username: "john",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenStr, err := token.SignedString([]byte("my-secret-key"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	issuer := NewJWTTokenIssuer("my-secret-key", 30, 60)
	_, err = issuer.Verify(expiredTokenStr)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
	if err != ErrorInvalidToken {
		t.Fatalf("expected ErrorInvalidToken, got %v", err)
	}
}

func TestJWTTokenIssuer_Verify_RefreshToken(t *testing.T) {
	issuer := NewJWTTokenIssuer("my-secret-key", 30, 60)
	_, refreshToken, err := issuer.Issue("u-1", "john")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims, err := issuer.Verify(refreshToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.UserID != "u-1" || claims.Username != "john" {
		t.Fatalf("expected claims u-1/john, got %s/%s", claims.UserID, claims.Username)
	}
}

func TestNewJWTTokenIssuer(t *testing.T) {
	issuer := NewJWTTokenIssuer("secret", 30, 60)
	if issuer == nil {
		t.Fatal("expected non-nil issuer")
	}
}
