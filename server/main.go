package main

import (
	"YoutubeMusicRichPresence/api"
	"net/http"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/song-data", server.ReceiveSongData)
	http.ListenAndServe(
		":8080",
		nil)
}
