package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PersistentUser represents a user for file storage (includes password hash)
type PersistentUser struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Role         UserRole  `json:"role"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

// PersistentData represents the data structure for file storage
type PersistentData struct {
	Users map[string]*PersistentUser `json:"users"`
}

// SaveToFile saves users to a JSON file
func (s *Store) SaveToFile(filename string) error {
	s.mu.RLock()
	data := PersistentData{
		Users: make(map[string]*PersistentUser),
	}
	
	// Copy users (excluding sessions as they should be temporary)
	for username, user := range s.users {
		data.Users[username] = &PersistentUser{
			ID:           user.ID,
			Username:     user.Username,
			Role:         user.Role,
			PasswordHash: user.PasswordHash,
			CreatedAt:    user.CreatedAt,
		}
	}
	s.mu.RUnlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// LoadFromFile loads users from a JSON file
func (s *Store) LoadFromFile(filename string) error {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist, that's okay for first run
		return nil
	}

	// Read file
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON
	var data PersistentData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	// Load users into store
	s.mu.Lock()
	if data.Users != nil {
		s.users = make(map[string]*User)
		for username, persistentUser := range data.Users {
			s.users[username] = &User{
				ID:           persistentUser.ID,
				Username:     persistentUser.Username,
				Role:         persistentUser.Role,
				PasswordHash: persistentUser.PasswordHash,
				CreatedAt:    persistentUser.CreatedAt,
			}
		}
	}
	s.mu.Unlock()

	return nil
}