package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/googollee/go-socket.io"
)

type post struct {
	Name string
	Text string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":5000"
	} else {
		port = ":" + port
	}

	server, err := socketio.NewServer(nil)
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

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(port, nil))
}
