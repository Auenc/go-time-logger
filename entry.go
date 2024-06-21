package main

import (
	"errors"
	"time"
)

type Entry struct {
	Id    uint      `json:"id"`
	Name  string    `json:"name"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func validateEntry(toValidate Entry) error {
	if toValidate.Name == "" {
		return errors.New("name must not be empty")
	}
	if toValidate.Start.After(toValidate.End) {
		return errors.New("start must be before end")
	}
	return nil
}
