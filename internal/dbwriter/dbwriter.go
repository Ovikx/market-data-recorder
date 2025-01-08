package dbwriter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ovikx/market-data-recorder/internal/adapter"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbwriter struct {
	ticks  chan adapter.Tick
	pool   *pgxpool.Pool
	errors chan error
}

func New(dbUrl string) (*dbwriter, error) {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %v", err)
	}
	return &dbwriter{ticks: make(chan adapter.Tick), pool: pool}, nil
}

func (d *dbwriter) Record(tableName string) {
	for {
		select {
		case t := <-d.ticks:
			go func() {
				log.Println("WRITING", t)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := d.pool.Exec(ctx, fmt.Sprintf("INSERT INTO %s (symbol, price, timestamp) VALUES ($1, $2, $3)", tableName), t.Symbol(), t.Price(), t.Timestamp().UnixNano())
				if err != nil {
					d.errors <- fmt.Errorf("failed to write to db: %v", err)
				}
			}()
		}
	}
}

func (d *dbwriter) Ticks() chan adapter.Tick {
	return d.ticks
}

func (d *dbwriter) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *dbwriter) Errors() chan error {
	return d.errors
}

func (d *dbwriter) Close() {
	d.pool.Close()
}
