package auth

import (
	"sync"
	"time"

	"bitbucket.org/lewington/autoroller/globals"

	"bitbucket.org/lewington/autoroller/assist"
)

var SessionDuration = time.Hour * 24 * 300

// SessionStore is an async-safe store of session tokens.
// It can tell you if a given string matches a currently
// stored token and removes tokens some time after they
// were added.
type SessionStore struct {
	mutex  *sync.Mutex
	tokens map[string]time.Time
}

// Add adds the given token to the stored list. It also
// sets an expiration date for that token. Adding an existing
// token refreshes it's expiration date.
func (s *SessionStore) Add(token string) {
	s.mutex.Lock()
	s.tokens[token] = assist.Timestamp().Add(SessionDuration)
	s.mutex.Unlock()
}

// Contains checks whether there is a stored token matching the
// given string. If so, the token is removed if it has expired.
// If the token exists and hasn't expired true is returned.
func (s *SessionStore) Contains(token string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	expirey, exists := s.tokens[token]

	if exists {
		if assist.Timestamp().After(expirey) {
			delete(s.tokens, token)
			return false
		}

		return true
	}

	return false
}

// NewSessionStore initializes an empty SessionStore.
func NewSessionStore() *SessionStore {
	return &SessionStore{
		&sync.Mutex{},
		map[string]time.Time{},
	}
}

// DefaultSessionStore initializes a SessionStore that already
// holds the default access token that never expires.
func DefaultSessionStore() *SessionStore {
	return &SessionStore{
		&sync.Mutex{},
		map[string]time.Time{
			globals.DefaultAccessToken: assist.Timestamp().Add(time.Hour * 100000),
		},
	}
}
