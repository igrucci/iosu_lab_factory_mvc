package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"os"
)

func ConnectDB(cfg *Config) (*pgxpool.Pool, error) {

	ConnectStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	dbpool, err := pgxpool.New(context.Background(), ConnectStr)
	if err != nil {
		fmt.Errorf("Unable to create connection pool: %v", err)
		os.Exit(1)
	}

	return dbpool, nil
}
