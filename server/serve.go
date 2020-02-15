package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/lewington/erosai/assist"
	"bitbucket.org/lewington/erosai/database"
	"bitbucket.org/lewington/erosai/globals"
	"bitbucket.org/lewington/erosai/shared"
)

type Server struct {
	AuthServer
	arch database.Archivist
}

func (s *Server) SpinUp() {
	s.SetAuthEndpoints()
	http.Handle("/", http.FileServer(http.Dir("public/static/")))

	http.HandleFunc("/register-attempt", s.register)
	http.HandleFunc("/details-data", s.IfAuthorized(s.usernames))

	fmt.Println("Server spun up and listening on port: ", globals.BetServerPort)

	if err := s.Srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}

	fmt.Println("Bet server closed")
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := assist.SafeBytes(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var reg shared.Details
	err = json.Unmarshal(reqBytes, &reg)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(reg)
	w.WriteHeader(http.StatusOK)

	if len(reg.Password) < 6 {
		w.Write(s.responseBytes("Password must be at least 6 characters"))
		return
	}

	if len(reg.Username) < 6 {

		w.Write(s.responseBytes("Username must be at least 6 characters"))
		return
	}

	if s.arch.DoesUserExist(reg) {
		w.Write(s.responseBytes("Username Already Taken"))
		return
	}

	s.arch.RegisterUser(reg)

	w.Write(s.responseBytes("Success"))

}

func (s *Server) responseBytes(message string) []byte {
	b, err := json.Marshal(Response{
		message,
	})
	assist.Check(err)
	return b
}

func (s *Server) usernames(w http.ResponseWriter, r *http.Request) {

}

func New() *Server {
	return &Server{
		*NewAuthServer(globals.BetServerPort),
		database.NewArchivist(),
	}
}
