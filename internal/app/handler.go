package app

import (
	"context"
	"encoding/json"
	"github.com/valen0k/url-shortener-swoyo/internal/utils"
	"log"
	"net/http"
	"strings"
)

func (a *App) newHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.shortening)

	return mux
}

func (a *App) shortening(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.post(w, r)
	case http.MethodGet:
		a.get(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Method:", r.Method, "Status: StatusMethodNotAllowed")
		return
	}
}

func (a *App) post(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("URI:", r.RequestURI, "Status: StatusBadRequest")
		return
	}

	event := struct {
		Url string `json:"url"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	hash := utils.GenerateHash(event.Url)
	a.Set(hash, event.Url)

	if a.d {
		query := `INSERT INTO test (id, url) VALUES ($1, $2)`
		if _, err := a.db.Exec(context.Background(), query, hash, event.Url); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	event.Url = a.url + hash

	marshal, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (a *App) get(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.RequestURI, "/")
	if len(split) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("URI:", r.RequestURI, "Status: StatusBadRequest")
		return
	}
	hash := split[1]
	get, ok := a.Get(hash)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		log.Println("Not found hash:", hash, "Status: StatusForbidden")
		return
	}

	event := struct {
		Url string `json:"url"`
	}{Url: get}

	marshal, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
