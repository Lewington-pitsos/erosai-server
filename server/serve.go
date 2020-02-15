package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/lewington/betback/globals"
	"bitbucket.org/lewington/betback/serve/server"
	"bitbucket.org/lewington/erosai/assist"
)

type Server struct {
	server.AuthServer
}

func (s *Server) SpinUp() {
	s.SetAuthEndpoints()
	http.Handle("/", http.FileServer(http.Dir("public/static/")))

	http.HandleFunc("/details-data", s.IfAuthorized(s.usernames))

	fmt.Println("Server spun up and listening on port: ", globals.BetServerPort)

	if err := s.Srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}

	fmt.Println("Bet server closed")
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
	}
}
