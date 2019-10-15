package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func newRouter(ds dataStore) *mux.Router {

	r := mux.NewRouter()

	r.Handle("/messages", getMessages(ds)).Methods("GET")
	r.Handle("/messages", newMessage(ds)).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}", updateMessage(ds)).Methods("PUT")
	r.Handle("/messages/{id:[0-9]+}", getMessage(ds)).Methods("GET")
	r.Handle("/messages/{id:[0-9]+}", deletMessage(ds)).Methods("DELETE")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	return r

}

func getMessages(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err, msg := db.GetMessages()
		log.Println(msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(msg); err != nil {
				panic(err)
			}
		}
	})

}

func getMessage(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			_, msg := db.GetMessage(id)
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(msg); err != nil {
				log.Println("falid to parse json", err)
			}
		}
	})
}

func newMessage(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		bodyString := string(body)
		log.Printf("new message: %s", bodyString)
		if err := db.NewMessage(bodyString); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
}

func updateMessage(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		body, _ := ioutil.ReadAll(r.Body)
		bodyString := string(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			db.UpdateMessage(id, bodyString)
			w.WriteHeader(http.StatusOK)
		}
	})
}

func deletMessage(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			db.DeleteMessage(id)
			w.WriteHeader(http.StatusOK)
		}
	})
}
