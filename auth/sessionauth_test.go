package auth

import (
	"net/http"
	"testing"

	"bitbucket.org/lewington/autoroller/assist"
	"bitbucket.org/lewington/autoroller/testhelp"
)

func TestSessionAuth(t *testing.T) {
	sessToken := "afdasd-gfsdr98043-fdsdf-sdf7sdf8"

	s := NewSessionAuth()
	s.Add(sessToken)

	req, err := http.NewRequest("POST", "localhost", nil)
	assist.Check(err)

	authenticated := s.IsAuthenticated(req)
	testhelp.ExpectFalse(t, authenticated, "expected request with no cookies not to be authenticated")

	req, err = http.NewRequest("Get", "localhost", nil)
	assist.Check(err)
	req.AddCookie(&http.Cookie{
		Name:  AccessCookieName,
		Value: "some random string that is gay",
	})
	authenticated = s.IsAuthenticated(req)
	testhelp.ExpectFalse(t, authenticated, "expected request with the wrong cookie value not to be authenticated")

	req, err = http.NewRequest("Get", "localhost", nil)
	assist.Check(err)
	req.AddCookie(&http.Cookie{
		Name:  AccessCookieName,
		Value: sessToken,
	})
	authenticated = s.IsAuthenticated(req)
	testhelp.ExpectTrue(t, authenticated, "expected request with the correct cookie key and value to be authenticated")

	req, err = http.NewRequest("Get", "localhost", nil)
	assist.Check(err)
	req.AddCookie(&http.Cookie{
		Name:  "wrong_cookie_name",
		Value: sessToken,
	})
	authenticated = s.IsAuthenticated(req)
	testhelp.ExpectFalse(t, authenticated, "expected request with the wrong cookie name but corretc value not to be authenticated")
}
