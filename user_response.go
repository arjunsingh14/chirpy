package main

import (
	"time"

	"github.com/arjunsingh14/chirpy/internal/database"
	"github.com/google/uuid"
)

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type loginResponse struct {
	userResponse
	Token string `json:"token"`
}

func newUserResponse(user database.User) userResponse {
	return userResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}
