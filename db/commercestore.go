package db

import (
	"time"
)

type CommerceStore interface {
	GetChains() []Chain
}

type Chain struct {
	Id        *string    `json:"id"`
	Name      *string    `json:"name,omitempty"`
	Email     *string    `json:"email,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
