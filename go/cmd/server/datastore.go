package main

import "github.com/peterahl/storytel/go/pkg/models"

type dataStore interface {
	GetMessages() (error, []models.Message)
	GetMessage(uint64) (error, models.Message)
	UpdateMessage(models.Message) error
	NewMessage(models.Message) error
	DeleteMessage(models.Message) error
}
