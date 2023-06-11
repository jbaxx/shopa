package main

import (
	"context"
	"log"
	"os"

	"gg.rocks/shopa/db"
	"gg.rocks/shopa/server"
)

func main() {
	port := ":5000"
	l := log.New(os.Stdout, "commerce: ", log.LstdFlags)

	// store := db.NewInMemoryCommerceStore()
	ctx := context.Background()
	store, err := db.NewRoachCommerceStore(ctx)
	if err != nil {
		l.Fatal(err)
	}
	server := server.NewCommerceServer(store, port, l)

	server.Run(port)
}
