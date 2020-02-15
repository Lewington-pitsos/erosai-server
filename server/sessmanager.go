package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"bitbucket.org/lewington/erosai/assist"
	"bitbucket.org/lewington/erosai/auth"
	"bitbucket.org/lewington/erosai/lg"
)

type sessManager struct {
	auth.SessionAuth
	validPasswords *auth.HashList
	throttle       auth.Throttle
}

func (s *sessManager) Authenticate(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := assist.SafeBytes(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var creds Credentials

	err = json.Unmarshal(reqBytes, &creds)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if s.throttle.Allow() && s.validPasswords.IsAuthorized(creds.Password) {
		s.throttle.Succeed()
		w.WriteHeader(http.StatusOK)
		sessionToken := uuid.New().String()
		s.Add(sessionToken)
		cookie := &http.Cookie{
			Name:    auth.AccessCookieName,
			Value:   sessionToken,
			Expires: time.Now().Add(auth.SessionDuration),
		}
		http.SetCookie(w, cookie)
		w.Write([]byte(cookie.String()))
		return
	}

	lg.L.Debug("failed login attempt with password #%v from address %v", creds.Password, r.RemoteAddr)

	s.throttle.Fail()
	w.WriteHeader(http.StatusUnauthorized)
}

func newSessManager() *sessManager {
	return &sessManager{
		*auth.NewSessionAuth(),
		auth.DefaultHashList(),
		auth.NewTimeThrottle(4, time.Minute*30),
	}
}
