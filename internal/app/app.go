package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/valen0k/url-shortener-swoyo/internal/config"
	"github.com/valen0k/url-shortener-swoyo/internal/strategy"
	"log"
	"net/http"
	"time"
)

type Store interface {
	Set(key, val string) error
	Get(key string) (string, bool)
}

type App struct {
	d          bool
	configFile string
	config     *config.Config
	engine     Store
	url        string
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
		db, err := app.newDBConnection()
		if err != nil {
			return nil, err
		}
		app.engine, err = strategy.NewPsqlMemStore(db)
		if err != nil {
			return nil, err
		}
	} else {
		app.engine = strategy.NewMemStore()
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

func (a *App) set(key, val string) error {
	return a.engine.Set(key, val)
}

func (a *App) getValue(key string) (string, bool) {
	return a.engine.Get(key)
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
