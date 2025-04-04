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
	ticks     chan adapter.Tick
	orders    chan adapter.Order
	pool      *pgxpool.Pool
	errors    chan error
	recording bool
}

func New(dbUrl string, recording bool, ticks chan adapter.Tick, orders chan adapter.Order) (*dbwriter, error) {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %v", err)
	}
	return &dbwriter{ticks: ticks, orders: orders, pool: pool, recording: recording}, nil
}

func (d *dbwriter) Record(ticksTableName string, ordersTableName string) {
	for {
		select {
		case t := <-d.ticks:
			go func() {
				if !d.recording {
					log.Println("RECEIVED", t)
				} else {
					log.Println("WRITING", t)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					_, err := d.pool.Exec(ctx, fmt.Sprintf("INSERT INTO %s (symbol, price, timestamp) VALUES ($1, $2, $3)", ticksTableName), t.Symbol(), t.Price(), t.Timestamp().UnixNano())
					if err != nil {
						d.errors <- fmt.Errorf("failed to write to db: %v", err)
					}
				}

			}()
		case o := <-d.orders:
			go func() {
				if o.Size() > 0 {
					if !d.recording {
						log.Println("RECEIVED", o)
					} else {
						log.Println("WRITING", o)
						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()
						_, err := d.pool.Exec(ctx, fmt.Sprintf("INSERT INTO %s (symbol, side, price, size, timestamp) VALUES ($1, $2, $3, $4, $5)", ordersTableName), o.Symbol(), o.Side(), o.Price(), o.Size(), o.Timestamp().UnixNano())
						if err != nil {
							d.errors <- fmt.Errorf("failed to write to db: %v", err)
						}
					}
				}
			}()
		}
	}
}

func (d *dbwriter) Ticks() chan adapter.Tick {
	return d.ticks
}

func (d *dbwriter) Orders() chan adapter.Order {
	return d.orders
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
