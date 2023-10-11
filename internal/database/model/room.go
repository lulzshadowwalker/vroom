package model

type Room struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Members   []User `json:"members"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updaated_at"`
}
