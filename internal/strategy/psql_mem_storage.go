package strategy

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

type PsqlMemStore struct {
	db  *pgx.Conn
	mem *MemStore
}

func NewPsqlMemStore(conn *pgx.Conn, mem *MemStore) (*PsqlMemStore, error) {
	query := `SELECT id, url FROM shortener`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	store := PsqlMemStore{conn, mem}
	var id, url string

	for rows.Next() {
		err = rows.Scan(&id, &url)
		if err != nil {
			return nil, err
		}
		if err = store.mem.Set(id, url); err != nil {
			return nil, err
		}
	}
	log.Println("record recovery completed")

	return &store, nil
}

func (p *PsqlMemStore) Set(key, val string) error {
	query := `INSERT INTO shortener (id, url) VALUES ($1, $2)`
	if _, err := p.db.Exec(context.Background(), query, key, val); err != nil {
		return err
	}

	if err := p.mem.Set(key, val); err != nil {
		return err
	}
	return nil
}

func (p *PsqlMemStore) Get(key string) (string, bool) {
	val, ok := p.mem.Get(key)
	if ok {
		return val, ok
	}

	query := `SELECT url FROM shortener WHERE id=$1`
	if err := p.db.QueryRow(context.Background(), query, key).Scan(&val); err != nil {
		return "", false
	}

	_ = p.mem.Set(key, val)
	return val, true
}
