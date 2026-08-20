package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/arjunsingh14/chirpy/internal/database"
	"github.com/google/uuid"
)

func TestLoginResponseMarshalsUserFieldsAndTokenAtTopLevel(t *testing.T) {
	userID := uuid.New()
	createdAt := time.Date(2026, time.August, 21, 10, 0, 0, 0, time.UTC)
	databaseUser := database.User{
		ID:        userID,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		Email:     "user@example.com",
	}
	response := loginResponse{
		userResponse: newUserResponse(databaseUser),
		Token:        "signed-token",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal login response: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal login response: %v", err)
	}

	if payload["id"] != userID.String() {
		t.Errorf("id = %v, want %s", payload["id"], userID)
	}
	if payload["email"] != databaseUser.Email {
		t.Errorf("email = %v, want %s", payload["email"], databaseUser.Email)
	}
	if payload["token"] != response.Token {
		t.Errorf("token = %v, want %s", payload["token"], response.Token)
	}
	if _, nested := payload["userResponse"]; nested {
		t.Error("user fields were nested instead of marshaled at the top level")
	}
}
