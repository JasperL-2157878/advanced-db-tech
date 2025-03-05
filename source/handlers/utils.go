package handlers

import (
	"encoding/json"

	"example.com/source/database"
)

var Db *database.DatabasePool = database.Connect()

func JSON(data any) []byte {
	payload, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return payload
}
