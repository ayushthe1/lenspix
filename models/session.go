package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/ayushthe1/lenspix/rand"
)

const (
	// the minimum number of bytes to be used for each session token
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserId int
	// Token is only set when creating a new session. When we look up a session this will be left empty, as we only store the hash of a session token in our database and we can't reverse it into a raw token.
	Token     string
	Tokenhash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
}

// When the user sign up/in ,we're going to pass in the userID & create a session and then use the token that was generated to set the cookie.
// Later on when the user comes back to the application, we're going to look up that token value and check if it matches the Tokenhash in our db.

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserId:    userID,
		Token:     token,
		Tokenhash: ss.hash(token),
	}

	// Try to update the user's session. If err, create a new session
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;`, session.UserId, session.Tokenhash)

	err = row.Scan(&session.ID)

	if err == sql.ErrNoRows {
		// Store the token in our DB
		row := ss.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2)
			RETURNING id;`, session.UserId, session.Tokenhash)

		err = row.Scan(&session.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

// User takes the token from the cookie and returns the User
func (ss *SessionService) User(token string) (*User, error) {

	// Hash the session tokens
	tokenhash := ss.hash(token)

	// Query for the session with that hash
	var user User
	row := ss.DB.QueryRow(`
		SELECT users.id,
    	users.email,
    	users.password_hash
	FROM sessions
    	JOIN users ON users.id = sessions.user_id
	WHERE sessions.token_hash = $1;`, tokenhash)

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	// Return the user
	return &user, nil
}

// function to delete the session (signout)
func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`
	DELETE FROM sessions
	WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

// function for hashing the session token using SHA256
func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
