package shared

import (
	"bitbucket.org/lewington/erosai/assist"
	"golang.org/x/crypto/bcrypt"
)

type Details struct {
	Username string
	Password string
}

func (d *Details) IsNull() bool {
	return d.Username == "" && d.Password == ""
}

func (d *Details) Combination() string {
	return d.Password + d.Username
}

func (d *Details) Hash() string {
	hash, err := bcrypt.GenerateFromPassword([]byte(d.Combination()), 9)

	assist.Check(err)

	return string(hash)
}
