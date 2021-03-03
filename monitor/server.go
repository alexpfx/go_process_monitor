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

)

type Msg struct {
	Command Command
	Process string
	Term    string
}

type Server struct {
	Port int
	User string
	Pass string
}

func (s Server) Start() (ch chan Msg) {
	ch = make(chan Msg)
	route := s.createRoute(ch)
	go s.serve(s.Port, route)
	fmt.Printf("iniciando serviço na porta : %d\n", s.Port)

	return
}

func (s Server) serve(port int, route *mux.Router) {
	err := http.ListenAndServe(fmt.Sprint(":", port), route)
	util.Check(err)
}

func (s Server) checkAuth(r *http.Request) bool {
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

func (s Server) createRoute(ch chan Msg) *mux.Router {

	route := mux.NewRouter()

	route.HandleFunc("/register/{process}/{term}", func(w http.ResponseWriter, r *http.Request) {
		if !s.checkAuth(r) {
			authError(w)
			return
		}

		params := mux.Vars(r)
		process := params["process"]
		term := params["term"]

		ch <- Msg{
			Command: Register,
			Process: process,
			Term:    term,
		}

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
