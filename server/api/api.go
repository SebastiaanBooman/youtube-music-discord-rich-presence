package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"YoutubeMusicRichPresence/discord_rpc"
	"YoutubeMusicRichPresence/song_data_types"
)

const (
	PresenceKillTimeout = time.Second * 15
)

type Server struct {
	SongPresenceInformation songdatatypes.SongPresenceInformation
	UpdatePendingMutex      sync.Mutex
	SongDataMutex           sync.Mutex
	StopTimer               *time.Timer
	PresenceActive          bool
}

func CreateServer() *Server {
	return &Server{
		SongPresenceInformation: songdatatypes.SongPresenceInformation{},
	}
}

func (server *Server) ReceiveSongData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Expected a POST request", http.StatusBadRequest)
	}

	updateTime := time.Now()

	var songData songdatatypes.SongData
	// give the decode method a memory adress to place the extracted data
	if err := json.NewDecoder(r.Body).Decode(&songData); err != nil {
		http.Error(w, "Wrong JSON format", http.StatusBadRequest)
		return
	}

	elapsedSeconds, totalSeconds, err := dissectTimeData(songData.TimeData)
	if err != nil {
		http.Error(w, "Issues converting time data", http.StatusBadRequest)
		return
	}

	// This covers a limitation regarding the MutationObserver used in the web extension.
	// When automatically running new video (chaining videos), multiple events are fired
	// Due to this the extension sends multiple API calls in quick succession.
	// One of these may contain 00:00 / 00:00 as its time value.
	if elapsedSeconds == time.Duration(0) && totalSeconds == time.Duration(0) {
		http.Error(w, "Time is empty", http.StatusBadRequest)
		return
	}

	// Do not terminate the ipc connection if continued playing
	if server.StopTimer != nil {
		server.StopTimer.Stop()
		server.StopTimer = nil
	}

	if !server.PresenceActive {
		discordrpc.Login()
		server.PresenceActive = true
	}

	songData.AppendNullCharacterToDataStrings()

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	smallImageKey := "mozilla-firefox"
	smallText := "Mozilla Firefox"
	if !songData.Playing {
		smallImageKey = "play-button-icon"
		smallText = "paused"
	}

	var timeLeft time.Duration = totalSeconds - elapsedSeconds

	server.SongDataMutex.Lock()
	server.SongPresenceInformation = songdatatypes.SongPresenceInformation{
		SongData:      songData,
		SmallImageKey: smallImageKey,
		SmallText:     smallText,
		StartTime:     updateTime.Add(-elapsedSeconds),
		EndTime:       updateTime.Add(timeLeft),
	}
	server.SongDataMutex.Unlock()
	server.UpdateRichPresence()

	// Kill the rich presence if no new song has been playing for 15 seconds and paused
	if !songData.Playing {
		server.StopTimer = time.AfterFunc(PresenceKillTimeout, func() {
			server.PresenceActive = false
			discordrpc.Logout()
		})
	}

	fmt.Fprint(w, "Processed song data succesfully")
}

// Returns elapsed time and total song length, or
// -1 for both if error happened during conversion
func dissectTimeData(timeData string) (time.Duration, time.Duration, error) {
	var splitTime []string = strings.Split(timeData, "/")
	var elapsedSeconds string = strings.Trim(splitTime[0], " ")
	var totalSongSeconds string = strings.Trim(splitTime[1], " ")

	convertedElapsedSeconds, err := convertTimeDataToSeconds(elapsedSeconds)
	if err != nil {
		return -1, -1, err
	}
	convertedTotalSongSeconds, err := convertTimeDataToSeconds(totalSongSeconds)
	if err != nil {
		return -1, -1, err
	}
	return convertedElapsedSeconds, convertedTotalSongSeconds, nil
}

// Converts a string containing a time value delimited by colons
// e.g ("03:33:33" or "03:00") into a time.Duration
func convertTimeDataToSeconds(timeData string) (time.Duration, error) {
	conversionArray := [...]int{3600, 60, 1}
	var totalSeconds int
	var splitTimeData = strings.Split(timeData, ":")
	var startingIndex int = len(conversionArray) - len(splitTimeData)
	for i := 0; i < len(splitTimeData); i++ {
		// strconv.Atoi does not allow integers starting with 0 that are not 0
		timeValue := strings.TrimLeft(splitTimeData[i], "0")
		var convertedInt int
		if timeValue == "" {
			convertedInt = 0
			continue
		}
		convertedInt, err := strconv.Atoi(splitTimeData[i])
		if err != nil {
			return -1, err
		}
		totalSeconds += (convertedInt * conversionArray[i+startingIndex])
	}

	return time.Duration(totalSeconds) * time.Second, nil
}

func (server *Server) UpdateRichPresence() {
	// TODO: negating this does not acquire the lock ?
	if server.UpdatePendingMutex.TryLock() {
		server.SongDataMutex.Lock()
		discordrpc.SetSongActivity(server.SongPresenceInformation)
		server.SongDataMutex.Unlock()
		server.UpdatePendingMutex.Unlock()
	} else {
		return
	}
}
