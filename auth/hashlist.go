package auth

import (
	"bitbucket.org/lewington/autoroller/assist"
	"bitbucket.org/lewington/autoroller/globals"
	"golang.org/x/crypto/bcrypt"
)

// HashList contains a list of hashed passwords, and checks
// whether password you give it match any of the list entries.
type HashList struct {
	passwords [][]byte
}

// Add adds a hashed version of the given password to the list.
func (a *HashList) Add(password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	assist.Check(err)
	a.passwords = append(a.passwords, hashedPassword)
}

// IsAuthorized returns whether a hash representation of the given
// password already exists in the list.
func (a *HashList) IsAuthorized(password string) bool {
	for _, existingHash := range a.passwords {
		err := bcrypt.CompareHashAndPassword(existingHash, []byte(password))

		if err == nil {
			return true
		}

	}

	return false
}

// NewHashList initializes an empty HashList.
func NewHashList() *HashList {
	return &HashList{}
}

// DefaultHashList initializes a HashList that already contains the
// default system-wide password.
func DefaultHashList() *HashList {
	h := &HashList{}
	h.Add(globals.Password)
	return h
}
