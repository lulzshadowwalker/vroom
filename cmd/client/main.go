package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/lulzshadowwalker/vroom/pkg/chat"
)

func main() {
	tcpAddress, err := net.ResolveTCPAddr("tcp4", "localhost:3000")
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

	go func() {
		buf := make([]byte, 512)
		for {
			n, err := con.Read(buf)
			if err != nil {
				fmt.Println("server hsa disconnected")
				return
			}

			fmt.Println(string(buf[:n]))
		}
	}()

	msg := chat.Message{
		RoomId:   "hello",
		SenderId: "Pepega",
	}
	for {
		fmt.Print("message: ")
		fmt.Scan(&msg.Content)
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot encode message %q", err)
			continue
		}

		// should be chat.Message
		_, err = con.Write([]byte(jsonMsg))
		if err != nil {
			fmt.Printf("canont write message to the server %q\n", err)
		}
	}
}
