package auth

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrSessionNotFound  = errors.New("session not found")
	ErrSessionExpired   = errors.New("session expired")
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidPassword  = errors.New("invalid password")
)

// Store handles user and session storage
type Store struct {
	users      map[string]*User    // username -> user
	sessions   map[string]*Session // token -> session
	mu         sync.RWMutex
	dataFile   string              // path to persistence file
	autoSave   bool                // whether to auto-save changes
}

// NewStore creates a new authentication store
func NewStore() *Store {
	return NewStoreWithFile("")
}

// NewStoreWithFile creates a new authentication store with file persistence
func NewStoreWithFile(dataFile string) *Store {
	store := &Store{
		users:    make(map[string]*User),
		sessions: make(map[string]*Session),
		dataFile: dataFile,
		autoSave: dataFile != "",
	}
	
	// Load existing data if file is specified
	if dataFile != "" {
		if err := store.LoadFromFile(dataFile); err != nil {
			// Log error but don't fail - we can continue with empty store
			// In a production app, you might want to handle this differently
		}
	}
	
	// Start cleanup goroutine for expired sessions
	go store.cleanupExpiredSessions()
	
	return store
}

// save persists the store to file if auto-save is enabled
func (s *Store) save() {
	if s.autoSave && s.dataFile != "" {
		// Save asynchronously to avoid blocking
		go func() {
			if err := s.SaveToFile(s.dataFile); err != nil {
				// In production, you'd want to log this properly
			}
		}()
	}
}

// CreateUser creates a new user
func (s *Store) CreateUser(username, password string) (*User, error) {
	return s.CreateUserWithRole(username, password, RoleUser)
}

// CreateUserWithRole creates a new user with a specific role
func (s *Store) CreateUserWithRole(username, password string, role UserRole) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Check if user already exists
	if _, exists := s.users[username]; exists {
		return nil, ErrUserExists
	}
	
	// Hash password
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	
	// Create user
	user := NewUser(username, hash, role)
	s.users[username] = user
	
	// Save to file if persistence is enabled
	s.save()
	
	return user, nil
}

// AuthenticateUser authenticates a user with username/password
func (s *Store) AuthenticateUser(username, password string) (*User, error) {
	s.mu.RLock()
	user, exists := s.users[username]
	s.mu.RUnlock()
	
	if !exists {
		return nil, ErrUserNotFound
	}
	
	if !VerifyPassword(password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}
	
	return user, nil
}

// CreateSession creates a new session for a user
func (s *Store) CreateSession(userID string) (*Session, error) {
	session := NewSession(userID)
	
	s.mu.Lock()
	s.sessions[session.Token] = session
	s.mu.Unlock()
	
	return session, nil
}

// GetSession retrieves a session by token
func (s *Store) GetSession(token string) (*Session, error) {
	s.mu.RLock()
	session, exists := s.sessions[token]
	s.mu.RUnlock()
	
	if !exists {
		return nil, ErrSessionNotFound
	}
	
	if session.IsExpired() {
		s.DeleteSession(token)
		return nil, ErrSessionExpired
	}
	
	return session, nil
}

// GetUserBySession gets a user by session token
func (s *Store) GetUserBySession(token string) (*User, error) {
	session, err := s.GetSession(token)
	if err != nil {
		return nil, err
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, user := range s.users {
		if user.ID == session.UserID {
			return user, nil
		}
	}
	
	return nil, ErrUserNotFound
}

// DeleteSession deletes a session
func (s *Store) DeleteSession(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

// GetUser gets a user by username
func (s *Store) GetUser(username string) (*User, error) {
	s.mu.RLock()
	user, exists := s.users[username]
	s.mu.RUnlock()
	
	if !exists {
		return nil, ErrUserNotFound
	}
	
	return user, nil
}

// ListUsers returns all users (admin only)
func (s *Store) ListUsers() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	
	return users
}

// ResetUserPassword resets a user's password (admin only)
func (s *Store) ResetUserPassword(username, newPassword string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Check if user exists
	user, exists := s.users[username]
	if !exists {
		return ErrUserNotFound
	}
	
	// Hash new password
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	
	// Update password
	user.PasswordHash = hash
	
	// Save to file if persistence is enabled
	s.save()
	
	return nil
}

// DeleteUser deletes a user by username
func (s *Store) DeleteUser(username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Check if user exists
	if _, exists := s.users[username]; !exists {
		return ErrUserNotFound
	}
	
	// Delete user
	delete(s.users, username)
	
	// Save to file if persistence is enabled
	s.save()
	
	return nil
}

// cleanupExpiredSessions runs periodically to clean up expired sessions
func (s *Store) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		s.mu.Lock()
		for token, session := range s.sessions {
			if session.IsExpired() {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}