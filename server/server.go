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
	store  db.CommerceStore
	server *http.Server
	l      *log.Logger
}

func NewCommerceServer(store db.CommerceStore, port string, l *log.Logger) *CommerceServer {
	cs := new(CommerceServer)

	cs.store = store
	cs.l = l

	router := http.NewServeMux()
	router.Handle("/chains", http.HandlerFunc(cs.chainsHandler))

	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
		ErrorLog:     l,
		Handler:      router,
	}

	cs.server = srv

	return cs
}

func (c *CommerceServer) chainsHandler(w http.ResponseWriter, r *http.Request) {
	c.l.Println("GET /chains")
	w.Header().Set("content-type", ContentApplicationJSON)
	json.NewEncoder(w).Encode(c.store.GetChains())
}

func (c *CommerceServer) Run(port string) {
	go func() {
		c.l.Printf("Starting server at port: %s", port)
		if err := c.server.ListenAndServe(); err != nil {
			c.l.Println(err)
		}
	}()

	is := make(chan os.Signal, 1)
	signal.Notify(is, os.Interrupt)
	signal.Notify(is, os.Kill)

	sig := <-is
	c.l.Printf("got signal %v, shuting down gracefully", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c.server.Shutdown(ctx)
	c.l.Println("Shut down")
}
