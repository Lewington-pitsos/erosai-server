package auth

import (
	"net/http"

	"bitbucket.org/lewington/erosai/database"
	"bitbucket.org/lewington/erosai/shared"
)

// SessionAuth takes requests and determines if they hold
// valid session tokens.
type SessionAuth struct {
	arch database.Archivist
}

// IsAuthenticated returns true iff the given request has
// a valid session token cookie.
func (s *SessionAuth) IsAuthenticated(r *http.Request) (*shared.Details, bool) {
	sessionToken, err := r.Cookie(AccessCookieName)

	if err != nil {
		return nil, false
	}

	t := s.arch.UserForToken(sessionToken.Value)
	return &t, true
}

// NewSessionAuth initializes a SessionAuth with an underlying
// session store that contains the default access token.
func NewSessionAuth() *SessionAuth {
	return &SessionAuth{
		database.NewArchivist(),
	}
}
