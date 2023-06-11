package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gg.rocks/shopa/db"
)

const ContentApplicationJSON = "application/json"

type CommerceServer struct {
	store db.CommerceStore
	// http.Handler
	server *http.Server
}

func NewCommerceServer(store db.CommerceStore, port string) *CommerceServer {
	cs := new(CommerceServer)

	cs.store = store

	router := http.NewServeMux()
	router.Handle("/chains", http.HandlerFunc(cs.chainsHandler))

	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      router,
	}

	cs.server = srv

	return cs
}

func (c *CommerceServer) chainsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", ContentApplicationJSON)
	json.NewEncoder(w).Encode(c.store.GetChains())
}

func (c *CommerceServer) Run(port string) {
	go func() {
		log.Printf("Starting server at port: %s", port)
		if err := c.server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	is := make(chan os.Signal, 1)
	signal.Notify(is, os.Interrupt)
	signal.Notify(is, os.Kill)

	sig := <-is
	log.Printf("got signal %v, shuting down gracefully", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c.server.Shutdown(ctx)
	log.Println("Shut down")
}
