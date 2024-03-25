package user

import (
	"time"
)

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,min=7,max=50"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type AuthenticateRequest struct {
	Email    string `json:"email" validate:"required,email,min=7,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
