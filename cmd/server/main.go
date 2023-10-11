package main

import (
	"fmt"
	"os"

	"github.com/lulzshadowwalker/vroom/pkg/chat"
)

func main() {
	s, err := chat.NewServer(":3124")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot instantiate a server %q", err)
		os.Exit(1)
	}

	s.Rooms["hello"] = &chat.Room{Id: "hello", Members: make(map[*chat.Client]bool)}

	l, err := s.Listen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer l.Close()

	go s.HandleBroadcast()

	for {
		con, err := l.AcceptTCP()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot accept incoming handshake %q\n", err)
			continue
		}

		client := chat.NewClient(s, con)
		client.JoinRoom("general")
		s.Rooms["hello"].Members[client] = true
		go client.Handle()
		s.Clients[client] = true
	}
}
