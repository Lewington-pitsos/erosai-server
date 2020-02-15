package server

import (
	"net/http"
)

// AuthServer contains a http server and a session
// authenticator, and serves a login endpoint where
// a user with the dws password can create authenticated
// sessions.
type AuthServer struct {
	Srv *http.Server
	sessManager
}

func (s *AuthServer) Login(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/static/html/login.html")
}

func (s *AuthServer) SetAuthEndpoints() {
	http.HandleFunc("/login-attempt", s.Authenticate)
}

func (s *AuthServer) IfAuthorized(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, isAuthenticated := s.IsAuthenticated(r)
		if isAuthenticated {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		}
	}
}

// NewAuthServer initializes an AuthServer.
func NewAuthServer(port string) *AuthServer {
	return &AuthServer{
		&http.Server{Addr: port},
		*newSessManager(),
	}
}
