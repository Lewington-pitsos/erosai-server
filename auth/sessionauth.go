package auth

import (
	"net/http"
)

// SessionAuth takes requests and determines if they hold
// valid session tokens.
type SessionAuth struct {
	SessionStore
}

// IsAuthenticated returns true iff the given request has
// a valid session token cookie.
func (s *SessionAuth) IsAuthenticated(r *http.Request) bool {
	sessionToken, err := r.Cookie(AccessCookieName)

	if err != nil {
		return false
	}

	return s.Contains(sessionToken.Value)
}

// NewSessionAuth initializes a SessionAuth with an underlying
// session store that contains the default access token.
func NewSessionAuth() *SessionAuth {
	return &SessionAuth{
		*DefaultSessionStore(),
	}
}
