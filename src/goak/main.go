package main

import (
	"flag"
	"log"
	"net/http"

	"goak/server"
)

type NullWriter int
func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	verbose := false
	flag.BoolVar(&verbose, "v", false, "verbose logging")

	port := "3333"
	flag.StringVar(&port, "port", "3333", "port")

	host := "localhost"
	flag.StringVar(&host, "host", "localhost", "host")

	flag.Parse()

	if !verbose {
		log.SetOutput(new(NullWriter))
	}

	server := server.New("http://"+host+":"+port)

	log.Printf("Starting on http://%s:%s", host, port)

	http.Handle("/", server.Handler())
	err := http.ListenAndServe(host+":"+port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
