package main

import (
	"encoding/json"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	lift "github.com/liftbridge-io/go-liftbridge"
	"github.com/peterahl/coding-assignment/go/pkg/models"
)

func newRouter(ds dataStore, client lift.Client) *mux.Router {

	r := mux.NewRouter()

	r.Handle("/messages", getMessages(ds)).Methods("GET")
	r.Handle("/list", getCmds(ds)).Methods("GET")
	r.Handle("/messages", newMessage(client)).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}", updateMessage(client)).Methods("PUT")
	r.Handle("/messages/{id:[0-9]+}", getMessage(ds)).Methods("GET")
	r.Handle("/messages/{id:[0-9]+}", deletMessage(client)).Methods("DELETE")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	return r

}

func getCmds(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err, msgs := db.GetCmds()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(msgs); err != nil {
				panic(err)
			}
		}
	})

}

func getMessages(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err, msgs := db.GetMessages()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(msgs); err != nil {
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

func newMessage(client lift.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var msg models.Message
		if err := json.Unmarshal([]byte(body), &msg); err != nil {
			log.Println(body, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		msg.Cmd = "create"
		data, err := proto.Marshal(&msg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := client.Publish(context.Background(), "foo", data); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func updateMessage(client lift.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var msg models.Message
		if err := json.Unmarshal([]byte(body), &msg); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		msg.Id = id
		msg.Cmd = "update"
		data, err := proto.Marshal(&msg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := client.Publish(context.Background(), "foo", data); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func deletMessage(client lift.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var msg models.Message
		if err := json.Unmarshal([]byte(body), &msg); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		msg.Id = id
		msg.Cmd = "delete"
		data, err := proto.Marshal(&msg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := client.Publish(context.Background(), "foo", data); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}
