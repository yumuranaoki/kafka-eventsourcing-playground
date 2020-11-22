package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	port                = ":8080"
	brokers             = []string{"localhost:9092"}
	group   goka.Group  = "collector"
	table   goka.Table  = goka.GroupTable(group)
	stream  goka.Stream = "message"
)

func main() {
	view, err := goka.NewView(brokers, table, new(codec.String))
	if err != nil {
		panic(err)
	}
	go view.Run(context.Background())

	emitter, err := goka.NewEmitter(brokers, stream, new(codec.String))
	if err != nil {
		panic(err)
	}
	defer emitter.Finish()

	router := mux.NewRouter()
	router.HandleFunc("/read", read(view)).Methods("GET")
	router.HandleFunc("/register", register(emitter)).Methods("POST")

	log.Printf("app is listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
