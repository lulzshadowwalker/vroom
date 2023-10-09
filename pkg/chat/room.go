package chat

type Room struct {
	Id      string           `json:"id"`
	Name    string           `json:"name"`
	Members map[*Client]bool `json:"members"`
}

func (r *Room) addMember(c *Client) {
	r.Members[c] = true
}
