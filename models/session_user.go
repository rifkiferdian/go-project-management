package models

// SessionUser is a minimal user payload stored in session cookies.
type SessionUser struct {
	UserID          int
	Name            string
	Email           string
	Initials        string
	Role            string
	IsAuthenticated bool
}
