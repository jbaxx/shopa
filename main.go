package main

import (
	"context"
	"log"
	"os"

	"github.com/jbaxx/shopa/db"
	"github.com/jbaxx/shopa/server"
)

func main() {
	l := log.New(os.Stdout, "commerce: ", log.LstdFlags)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	port = ":" + port

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	if env != "dev" && env != "prod" {
		l.Fatalf("ENV must be set to either dev or prod, set to: %s", env)
	}
	l.Printf("Running server in a %q environment", env)

	var store db.CommerceStore
	var err error

	if env == "dev" {
		store = db.NewInMemoryCommerceStore()
	} else {
		ctx := context.Background()
		store, err = db.NewRoachCommerceStore(ctx)
		if err != nil {
			l.Fatal(err)
		}
	}

	server := server.NewCommerceServer(store, port, l)

	server.Run(port)
}
