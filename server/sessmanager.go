package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/lewington/erosai-server/assist"
	"bitbucket.org/lewington/erosai-server/auth"
	"bitbucket.org/lewington/erosai-server/database"
	"bitbucket.org/lewington/erosai-server/lg"
	"bitbucket.org/lewington/erosai-server/shared"
	"github.com/google/uuid"
)

type sessManager struct {
	auth.SessionAuth
	validPasswords *auth.HashList
	arch           database.Archivist
}

func (s *sessManager) Authenticate(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := assist.SafeBytes(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var details shared.Details

	err = json.Unmarshal(reqBytes, &details)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if s.arch.DoesUserExist(details) {
		w.WriteHeader(http.StatusOK)
		sessionToken := uuid.New().String()
		s.arch.SetUserToken(details, sessionToken)

		cookie := &http.Cookie{
			Name:    auth.AccessCookieName,
			Value:   sessionToken,
			Expires: time.Now().Add(auth.SessionDuration),
		}
		http.SetCookie(w, cookie)
		fmt.Println(sessionToken)
		w.Write([]byte(sessionToken))
		return
	}

	lg.L.Debug("failed login attempt with password #%v from address %v", details.Password, r.RemoteAddr)

	w.WriteHeader(http.StatusUnauthorized)
}

func newSessManager() *sessManager {
	return &sessManager{
		*auth.NewSessionAuth(),
		auth.DefaultHashList(),
		database.NewArchivist(),
	}
}
