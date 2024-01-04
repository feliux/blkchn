package utils

import (
	"encoding/json"
	"log"
)

func JsonStatus(message string) []byte {
	m, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	if err != nil {
		log.Printf("ERROR marshaling data: %s" + err.Error())
	}
	return m
}
