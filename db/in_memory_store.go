package db

import "sync"

type InMemoryCommerceStore struct {
	mu     sync.Mutex
	chains map[int]Chain
}

func NewInMemoryCommerceStore() *InMemoryCommerceStore {
	return &InMemoryCommerceStore{chains: map[int]Chain{}}
}

func (i *InMemoryCommerceStore) GetChains() []Chain {
	c := []Chain{}
	i.mu.Lock()
	defer i.mu.Unlock()
	for _, v := range i.chains {
		c = append(c, v)
	}
	return c
}
