package handlers

import (
	"encoding/json"

	db "example.com/source/database"
)

var Db *db.PostgresConnection = db.NewPostgresConnection()

func JSON(data any) []byte {
	payload, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return payload
}
