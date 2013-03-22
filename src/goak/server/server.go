package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Server struct {
	values  map[string]string
	peer    string
	running bool
}

func New() *Server {
	return &Server{
		values:  make(map[string]string),
		peer:    "",
	}
}

type GoakRequest struct {
	*http.Request
}

func (r *GoakRequest) fetchValue() string {
	body, _ := ioutil.ReadAll(r.Body)
	return string(body)
}

func (r *GoakRequest) fetchKey() string {
	uri := r.RequestURI
	uriParts := strings.Split(uri, "/")
	key := uriParts[len(uriParts)-1]
	return key
}

func (server *Server) AddPeer(peer string) {
	server.peer = peer
}

func (server *Server) hasPeer() bool {
	return server.peer != ""
}

func (server *Server) Listen(port string) {
	http.ListenAndServe(port, server.Handler())
}

func (server *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/heartbeat", func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "OK")
	})

	mux.HandleFunc("/data/", func(w http.ResponseWriter, request *http.Request) {
		r := GoakRequest{request}
		var value string

		switch r.Method {
		case "GET":
			key := r.fetchKey()
			value = server.values[key]

			if value == "" {
				if server.hasPeer() {
					response, _ := http.Get(server.peer + "/data/" + key)
					if response.StatusCode == 200 {
						body, _ := ioutil.ReadAll(response.Body)
						value = string(body)
					} else {
						w.WriteHeader(404)
					}
				} else {
					w.WriteHeader(404)
				}
			} else {
				w.WriteHeader(200)
			}
		case "PUT":
			_, keyExists := server.values[r.fetchKey()]
			value = r.fetchValue()
			server.values[r.fetchKey()] = value
			if keyExists {
				w.WriteHeader(200)
			} else {
				w.WriteHeader(201)
			}
		}

		fmt.Fprintf(w, "%s", value)
	})

	return mux
}
