package main

import(
	"time"
	"github.com/google/uuid"
	"github.com/arjunsingh14/chirpy/internal/database"
)

type chirp struct {
	ID			uuid.UUID `json:"id"`
	User_id		uuid.UUID `json:"user_id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Email		string    `json:"email"`
	Body		string    `json:"body"`
}

func buildChirp(c database.Chirp) chirp {
	return chirp{
		ID: c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.CreatedAt,
		Body: c.Body,
		User_id: c.UserID,
	}
}