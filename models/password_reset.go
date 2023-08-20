package models

import (
	"database/sql"
	"fmt"
	"time"
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a passwordReset is created
	Token     string
	Tokenhash string
	ExpiresAt time.Time
}

const (
	// DefaultResetDuration is the default time that a PasswordReset is
	// valid for.
	DefaultResetDuration = 1 * time.Hour
)

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each password reset token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: Complete the function")
}

// function to take an exisiting password reset token and to use it
// It will consume a token and return the user associated with it, or return an error if the token wasn't valid for any reason. Later when the password reset token is consumed ,we can call the userService to update the user
func (service *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO")
}
