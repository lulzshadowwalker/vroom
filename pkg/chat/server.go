package chat

import (
	"fmt"
	"net"
)

type Server struct {
	Address          *net.TCPAddr
	BroadcastChannel chan Message
	Clients          map[*Client]bool
	Rooms            map[string]*Room
}

func NewServer(address string) (*Server, error) {
	tcpAddress, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve tcp address %w", err)
	}

	return &Server{
		Address:          tcpAddress,
		BroadcastChannel: make(chan Message),
		Clients:          make(map[*Client]bool),
		Rooms:            make(map[string]*Room),
	}, nil
}

func (s *Server) Listen() (*net.TCPListener, error) {
	l, err := net.ListenTCP("tcp4", s.Address)
	if err != nil {
		return nil, fmt.Errorf("cannot listen on %s %w\n", s.Address.String(), err)
	}

	fmt.Printf("🥑 listening on %s\n", s.Address.String())

	return l, nil
}

func (s *Server) HandleBroadcast() {
	for {
		msg := <-s.BroadcastChannel

		if s.Rooms[msg.RoomId] == nil {
			continue
		}

		for m := range s.Rooms[msg.RoomId].Members {
			m.Send(msg)
		}
	}
}

func (s *Server) Run() {
	for {
		msg := <-s.BroadcastChannel
		if s.Rooms[msg.RoomId] == nil {
			continue
		}

		for m := range s.Rooms[msg.RoomId].Members {
			m.Send(msg)
		}
	}
}
