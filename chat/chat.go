package chat

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/googollee/go-socket.io"
)

type post struct {
	Name string
	Text string
}

var server *socketio.Server

func init() {
	var err error
	server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	var users = make(map[string]string)

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Join("chat")

		so.On("new user", func(username string) {
			outputMsg := fmt.Sprintf("%v [ENTROU]", username)

			id := so.Id()

			users[id] = username

			so.BroadcastTo("chat", "new user", outputMsg)
		})

		so.On("chat message", func(msg string) {
			jsonBytes := []byte(msg)

			var posts post
			json.Unmarshal(jsonBytes, &posts)

			outputMsg := fmt.Sprintf("%v: %v", posts.Name, posts.Text)

			log.Println("emit:", so.Emit("chat message", outputMsg))

			so.BroadcastTo("chat", "chat message", outputMsg)
		})

		so.On("disconnection", func() {
			id := so.Id()

			username := users[id]

			outputMsg := fmt.Sprintf("%v [SAIU]", username)

			so.BroadcastTo("chat", "new user", outputMsg)
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
}
