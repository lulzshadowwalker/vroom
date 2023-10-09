package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	Id         string
	Name       string
	Server     *Server
	Connection *net.TCPConn
	Rooms      map[string]*Room
	ActiveRoom *Room
}

func NewClient(server *Server, con *net.TCPConn) *Client {
	return &Client{
		Server:     server,
		Connection: con,
		Rooms:      make(map[string]*Room),
	}
}

func (c *Client) Handle() {
	const maxPacketSize = 512
	buf := make([]byte, maxPacketSize)
	for {
		n, err := c.Connection.Read(buf)
		if err != nil {
			fmt.Println("client has disconnected")
			return
		}

		var msg Message
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			fmt.Fprintf(c.Connection, "invalid message format %q\n", err)
			continue
		}

		const (
			prefixCreate     = "create"
			prefixJoin       = "join"
			prefixLeave      = "leave"
			prefixDisconnect = "dis"
		)

		if !strings.HasPrefix(msg.Content, "/") {
			c.Server.BroadcastChannel <- msg
			continue
		}

		cmd := strings.TrimPrefix(msg.Content, "/")
		switch {
		case strings.HasPrefix(cmd, prefixJoin):
			roomId := strings.TrimPrefix(cmd, prefixJoin)
			err = c.JoinRoom(roomId)
			if err != nil {
				if errors.Is(err, ErrRoomNotExist) {
					fmt.Fprintln(c.Connection, "room does not exist")
					continue
				}

				fmt.Fprintln(c.Connection, "cannot join room")
				fmt.Printf("cannot join room %q\n", err)
			}

		case strings.HasPrefix(cmd, prefixLeave):
			err = c.LeaveRoom()
			if err != nil {
				if errors.Is(err, ErrNoActiveRoom) {
					fmt.Fprintln(c.Connection, "you're not connected to any room")
					continue
				}

				fmt.Fprintln(c.Connection, "cannot leave room")
				fmt.Printf("cannot leave room %q\n", err)
			}
		case strings.HasPrefix(cmd, prefixDisconnect):
			err = c.Disconnect()
			if err != nil {
				if errors.Is(err, ErrNoActiveRoom) {
					fmt.Fprintln(c.Connection, "you're not connected to any room")
					continue
				}

				fmt.Fprintln(c.Connection, "cannot disconnect from room")
				fmt.Printf("cannot disconnect from room %q\n", err)
			}
		default:
			fmt.Fprintln(c.Connection, "unrecognized command")
		}
		fmt.Println(string(buf[:n]))
	}
}

func (c *Client) String() string {
	return fmt.Sprintf("%s (%s)", c.Name, c.Id)
}

func (c *Client) Send(msg Message) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("cannot encode message %v %w", msg, err)
	}

	_, err = c.Connection.Write(jsonMsg)
	if err != nil {
		return fmt.Errorf("cannot send message %v %w", msg, err)
	}

	return nil
}

func (c *Client) JoinRoom(roomId string) error {
	if c.Server.Rooms[roomId] == nil {
		return ErrRoomNotExist
	}

	c.ActiveRoom = c.Server.Rooms[roomId]
	c.Server.Rooms[roomId].addMember(c)

	return nil
}

func (c *Client) LeaveRoom() error {
	if c.ActiveRoom == nil {
		return ErrNoActiveRoom
	}

	delete(c.Rooms, c.ActiveRoom.Id)
	delete(c.Server.Rooms, c.ActiveRoom.Id)
	c.ActiveRoom = nil

	return nil
}

func (c *Client) Disconnect() error {
	if c.ActiveRoom == nil {
		return ErrNoActiveRoom
	}

	c.ActiveRoom = nil
	return nil
}
