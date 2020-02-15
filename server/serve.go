package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/lewington/betback/globals"
	"bitbucket.org/lewington/betback/serve/server"
	"bitbucket.org/lewington/erosai/assist"
	"bitbucket.org/lewington/erosai/database"
	"bitbucket.org/lewington/erosai/shared"
)

type Server struct {
	server.AuthServer
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

	var reg shared.Registration
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
	names := map[globals.BookieName]map[string]string{}

	nameBytes, err := json.Marshal(names)
	assist.Check(err)

	w.Write(nameBytes)
}

func New() *Server {
	return &Server{
		*server.NewAuthServer(globals.BetServerPort),
		database.NewArchivist(),
	}
}
