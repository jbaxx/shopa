package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type RoachCommerceStore struct {
	mu sync.Mutex

	dbConn *pgx.Conn
}

func checkVariable(s string) error {
	if _, ok := os.LookupEnv(s); !ok {
		return fmt.Errorf("%s environment variable not set", s)
	}
	return nil
}

const PGUser = "PGUSER"
const PGPassword = "PGPASSWORD"
const PGHost = "PGHOST"
const PGPort = "PGPORT"
const PGDatabase = "PGDATABASE"

func NewRoachCommerceStore(ctx context.Context) (*RoachCommerceStore, error) {

	if err := checkVariable(PGUser); err != nil {
		return nil, err
	}
	if err := checkVariable(PGPassword); err != nil {
		return nil, err
	}
	if err := checkVariable(PGHost); err != nil {
		return nil, err
	}
	if err := checkVariable(PGPort); err != nil {
		return nil, err
	}
	if err := checkVariable(PGDatabase); err != nil {
		return nil, err
	}

	pgConnString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv(PGUser),
		os.Getenv(PGPassword),
		os.Getenv(PGHost),
		os.Getenv(PGPort),
		os.Getenv(PGDatabase),
	)

	dbConn, err := pgx.Connect(
		ctx,
		pgConnString,
	)
	if err != nil {
		return nil, err
	}

	return &RoachCommerceStore{
		dbConn: dbConn,
	}, nil
}

func (c *RoachCommerceStore) GetChains() []Chain {
	c.mu.Lock()
	defer c.mu.Unlock()

	var ch []Chain

	err := pgxscan.Select(context.Background(), c.dbConn, &ch, "select * from chains")
	if err != nil {
		fmt.Printf("error getting chains: %v\n", err)
	}

	return ch
}
