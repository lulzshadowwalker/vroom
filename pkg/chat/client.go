// TODO, init server with rooms from db 
// TODO, register room operation with db
package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/go-playground/validator"
	"github.com/lulzshadowwalker/vroom/internal/database"
	"github.com/lulzshadowwalker/vroom/pkg/repo"
	"github.com/lulzshadowwalker/vroom/pkg/utils"
)

func init() {
	validate = validator.New()
}

var validate *validator.Validate

type Client struct {
	Id         int
	Name       string
	Server     *Server
	Connection *net.TCPConn
	ActiveRoom *Room
}

func NewClient(server *Server, con *net.TCPConn) *Client {
	return &Client{
		Server:     server,
		Connection: con,
	}
}

func (c *Client) Handle() {
	const maxPacketSize = 512
	buf := make([]byte, maxPacketSize)
	for {
		n, err := c.Connection.Read(buf[0:])
		if err != nil {
			fmt.Println("client has disconnected")
			return
		}

		var msg Message
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			c.Send(*c.Server.message(fmt.Sprintf("invalid message format %q\n", buf[:n])))
			continue
		}

		c.Id = msg.SenderId
		c.Name = msg.SenderName

		const (
			prefixCreate     = "create"
			prefixJoin       = "join"
			prefixLeave      = "leave"
			prefixDisconnect = "dis"
			prefixHelp       = "help"
		)

		if !strings.HasPrefix(msg.Content, "~") {
			if c.ActiveRoom == nil {
				c.Send(*c.Server.message("you are not connected to any room (cmd: ~ join/create {room_name})"))
			}

			c.Server.BroadcastChannel <- msg
			continue
		}

		cmd := strings.Trim(msg.Content, "~ ")
		switch {
		case strings.HasPrefix(cmd, prefixCreate):
			name := strings.TrimPrefix(cmd, prefixCreate)
			if len(strings.TrimSpace(name)) == 0 {
				c.Send(*c.Server.message("~ create {room_name}"))
				continue
			}

			roomId, err := c.CreateRoom(name)
			if err != nil {
				// TODO add logging
				c.Send(*c.Server.message(fmt.Sprintf("cannot create room %q", err)))
				continue
			}

			c.Send(*c.Server.message("").WithMeta(map[MetaKey]any{
				MkRoomId: roomId,
			}))
		case strings.HasPrefix(cmd, prefixJoin):
			roomId := strings.TrimSpace(strings.TrimPrefix(cmd, prefixJoin))
			if len(strings.TrimSpace(roomId)) == 0 {
				c.Send(*c.Server.message("~ join {room_name}"))
				continue
			}

			err = c.JoinRoom(strings.TrimSpace(roomId))
			if err != nil {
				if errors.Is(err, ErrRoomNotExist) {
					c.Send(*c.Server.message("room does not exist"))
					continue
				}

				c.Send(*c.Server.message("cannot join room"))
				fmt.Printf("cannot join room %q\n", err)
			}

			c.Send(*c.Server.message("").WithMeta(map[MetaKey]any{
				MkRoomId: roomId,
			}))
		default:
			c.Send(*c.Server.message("unrecognized command"))
		}
	}
}

func (c *Client) String() string {
	return fmt.Sprintf("%s (%d)", c.Name, c.Id)
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

func (c *Client) CreateRoom(name string) (roomId string, err error) {
	err = validate.Var(name, "required,min=3")
	if err != nil {
		return "", utils.NewAppErr("room name has to be at least 3 characters long")
	}

	r := repo.RoomsRepo{
		Db: database.Db,
	}

	roomId, err = r.Create(name, c.Id)
	if err != nil {
		return "", err
	}

	room := &Room{
		Id:   roomId,
		Name: name,
	}
	room.addMember(c)

	c.ActiveRoom = room
	c.Server.Rooms[roomId] = room

	return
}

func (c *Client) JoinRoom(roomId string) error {
	if c.Server.Rooms[roomId] == nil {
		return ErrRoomNotExist
	}

	c.ActiveRoom = c.Server.Rooms[roomId]
	c.Server.Rooms[roomId].addMember(c)

	return nil
}

func (c *Client) message(content string) *Message {
	if c.ActiveRoom == nil {
		panic("you should not send a message without an active room")
	}

	return &Message{
		SenderId: c.Id,
		RoomId:   c.ActiveRoom.Id,
		Content:  content,
	}
}

// TODO
// customize colorscheme
// fetch messages on login
// fetch rooms on server init
