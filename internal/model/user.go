package model

import (
	"time"
)

// UserRole type for user roles
type UserRole string

const (
	// RoleClient represents a regular client user
	RoleClient UserRole = "client"
	// RoleModerator represents a moderator user
	RoleModerator UserRole = "moderator"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Password is not returned in JSON
	Role      UserRole  `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
} 