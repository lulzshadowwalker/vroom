package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/lulzshadowwalker/vroom/internal/database"
	"github.com/lulzshadowwalker/vroom/internal/database/model"
	"github.com/lulzshadowwalker/vroom/pkg/auth"
	"github.com/lulzshadowwalker/vroom/pkg/chat"
	"github.com/lulzshadowwalker/vroom/pkg/repo"
	"github.com/lulzshadowwalker/vroom/pkg/utils"
	"github.com/manifoldco/promptui"
)

func main() {
	tcpAddress, err := net.ResolveTCPAddr("tcp4", "localhost:3124")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot instantiate a server %q", err)
		os.Exit(1)
	}

	con, err := net.DialTCP("tcp4", nil, tcpAddress)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer con.Close()

	user, err := handleAuth()
	if err != nil {
		// TODO better auth ux
		fmt.Println(err)
		return
	}

	outMsg := chat.Message{
		SenderId:   user.Id,
		SenderName: user.Username,
	}

	rooms, err := getUserRooms(user.Id)
	if err != nil && len(rooms) != 0 {
		r, err := runSelectRoomPrompt(rooms)
		if err != nil {
			outMsg.Content = "~ join " + r.Id
			jsonMsg, err := json.Marshal(outMsg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot encode message %q", err)
			}

			_, err = con.Write(jsonMsg)
			if err != nil {
				fmt.Printf("cannot join room "+r.Name, err)
			}
		}
	}

	go func() {
		buf := make([]byte, 512)
		var msg chat.Message

		for {
			n, err := con.Read(buf[:])
			if err != nil {
				fmt.Println("server has disconnected")
				return
			}

			err = json.Unmarshal(buf[:n], &msg)
			if err != nil {
				fmt.Println("cannot recieve message", err, string(buf[:n]))
				continue
			}

			if msg.Meta != nil {
				roomId := msg.Meta[chat.MkRoomId]

				if roomId != nil {
					if id, ok := roomId.(string); ok {
						outMsg.RoomId = id
					}

					fmt.Println("room:", roomId)
				}
			}

			if len(strings.TrimSpace(msg.Content)) != 0 {
				fmt.Println(msg.Content)
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("cannot read input", err)
		}
		outMsg.Content = string(l)

		jsonMsg, err := json.Marshal(outMsg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot encode message %q", err)
			continue
		}

		if len(strings.Trim(outMsg.Content, " ")) == 0 {
			continue
		}

		_, err = con.Write(jsonMsg)
		if err != nil {
			fmt.Printf("canont write message to the server %q\n", err)
		}
	}
}

func handleAuth() (*model.User, error) {
	authHandler := auth.AuthHandler{
		Repo: &repo.AuthRepo{
			Db: database.Db,
		},
	}

Retry:
	user, err := authHandler.Trigger()
	if err != nil {
		if _, ok := err.(*utils.AppErr); ok {
			fmt.Println(err)

			fmt.Println("try again")
			goto Retry
		} else {
			return nil, fmt.Errorf("cannot sign in %w", err)
			// TODO add logger
		}
	}
	return user, nil
}

func getUserRooms(userId int) ([]model.Room, error) {
	roomsRepo := repo.RoomsRepo{
		Db: database.Db,
	}

	return roomsRepo.GetUserRooms(userId)
}

func runSelectRoomPrompt(rooms []model.Room) (*model.Room, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F336 {{ .Name | yellow | cyan }}",
		Details: `
--------- Room ----------
{{ "Name:" | faint }}	{{ .Name }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := rooms[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Room",
		Items:     rooms,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, fmt.Errorf("cannot run room prompt %w", err)
	}

	return &rooms[i], nil
}
