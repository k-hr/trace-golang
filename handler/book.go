package handler

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/environment"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/model"
)

func CreateBook(env environment.Environment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var newBook model.Book
		err := json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = env.DBProxy.InsertBook(r.Context(), newBook)
		if err != nil {
			log.Printf("Error while inserting book - %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newBook)
	}
}

func ReadBook(env environment.Environment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]

		if title == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		book, err := env.DBProxy.ReadBook(r.Context(), title)
		if err == pg.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			log.Printf("Error while reading book - %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(book)
	}
}
