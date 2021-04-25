package monitor

import (
	"encoding/base64"
	"fmt"
	"github.com/alexpfx/go_listener/util"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Command int

const (
	Register Command = iota
	Unregister
	Quit
)

type Observer struct {
	Command Command `json:"command"`
	Process string  `json:"process"`
	Term    string  `json:"term"`
	Pattern string  `json:"pattern"`
}

type OldServer struct {
	Port int
	User string
	Pass string
}

func (s OldServer) Start(ch chan Observer) {
	route := s.createRoute(ch)
	s.serve(s.Port, route)
	fmt.Printf("iniciando serviço na porta : %d\n", s.Port)
}

func (s OldServer) serve(port int, route *mux.Router) {
	err := http.ListenAndServe(fmt.Sprint(":", port), route)
	util.Check(err)
}

func (s OldServer) checkAuth(r *http.Request) bool {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 {
		return false

	}

	ds, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(ds), ":", 2)

	return pair[0] == s.User && pair[1] == s.Pass
}

func (s OldServer) createRoute(ch chan Observer) *mux.Router {

	route := mux.NewRouter()

	route.Path("/register").
		Queries("ps", "{ps}").
		Queries("term", "{term}").
		Queries("p", "{pattern}").

		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !s.checkAuth(r) {
				authError(w)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			vs := r.URL.Query()

			process := vs.Get("ps")
			term := vs.Get("term")
			pattern := vs.Get("p")

			ch <- Observer{
				Command: Register,
				Process: process,
				Term:    term,
				Pattern: pattern,
			}

			msg := "registrando: \nps %s\ntermo: %s\npattern: %s\n"
			w.Write([]byte(fmt.Sprintf(msg, process, term, pattern)))
		})

	return route
}

func splitParam(callback string) []string {
	if strings.Contains(callback, " ") {
		return strings.Split(callback, " ")
	}
	return strings.Split(callback, "+")
}

func authError(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}