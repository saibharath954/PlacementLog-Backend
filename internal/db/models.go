package db

import "encoding/json"

type User struct {
	ID       string
	Username string
}

type Post struct {
	ID       string          `json:"id"`
	UserID   string          `json:"user_id"`
	PostBody json.RawMessage `json:"post_body"`
}

type Round struct {
	Content string `json:"content"`
}

type PostBody struct {
	Company string  `json:"company"`
	Role    string  `json:"role"`
	Rounds  []Round `json:"rounds"`
}
