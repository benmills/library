package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"library/server"
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


	logger := log.New(os.Stdin, "["+host+":"+port+"] ", log.LstdFlags)
	if !verbose {
		logger = log.New(new(NullWriter), "", 0)
	}

	server := server.New("http://"+host+":"+port, logger)

	logger.Printf("Starting")

	http.Handle("/", server.Handler())
	err := http.ListenAndServe(host+":"+port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
