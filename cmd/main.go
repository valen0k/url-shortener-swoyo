package main

import (
	"github.com/valen0k/url-shortener-swoyo/internal/app"
	"log"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatalln(err)
	}
	if err = app.Run(); err != nil {
		log.Fatalln(err)
	}
}
