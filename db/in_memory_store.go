package db

import (
	"sync"
	"time"
)

type InMemoryCommerceStore struct {
	mu     sync.Mutex
	chains map[string]Chain
}

func String(a string) *string {
	return &a
}

func NewInMemoryCommerceStore() *InMemoryCommerceStore {
	created := time.Date(2022, time.August, 29, 19, 0, 0, 0, time.UTC)
	updated := time.Date(2022, time.December, 29, 19, 0, 0, 0, time.UTC)
	m := map[string]Chain{
		"5bb2b713-e8e3-4b65-8fd1-27975a07093b": {
			Id:        String("5bb2b713-e8e3-4b65-8fd1-27975a07093b"),
			Name:      String("Gigazoom"),
			Email:     String("sales@gigazoom"),
			CreatedAt: &created,
			UpdatedAt: &updated,
		},
		"c70b5c82-b286-4665-bad9-8743df3e35ff": {
			Id:        String("c70b5c82-b286-4665-bad9-8743df3e35ff"),
			Name:      String("Thoughtmix"),
			Email:     String("sales@thoughtmix"),
			CreatedAt: &created,
			UpdatedAt: &updated,
		},
		"876b577c-fb70-49c2-b2fb-a315754ea0e0": {
			Id:        String("876b577c-fb70-49c2-b2fb-a315754ea0e0"),
			Name:      String("Divavu"),
			Email:     String("sales@divavu"),
			CreatedAt: &created,
			UpdatedAt: &updated,
		},
		"6b7650d3-f339-40c8-b6b1-93c90818ba54": {
			Id:        String("6b7650d3-f339-40c8-b6b1-93c90818ba54"),
			Name:      String("Skyba"),
			Email:     String("sales@skyba"),
			CreatedAt: &created,
			UpdatedAt: &updated,
		},
		"024b4eba-5cf4-4d49-823b-ffde6d710ee1": {
			Id:        String("024b4eba-5cf4-4d49-823b-ffde6d710ee1"),
			Name:      String("Meevee"),
			Email:     String("sales@meevee"),
			CreatedAt: &created,
			UpdatedAt: &updated,
		},
	}

	return &InMemoryCommerceStore{chains: m}
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
