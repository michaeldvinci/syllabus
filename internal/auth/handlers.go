package auth

import (
	"encoding/json"
	"html/template"
	"net/http"
)

// AuthHandlers provides authentication-related HTTP handlers
type AuthHandlers struct {
	store *Store
}

// NewAuthHandlers creates new authentication handlers
func NewAuthHandlers(store *Store) *AuthHandlers {
	return &AuthHandlers{store: store}
}

// HandleLogin serves the login page and handles login POST requests
func (h *AuthHandlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.serveLoginPage(w, r)
	case "POST":
		h.handleLoginPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// serveLoginPage serves the login HTML page
func (h *AuthHandlers) serveLoginPage(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	if cookie, err := r.Cookie("session_token"); err == nil {
		if _, err := h.store.GetUserBySession(cookie.Value); err == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	tmpl := template.Must(template.New("login").Parse(LoginHTML))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleLoginPost processes login form submissions
func (h *AuthHandlers) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Handle both JSON and form data
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		req.Username = r.FormValue("username")
		req.Password = r.FormValue("password")
	}

	// Authenticate user
	user, err := h.store.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		response := LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}

		if contentType == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Redirect back to login with error
		http.Redirect(w, r, "/login?error=invalid", http.StatusSeeOther)
		return
	}

	// Create session
	session, err := h.store.CreateSession(user.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
	})

	response := LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   session.Token,
		User:    user,
	}

	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleLogout handles logout requests
func (h *AuthHandlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Delete session from store
		h.store.DeleteSession(cookie.Value)
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
	})

	// Check if this is an API request
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Logout successful",
		})
		return
	}

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// HandleAPI returns current user info as JSON
func (h *AuthHandlers) HandleAPI(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user":          user,
	})
}

// HandleCreateUser handles user creation requests (admin only)
func (h *AuthHandlers) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateUserResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateUserResponse{
			Success: false,
			Message: "Username and password are required",
		})
		return
	}

	// Default role is user if not specified
	if req.Role == "" {
		req.Role = RoleUser
	}

	// Create user
	user, err := h.store.CreateUserWithRole(req.Username, req.Password, req.Role)
	if err != nil {
		var statusCode int
		var message string

		switch err {
		case ErrUserExists:
			statusCode = http.StatusConflict
			message = "Username already exists"
		default:
			statusCode = http.StatusInternalServerError
			message = "Failed to create user"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(CreateUserResponse{
			Success: false,
			Message: message,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateUserResponse{
		Success: true,
		Message: "User created successfully",
		User:    user,
	})
}

// HandleListUsers handles user listing requests (admin only)
func (h *AuthHandlers) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users := h.store.ListUsers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListUsersResponse{
		Users: users,
	})
}

// HandleResetPassword handles password reset requests (admin only)
func (h *AuthHandlers) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResetPasswordResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate input
	if req.Username == "" || req.NewPassword == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResetPasswordResponse{
			Success: false,
			Message: "Username and new password are required",
		})
		return
	}

	// Reset password
	err := h.store.ResetUserPassword(req.Username, req.NewPassword)
	if err != nil {
		var statusCode int
		var message string

		switch err {
		case ErrUserNotFound:
			statusCode = http.StatusNotFound
			message = "User not found"
		default:
			statusCode = http.StatusInternalServerError
			message = "Failed to reset password"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(ResetPasswordResponse{
			Success: false,
			Message: message,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResetPasswordResponse{
		Success: true,
		Message: "Password reset successfully",
	})
}

// HandleDeleteUser handles user deletion requests (admin only)
func (h *AuthHandlers) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}

	// Validate input
	if req.Username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Username is required",
		})
		return
	}

	// Prevent deletion of admin user
	if req.Username == "admin" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Cannot delete admin user",
		})
		return
	}

	err := h.store.DeleteUser(req.Username)
	if err != nil {
		var statusCode int
		var message string

		switch err {
		case ErrUserNotFound:
			statusCode = http.StatusNotFound
			message = "User not found"
		default:
			statusCode = http.StatusInternalServerError
			message = "Failed to delete user"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": message,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}

const LoginHTML = `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Syllabus - Login</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
:root {
  --bg: #ffffff;
  --text: #111827;
  --muted: #6b7280;
  --line: #e5e7eb;
  --primary: #3b82f6;
  --primary-hover: #2563eb;
}

[data-theme="dark"] {
  --bg: #1f2937;
  --text: #f9fafb;
  --muted: #9ca3af;
  --line: #374151;
  --primary: #3b82f6;
  --primary-hover: #2563eb;
}

* {
  box-sizing: border-box;
}

body {
  font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial;
  margin: 0;
  padding: 0;
  background: var(--bg);
  color: var(--text);
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-container {
  background: var(--bg);
  border: 1px solid var(--line);
  border-radius: 0.5rem;
  padding: 2rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.logo {
  text-align: center;
  margin-bottom: 2rem;
}

.logo h1 {
  margin: 0;
  font-size: 1.875rem;
  font-weight: 700;
  color: var(--text);
}

.form-group {
  margin-bottom: 1rem;
}

label {
  display: block;
  margin-bottom: 0.25rem;
  font-weight: 500;
  color: var(--text);
}

input[type="text"],
input[type="password"] {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--line);
  border-radius: 0.375rem;
  background: var(--bg);
  color: var(--text);
  font-size: 1rem;
}

input[type="text"]:focus,
input[type="password"]:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.btn {
  width: 100%;
  padding: 0.75rem;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 0.375rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn:hover {
  background: var(--primary-hover);
}

.error {
  color: #dc2626;
  font-size: 0.875rem;
  margin-top: 0.5rem;
  text-align: center;
}

.footer {
  margin-top: 2rem;
  text-align: center;
  color: var(--muted);
  font-size: 0.875rem;
}
</style>
</head>
<body>
  <div class="login-container">
    <div class="logo">
      <img src="/static/syllabus_logo.png" alt="Syllabus Logo" style="height: 3rem; width: auto; margin-bottom: 0.5rem;">
      <h1>Syllabus</h1>
    </div>
    
    <form method="POST" action="/login">
      <div class="form-group">
        <label for="username">Username</label>
        <input type="text" id="username" name="username" required>
      </div>
      
      <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" name="password" required>
      </div>
      
      <button type="submit" class="btn">Sign In</button>
      
      <script>
        // Check for error in URL params
        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.get('error') === 'invalid') {
          const errorDiv = document.createElement('div');
          errorDiv.className = 'error';
          errorDiv.textContent = 'Invalid username or password';
          document.querySelector('form').appendChild(errorDiv);
        }
      </script>
    </form>
    
    <div class="footer">
    </div>
  </div>
</body>
</html>
`
