package database

import (
	"context"

	"github.com/de1phin/iam/pkg/logger"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Database interface {
	GetSingle(ctx context.Context, pointerOnDst any, query string, args ...any) error
	GetSlice(ctx context.Context, pointerOnSliceDst any, query string, args ...any) error

	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	Close()
}

type database struct {
	conn *pgxpool.Pool
}

func NewDatabase(ctx context.Context, dsn string) (Database, error) {
	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		logger.Error("database connection: ", zap.Error(err))
		return nil, err
	}

	db := &database{conn: conn}

	return db, nil
}

func (d *database) GetSingle(ctx context.Context, pointerOnDst any, query string, args ...any) error {
	err := pgxscan.Get(
		ctx,
		d.conn,
		pointerOnDst,
		query,
		args...,
	)
	return err
}

func (d *database) GetSlice(ctx context.Context, pointerOnSliceDst any, query string, args ...any) error {
	err := pgxscan.Select(
		ctx,
		d.conn,
		pointerOnSliceDst,
		query,
		args...,
	)
	return err
}

func (d *database) Exec(ctx context.Context, query string, args ...any) (commandTag pgconn.CommandTag, err error) {
	return d.conn.Exec(ctx, query, args...)
}

func (d *database) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return d.conn.Query(ctx, query, args...)
}

func (d *database) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.conn.QueryRow(ctx, query, args...)
}

func (d *database) Close() {
	d.conn.Close()
}
