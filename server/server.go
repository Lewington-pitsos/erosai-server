package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/lewington/erosai/assist"
	"bitbucket.org/lewington/erosai/database"
	"bitbucket.org/lewington/erosai/globals"
	"bitbucket.org/lewington/erosai/lg"
	"bitbucket.org/lewington/erosai/scanner"
	"bitbucket.org/lewington/erosai/shared"
)

type Server struct {
	AuthServer
	arch        database.Archivist
	newURLInput chan shared.Link
	scan        *scanner.Scanner
}

func (s *Server) SpinUp() {
	s.SetAuthEndpoints()
	http.Handle("/", http.FileServer(http.Dir("public/static/")))

	http.HandleFunc("/register-attempt", s.register)
	http.HandleFunc("/process-url", s.processURL)

	http.HandleFunc("/get-recommendations", s.recommendations)

	fmt.Println("Server spun up and listening on port: ", globals.BetServerPort)

	if err := s.Srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}

	fmt.Println("Bet server closed")
}

func (s *Server) recommendations(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	sessionToken := r.URL.Query().Get("token")

	fmt.Println(sessionToken)

	userID := s.arch.GetIdForToken(sessionToken)

	if userID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	links := s.arch.GetReccomendations(userID)

	linkBytes, err := json.Marshal(links)
	assist.Check(err)

	fmt.Println(links)

	w.WriteHeader(http.StatusOK)
	w.Write(linkBytes)
}

func (s *Server) processURL(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := assist.SafeBytes(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var payload Payload
	err = json.Unmarshal(reqBytes, &payload)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(payload)

	userID := s.arch.GetIdForToken(payload.Token)

	fmt.Println(userID)

	if userID == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		isNew := s.arch.URLIsNew(payload.URL)

		linkID := s.arch.AddURL(payload.URL)
		s.arch.AddVisit(userID, linkID)

		if isNew {
			s.newURLInput <- shared.NewUnscannedLink(linkID, payload.URL)
		}
	}

	w.Write(s.responseBytes("Success"))
}

func (s *Server) extractDetails(w http.ResponseWriter, r *http.Request) (shared.Details, error) {
	reqBytes, err := assist.SafeBytes(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return shared.Details{}, err
	}

	var reg shared.Details
	err = json.Unmarshal(reqBytes, &reg)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return shared.Details{}, err
	}

	return reg, nil
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	reg, err := s.extractDetails(w, r)

	if err != nil {
		lg.L.Debug(err.Error())
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

func New() *Server {
	newURLInput := make(chan shared.Link, 200)

	return &Server{
		*NewAuthServer(globals.BetServerPort),
		database.NewArchivist(),
		newURLInput,
		scanner.New(newURLInput),
	}
}
