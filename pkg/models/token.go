package models

import (
	"time"
)

// NewToken create to token model instance
func NewToken(jwtType TokenType) *Token {
	return &Token{
		Type: jwtType,
	}
}

// TokenType defines the use of generated token
type TokenType string

var (
	// TokenLogin will be used for login purposes
	TokenLogin TokenType = "Login"
	// TokenEmail will be used as authenticity for the link in emails
	TokenEmail TokenType = "Email"
)

// Token token model
type Token struct {
	Access          string        `bson:"Access"`
	AccessCreateAt  time.Time     `bson:"AccessCreateAt"`
	AccessExpiresIn time.Duration `bson:"AccessExpiresIn"`
	Type            TokenType     `bson:"type"`
}

// GetAccess access Token
func (t *Token) GetAccess() string {
	return t.Access
}

// SetAccess access Token
func (t *Token) SetAccess(access string) {
	t.Access = access
}

// GetAccessCreateAt create Time
func (t *Token) GetAccessCreateAt() time.Time {
	return t.AccessCreateAt
}

// SetAccessCreateAt create Time
func (t *Token) SetAccessCreateAt(createAt time.Time) {
	t.AccessCreateAt = createAt
}

// GetAccessExpiresIn the lifetime in seconds of the access token
func (t *Token) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresIn
}

// SetAccessExpiresIn the lifetime in seconds of the access token
func (t *Token) SetAccessExpiresIn(exp time.Duration) {
	t.AccessExpiresIn = exp
}
