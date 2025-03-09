package handlers

import (
	db "example.com/source/database"
)

var Db *db.PostgresConnection = db.NewPostgresConnection()
