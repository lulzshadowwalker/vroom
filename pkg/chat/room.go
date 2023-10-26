package chat

type Room struct {
	Id      string           `json:"id"`
	Name    string           `json:"name"`
	Members map[*Client]bool `json:"members"`
}

func (r *Room) addMember(c *Client) {
	if r.Members == nil {
		r.Members = make(map[*Client]bool)
	}

	r.Members[c] = true
}
