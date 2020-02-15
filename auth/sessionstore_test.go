package auth

import (
	"testing"

	"bitbucket.org/lewington/autoroller/assist"
	"bitbucket.org/lewington/autoroller/testhelp"
)

func TestSessionStore(t *testing.T) {
	s := NewSessionStore()
	sess1 := "asdasdasd"

	testhelp.ExpectFalse(t, s.Contains(sess1), "expected empty session not to have any token")

	s.Add(sess1)

	testhelp.ExpectTrue(t, s.Contains(sess1), "expected session to contain recently added token")

	sess2 := "8r7yw8ruew8"

	testhelp.ExpectFalse(t, s.Contains(sess2), "session not to have a particular token stored")

	s.Add(sess2)

	testhelp.ExpectTrue(t, s.Contains(sess2), "expected session to contain recently added token")
	testhelp.ExpectTrue(t, s.Contains(sess1), "expected session to contain recently added token")

	s.tokens[sess2] = assist.Timestamp()
	testhelp.ExpectFalse(t, s.Contains(sess2), "expected expired token not to be stored still")
	testhelp.ExpectFalse(t, s.Contains(sess2), "expected expired token not to be stored still")
	testhelp.ExpectTrue(t, s.Contains(sess1), "expected session to contain recently added token")

	s.Add(sess2)
	s.tokens[sess2] = assist.Timestamp()
	s.Add(sess2)
	testhelp.ExpectTrue(t, s.Contains(sess2), "expected session to contain recently refreshed token")

}
