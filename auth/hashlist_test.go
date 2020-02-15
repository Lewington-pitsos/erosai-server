package auth

import (
	"testing"

	"bitbucket.org/lewington/autoroller/globals"

	"bitbucket.org/lewington/autoroller/testhelp"
)

func TestHashList(t *testing.T) {
	a := NewHashList()
	pass1 := "1m@N3b$V5c"
	a.Add(pass1)

	if string(a.passwords[0]) == pass1 {
		t.Fatalf("password to be stored in a secure hash got %v", pass1)
	}

	passes := a.IsAuthorized(pass1)

	testhelp.ExpectTrue(t, passes, "the same password to be authorized")

	passes = a.IsAuthorized("randomstring")
	testhelp.ExpectFalse(t, passes, "new random password not to be authorized")

	str2 := "Zasjdoi"

	passes = a.IsAuthorized(str2)
	testhelp.ExpectFalse(t, passes, "new random password not to be authorized")

	a.Add(str2)

	passes = a.IsAuthorized(str2)
	testhelp.ExpectTrue(t, passes, "added password to be authorized")
}

func TestDefaultHashList(t *testing.T) {
	a := DefaultHashList()

	passes := a.IsAuthorized(globals.Password)

	testhelp.ExpectTrue(t, passes, "the default password to be authorized")
}
