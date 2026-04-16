package domain

import "time"

// User represents the user domain entity.
type User struct {
	ID           uint64
	Email        string
	Username     string
	PasswordHash string
	AvatarURL    string
	OAuthProvider string
	OAuthID      string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

const (
	OAuthProviderNone   = "none"
	OAuthProviderGitHub = "github"
	OAuthProviderGoogle = "google"
	RoleUser            = "user"
	RoleAdmin           = "admin"
)
