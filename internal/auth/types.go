package auth

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Role         UserRole  `json:"role"`
	PasswordHash string    `json:"-"` // Never serialize in API responses
	CreatedAt    time.Time `json:"created_at"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// CreateUserRequest represents a user creation request
type CreateUserRequest struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Role     UserRole `json:"role,omitempty"` // Optional, defaults to user
}

// CreateUserResponse represents a user creation response
type CreateUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// ListUsersResponse represents a user list response
type ListUsersResponse struct {
	Users []*User `json:"users"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
}

// ResetPasswordResponse represents a password reset response
type ResetPasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// NewUser creates a new user with a generated ID
func NewUser(username, passwordHash string, role UserRole) *User {
	if role == "" {
		role = RoleUser
	}
	return &User{
		ID:           uuid.New().String(),
		Username:     username,
		Role:         role,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// NewSession creates a new session for a user
func NewSession(userID string) *Session {
	return &Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour sessions
		CreatedAt: time.Now(),
	}
}

// IsExpired checks if a session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}