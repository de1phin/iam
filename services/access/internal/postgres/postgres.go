package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Options struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Storage struct {
	conn *sql.DB
}

func New(ctx context.Context, options Options) (*Storage, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		options.Host,
		options.Port,
		options.User,
		options.Password,
		options.DBName,
	)
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	err = conn.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	s := &Storage{conn: conn}
	err = s.initSchema(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) initSchema(ctx context.Context) error {
	initSchemaSQL := `
	CREATE TABLE IF NOT EXISTS roles (
		name TEXT PRIMARY KEY,
		permissions TEXT[]
	);
	
	CREATE TABLE IF NOT EXISTS access_bindings (
		user_id TEXT,
		role_name TEXT,
		resource TEXT
	);`
	_, err := s.conn.ExecContext(ctx, initSchemaSQL)
	return err
}
