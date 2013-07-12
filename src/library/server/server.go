package server

import (
	"github.com/bmizerany/pat"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"library/httpclient"
	"library/peer"
)

type Server struct {
	*peer.Peer
	values map[string]string
	logger *log.Logger
}

func New(url string, logger *log.Logger) *Server {
	values := make(map[string]string)
	peer := peer.New(url, values, logger)
	return &Server{peer, values, logger}
}

func (server *Server) handoffKey(address string, key string, value string) {
	server.logger.Printf("Passing off '%s'->'%s' to %s", key, value, address)

	statusCode, _ := httpclient.Put(address+"/set/"+key, value)
	if statusCode == 0 {
		server.NotifyDown(address)
	}
}

func (server *Server) Handler() http.Handler {
	m := pat.New()

	server.Peer.Handler(m)

	m.Put("/set/:key", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get(":key")
		body, _ := ioutil.ReadAll(request.Body)
		value := string(body)

		server.logger.Printf("Setting '%s'->'%s'", key, value)
		server.values[key] = value
		w.WriteHeader(201)
		io.WriteString(w, value)

	}))

	m.Put("/data/:key", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get(":key")
		body, _ := ioutil.ReadAll(request.Body)
		value := string(body)

		for _, address := range(server.PreferenceListForKey(key)) {
			if server.URL() == address {
				server.logger.Printf("Storing '%s'->'%s'", key, value)
				server.values[key] = value
			} else {
				server.handOffKey(address, key, value)
			}
		}

		w.WriteHeader(201)
		io.WriteString(w, value)

	}))

	m.Get("/data/:key", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get(":key")


		for _, address := range(server.PreferenceListForKey(key)) {
			if server.URL() == address {
				value, ok := server.values[key]
				if ok {
					server.logger.Printf("Get key '%s' found value '%s'", key, value)
					w.WriteHeader(200)
					io.WriteString(w, value)
				} else {
					server.logger.Printf("Key '%s' not found", key)
					w.WriteHeader(404)
				}

				return;
			}
		}

		destinationAddress := server.PeerAddressForKey(key)

		server.logger.Printf("Passing off get of key '%s' to %s", key, destinationAddress)
		statusCode, response := httpclient.Get(destinationAddress+"/data/"+key, "")
		w.WriteHeader(statusCode)
		io.WriteString(w, response)
	}))

	m.Get("/stats/keys", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		stats := httpclient.JsonData{
			"count": len(server.values),
			"data": server.values,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		io.WriteString(w, stats.Encode())
	}))

	return m
}
