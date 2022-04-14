package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Dagetby/simplesurance/internal/counter"
	"github.com/Dagetby/simplesurance/internal/handler"
)

const (
	path = "dates.txt"
	port = ":8080"
)

func main() {
	ctx := context.Background()

	co := counter.MustCounter(ctx, path)
	dummy := handler.New(co)

	mux := http.DefaultServeMux
	mux.HandleFunc("/simplesurance", dummy.DummyHandler)

	log.Printf("Server staring at port: %s\n", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Server shutdownt")
}
