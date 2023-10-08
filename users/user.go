package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" field:"id"`
	CreatedAt time.Time `json:"created_at" field:"created_at"`
	UpdatedAt time.Time `json:"updated_at" field:"updated_at"`
	Name      string    `json:"name" field:"name"`
}
