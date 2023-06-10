package main

import (
	"encoding/json"
	"net/http"
	"time"
)

const ContentApplicationJSON = "application/json"

type CommerceStore interface {
	GetChains() []Chain
}

type CommerceServer struct {
	store CommerceStore
	http.Handler
}

type Chain struct {
	Id        *int       `json:"id"`
	Name      *string    `json:"name,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func NewCommerceServer(store CommerceStore) *CommerceServer {
	c := new(CommerceServer)
	c.store = store

	router := http.NewServeMux()
	router.Handle("/chains", http.HandlerFunc(c.chainsHandler))

	c.Handler = router

	return c
}

func (c *CommerceServer) chainsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", ContentApplicationJSON)
	json.NewEncoder(w).Encode(c.store.GetChains())
}
