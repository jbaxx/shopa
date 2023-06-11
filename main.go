package main

import (
	"log"
	"os"

	"gg.rocks/shopa/db"
	"gg.rocks/shopa/server"
)

func main() {
	port := ":5000"
	l := log.New(os.Stdout, "commerce: ", log.LstdFlags)

	store := db.NewInMemoryCommerceStore()
	server := server.NewCommerceServer(store, port, l)

	server.Run(port)
}
