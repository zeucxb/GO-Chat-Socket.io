package main

import (
	"log"
	"net/http"
	"os"

	"GO-Chat-Socket.io/chat"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":5000"
	} else {
		port = ":" + port
	}

	http.Handle("/socket.io/", chat.Server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(port, nil))
}
