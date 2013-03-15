package server

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
)

type SemanticTuple struct {
	Key string
	Value string
}

type Server struct {
	values map[string]string
}

func New() *Server {
	return &Server{
		values: make(map[string]string),
	}
}

type GoakRequest struct {
	*http.Request
}

func(r *GoakRequest) fetchValue() string {
	body, _ := ioutil.ReadAll(r.Body)
	return string(body)
}

func(r *GoakRequest) fetchKey() string {
	uri := r.RequestURI
	uriParts := strings.Split(uri, "/")
	key := uriParts[len(uriParts) - 1]
	return key
}

func (server *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		r := GoakRequest{request}
		var value string

		switch r.Method {
		case "GET":
			value = server.values[r.fetchKey()]
			if value == "" {
				w.WriteHeader(404)
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
