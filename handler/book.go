package handler

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/environment"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/model"
	"strconv"
	"time"
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
		queryParams := r.URL.Query()
		apiDelay, _ := strconv.ParseBool(queryParams.Get("api-delay"))
		dbDelay, _ := strconv.ParseBool(queryParams.Get("db-delay"))
		runtimeErr, _ := strconv.ParseBool(queryParams.Get("error"))

		if title == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if apiDelay {
			time.Sleep(5 * time.Second)
		}

		book, err := env.DBProxy.ReadBook(r.Context(), title, dbDelay, runtimeErr)
		if err == pg.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if apiDelay {
			time.Sleep(5 * time.Second)
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
