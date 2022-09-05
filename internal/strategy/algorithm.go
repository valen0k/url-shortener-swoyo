package strategy

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"sync"
)

type EvictionAlgo interface {
	Set(key, val string) error
	Get(key string) (string, bool)
}

type Memory struct {
	sync.RWMutex
	memory map[string]string
}

type PsqlMemStore struct {
	db  *pgx.Conn
	mem *Memory
}

func NewPsqlMemStore(conn *pgx.Conn) (*PsqlMemStore, error) {
	mem := &Memory{}

	query := `SELECT id, url FROM test`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var id, url string
	mem.memory = make(map[string]string)

	for rows.Next() {
		err = rows.Scan(&id, &url)
		if err != nil {
			return nil, err
		}
		mem.memory[id] = url
	}
	log.Println("record recovery completed")

	return &PsqlMemStore{conn, mem}, nil
}

func (p *PsqlMemStore) Set(key, val string) error {
	query := `INSERT INTO test (id, url) VALUES ($1, $2)`
	if _, err := p.db.Exec(context.Background(), query, key, val); err != nil {
		return err
	}

	p.mem.Lock()
	defer p.mem.Unlock()
	p.mem.memory[key] = val
	log.Println("recorded in memory")
	return nil
}

func (p *PsqlMemStore) Get(key string) (string, bool) {
	//для примера, что можно подтягивать напрямую из бд,
	//но такой вариант не эффективен, поскольку есть memory
	var val string
	query := `SELECT url FROM test WHERE id=$1`
	if err := p.db.QueryRow(context.Background(), query, key).Scan(&val); err != nil {
		return "", false
	}
	return val, true
}

type MemStore struct {
	mem *Memory
}

func NewMemStore() *MemStore {
	mem := &Memory{}
	mem.memory = make(map[string]string)

	return &MemStore{mem}
}

func (m *MemStore) Set(key, val string) error {
	m.mem.Lock()
	defer m.mem.Unlock()
	m.mem.memory[key] = val
	log.Println("recorded in memory")
	return nil
}

func (m *MemStore) Get(key string) (string, bool) {
	m.mem.RLock()
	defer m.mem.RUnlock()

	value, ok := m.mem.memory[key]
	if ok {
		return value, ok
	}
	return "", ok
}
