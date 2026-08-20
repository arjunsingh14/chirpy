package auth_test

import (
	"testing"
	"time"

	"github.com/arjunsingh14/chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := time.Hour
	issuedAfter := time.Now().Add(-time.Second)

	tokenString, err := auth.MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() returned an unexpected error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("MakeJWT() returned an empty token")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		t.Fatalf("parsing token returned by MakeJWT(): %v", err)
	}
	if !token.Valid {
		t.Fatal("MakeJWT() returned an invalid token")
	}
	if claims.Issuer != "chirpy-access" {
		t.Errorf("issuer = %q, want %q", claims.Issuer, "chirpy-access")
	}
	if claims.Subject != userID.String() {
		t.Errorf("subject = %q, want %q", claims.Subject, userID.String())
	}
	if claims.IssuedAt == nil {
		t.Fatal("issued-at claim is missing")
	}
	if claims.IssuedAt.Time.Before(issuedAfter) || claims.IssuedAt.Time.After(time.Now().Add(time.Second)) {
		t.Errorf("issued-at claim %v is outside the expected time window", claims.IssuedAt.Time)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("expiration claim is missing")
	}
	wantExpiration := time.Now().Add(expiresIn)
	if difference := claims.ExpiresAt.Time.Sub(wantExpiration); difference < -2*time.Second || difference > 2*time.Second {
		t.Errorf("expiration = %v, want approximately %v", claims.ExpiresAt.Time, wantExpiration)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	t.Run("valid token", func(t *testing.T) {
		tokenString, err := auth.MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT() returned an unexpected error: %v", err)
		}

		gotUserID, err := auth.ValidateJWT(tokenString, secret)
		if err != nil {
			t.Fatalf("ValidateJWT() returned an unexpected error: %v", err)
		}
		if gotUserID != userID {
			t.Errorf("ValidateJWT() user ID = %s, want %s", gotUserID, userID)
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		tokenString, err := auth.MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT() returned an unexpected error: %v", err)
		}

		gotUserID, err := auth.ValidateJWT(tokenString, "wrong-secret")
		if err == nil {
			t.Fatal("ValidateJWT() returned nil error for a token signed with a different secret")
		}
		if gotUserID != uuid.Nil {
			t.Errorf("ValidateJWT() user ID = %s, want uuid.Nil", gotUserID)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		tokenString, err := auth.MakeJWT(userID, secret, -time.Minute)
		if err != nil {
			t.Fatalf("MakeJWT() returned an unexpected error: %v", err)
		}

		gotUserID, err := auth.ValidateJWT(tokenString, secret)
		if err == nil {
			t.Fatal("ValidateJWT() returned nil error for an expired token")
		}
		if gotUserID != uuid.Nil {
			t.Errorf("ValidateJWT() user ID = %s, want uuid.Nil", gotUserID)
		}
	})

	t.Run("malformed token", func(t *testing.T) {
		gotUserID, err := auth.ValidateJWT("not-a-jwt", secret)
		if err == nil {
			t.Fatal("ValidateJWT() returned nil error for a malformed token")
		}
		if gotUserID != uuid.Nil {
			t.Errorf("ValidateJWT() user ID = %s, want uuid.Nil", gotUserID)
		}
	})

	t.Run("subject is not a UUID", func(t *testing.T) {
		claims := jwt.RegisteredClaims{
			Subject:   "not-a-uuid",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		}
		tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("signing test token: %v", err)
		}

		gotUserID, err := auth.ValidateJWT(tokenString, secret)
		if err == nil {
			t.Fatal("ValidateJWT() returned nil error for a non-UUID subject")
		}
		if gotUserID != uuid.Nil {
			t.Errorf("ValidateJWT() user ID = %s, want uuid.Nil", gotUserID)
		}
	})
}
