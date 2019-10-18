package main

import "github.com/peterahl/coding-assignment/go/pkg/models"

type dataStore interface {
	GetMessages() (error, []models.Message)
	GetCmds() (error, []models.Message)
	AddCommand(models.Message) error
	GetMessage(uint64) (error, models.Message)
	UpdateMessage(models.Message) error
	NewMessage(models.Message) error
	DeleteMessage(models.Message) error
}
