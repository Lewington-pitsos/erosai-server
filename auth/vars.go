package auth

import (
	"net/http"
	"time"

	"bitbucket.org/lewington/autoroller/globals"
)

const AccessCookieName = "erosai_access_token"

var DefaultCookie = &http.Cookie{
	Name:    AccessCookieName,
	Value:   globals.DefaultAccessToken,
	Expires: time.Now().Add(SessionDuration),
}
