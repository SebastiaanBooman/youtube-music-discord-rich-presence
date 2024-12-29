package main

import (
	"YoutubeMusicRichPresence/api"
	"net/http"
)

func main() {
	server := api.CreateServer()
	http.HandleFunc("/song-data", server.ReceiveSongData)
	err := http.ListenAndServe(
		":8080",
		nil)
	if err != nil {
		panic(err)
	}
}
