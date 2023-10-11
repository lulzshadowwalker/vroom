package model

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Rooms     []Room `json:"rooms"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updaated_at"`
}
