package db

import "encoding/json"

/*
User represents a user in the system.
Contains basic user information for authentication and identification.
*/
type User struct {
	ID        string `json:"id"`
	Regno     string `json:"regno"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
}

/*
Post represents a placement log post in the system.
Contains post content, ownership information, and review status.
*/
type Post struct {
	ID       string          `json:"id"`
	UserID   string          `json:"user_id"`
	PostBody json.RawMessage `json:"post_body"`
	Reviewed bool            `json:"reviewed"`
}

/*
Round represents a single round in a placement process.
Contains the content/description of the round.
*/
type Round struct {
	Content string `json:"content"`
}

/*
PostBody represents the structured content of a placement log post.
Contains company information, role details, and interview rounds.
*/
type PostBody struct {
	Company string  `json:"company"`
	Role    string  `json:"role"`
	Rounds  []Round `json:"rounds"`
}

/*
Admin represents an admin user in the system.
Contains basic admin information for authentication and identification.
*/
type Admin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

/*
PlacementCompany represents a placement event (company, ctc, date, etc).
*/
type PlacementCompany struct {
	ID            int           `json:"id"`
	Company       string        `json:"company"`
	CTC           float64       `json:"ctc"`
	PlacementDate string        `json:"placement_date"`
	CreatedAt     string        `json:"created_at"`
	BranchCounts  []BranchCount `json:"branch_counts,omitempty"`
}

/*
BranchCount represents the number of students placed from a branch in a placement event.
*/
type BranchCount struct {
	Branch string `json:"branch"`
	Count  int    `json:"count"`
}
