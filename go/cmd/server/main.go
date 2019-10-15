package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/peterahl/storytel/go/pkg/memstore"
)

type dataStore interface {
	GetMessages() (error, []memstore.Message)
	GetMessage(uint64) (error, memstore.Message)
	UpdateMessage(uint64, string) error
	NewMessage(string) error
	DeleteMessage(uint64) error
}

type service struct {
	db dataStore
}

func main() {

	db := &memstore.Store{
		Messages: make(map[uint64]memstore.Message),
	}

	r := newRouter(db)

	fmt.Println("Starting server")

	log.Fatal(http.ListenAndServe(":3000", r))
}
