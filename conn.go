package database

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	pgxgeom "github.com/twpayne/pgx-geom"
)

// Conn represents a connection to the database.
type Conn struct {
	pool *pgxpool.Pool
}

// ConnParams holds the parameters for connecting to the database.
type ConnParams struct {
	Addr       string
	HasPostgis bool
}

// NewConn creates a new database connection using the provided parameters.
func NewConn(ctx context.Context, p *ConnParams) (*Conn, error) {
	config, err := pgxpool.ParseConfig(p.Addr)
	if err != nil {
		return nil, err
	}

	if p.HasPostgis {
		config.AfterConnect = pgxgeom.Register
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Conn{pool: pool}, nil
}

// Close the database connection.
func (c *Conn) Close() {
	c.pool.Close()
}

// Pool returns the underlying pgxpool.Pool.
func (c *Conn) Pool() *pgxpool.Pool {
	return c.pool
}

// Ping checks the database connection.
func (c *Conn) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

// Handle returns a *sql.DB instance for compatibility with database/sql package, acting as a connection handle.
func (c *Conn) Handle() *sql.DB {
	return stdlib.OpenDBFromPool(c.pool)
}
