package db

import "encoding/json"

type User struct {
	ID       string
	Username string
}

type Post struct {
	ID       string
	UserID   string
	PostBody json.RawMessage
}
