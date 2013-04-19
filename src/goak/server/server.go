package server

import (
	"github.com/bmizerany/pat"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	values  map[string]string
	peer    string
	running bool
}

func New() *Server {
	return &Server{
		values: make(map[string]string),
		peer:   "",
	}
}

func (server *Server) AddPeer(peer string) {
	server.peer = peer
}

func (server *Server) hasPeer() bool {
	return server.peer != ""
}

func (server *Server) Handler() http.Handler {
	m := pat.New()

	m.Put("/data/:key", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get(":key")
		body, _ := ioutil.ReadAll(request.Body)
		value := string(body)

		server.values[key] = value
		w.WriteHeader(201)
		io.WriteString(w, value)
	}))

	m.Get("/data/:key", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get(":key")
		value, ok := server.values[key]

		if !ok {
			if server.hasPeer() {
				response, _ := http.Get(server.peer + "/data/" + key)
				if response.StatusCode == 200 {
					body, _ := ioutil.ReadAll(response.Body)
					value = string(body)
					w.WriteHeader(200)
					io.WriteString(w, value)
				} else {
					w.WriteHeader(404)
				}
			} else {
				w.WriteHeader(404)
			}
		} else {
			w.WriteHeader(200)
			io.WriteString(w, value)
		}
	}))

	return m
}
