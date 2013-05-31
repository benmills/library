package server

import (
	"github.com/bmizerany/pat"
	"io"
	"io/ioutil"
	"net/http"

	"goak/peer"
)

type Server struct {
	*peer.Peer
	values map[string]string
}

func New(url string) *Server {
	return &Server{peer.New(url), make(map[string]string)}
}

func (server *Server) Handler() http.Handler {
	m := pat.New()

	server.Peer.Handler(m)

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
			if server.HasPeer() {
				response, _ := http.Get(server.Peers[0] + "/data/" + key)
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
