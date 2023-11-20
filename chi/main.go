package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	srv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}
	srv.ListenAndServe()
	// http.ListenAndServe(":3000", r)
}
