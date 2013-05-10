package server

import (
	"github.com/bmizerany/pat"
	"io"
	"io/ioutil"
	"net/http"
)

type Server struct {
	values  map[string]string
	peers    []string
	running bool
}

func New() *Server {
	return &Server{
		values: make(map[string]string),
		peers:   []string{},
	}
}

func (server *Server) addPeer(peer string) {
	server.peers = append(server.peers, peer)
}

func (server *Server) hasPeer() bool {
	return len(server.peers) > 0
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
				response, _ := http.Get(server.peers[0] + "/data/" + key)
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

	m.Put("/peers", http.HandlerFunc(func (w http.ResponseWriter, request *http.Request) {
		body, _ := ioutil.ReadAll(request.Body)
		newPeerURL := string(body)
		server.addPeer(newPeerURL)

		w.WriteHeader(201)
	}))

	m.Get("/peers", http.HandlerFunc(func (w http.ResponseWriter, request *http.Request) {
		if server.hasPeer() {
			w.WriteHeader(200)
			io.WriteString(w, JsonData{"peers": server.peers}.Encode())
		} else {
			w.WriteHeader(404)
		}
	}))

	return m
}
