package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/valen0k/url-shortener-swoyo/internal/config"
	"log"
	"net/http"
	"sync"
	"time"
)

type App struct {
	d          bool
	configFile string
	config     *config.Config
	db         *pgx.Conn
	buf        struct {
		sync.RWMutex
		memory map[string]string
	}
	url string
}

func NewApp() (*App, error) {
	app := App{}
	flag.BoolVar(&app.d, "d", false, "хранить информацию в postgres")
	flag.StringVar(&app.configFile, "c", "config.json", "конфигурационный файл")
	flag.Parse()

	err := app.uploadConfig()
	if err != nil {
		return nil, err
	}

	if app.d {
		app.db, err = app.newDBConnection()
		if err != nil {
			return nil, err
		}
		if err = app.memoryRecovery(); err != nil {
			return nil, err
		}
	} else {
		app.buf.memory = make(map[string]string)
	}

	app.url = fmt.Sprintf("http://%s:%s/", app.config.Server.Host, app.config.Server.Port)

	return &app, nil
}

func (a *App) Run() error {
	server := http.Server{
		Addr:         a.config.Server.Host + ":" + a.config.Server.Port,
		Handler:      a.newHandler(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("app started")
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) uploadConfig() error {
	config, err := config.NewConfig(a.configFile, a.d)
	if err != nil {
		return err
	}
	a.config = config
	return nil
}

func (a *App) memoryRecovery() error {
	query := `SELECT id, url FROM test`
	rows, err := a.db.Query(context.Background(), query)
	if err != nil {
		return err
	}

	var id, url string
	a.buf.memory = make(map[string]string)

	for rows.Next() {
		err = rows.Scan(&id, &url)
		if err != nil {
			return err
		}
		a.Set(id, url)
	}

	log.Println("record recovery completed")
	return nil
}

func (a *App) Set(key, val string) {
	a.buf.Lock()
	defer a.buf.Unlock()

	log.Println("recorded in memory")
	a.buf.memory[key] = val
}

func (a *App) Get(key string) (string, bool) {
	a.buf.RLock()
	defer a.buf.RUnlock()

	value, ok := a.buf.memory[key]
	if ok {
		return value, ok
	}
	return "", ok
}

func (a *App) newDBConnection() (*pgx.Conn, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			a.config.Storage.User, a.config.Storage.Password,
			a.config.Storage.Host, a.config.Storage.Port,
			a.config.Storage.Database))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
