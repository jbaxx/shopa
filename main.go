package main

import (
	"gg.rocks/shopa/db"
	"gg.rocks/shopa/server"
)

func main() {
	port := ":5000"

	store := db.NewInMemoryCommerceStore()
	server := server.NewCommerceServer(store, port)

	server.Run(port)
}
