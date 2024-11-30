package songdatatypes

import "time"

type SongData struct {
	//TODO: Change the value to bool in web extension
	Playing  bool
	Artist   string
	Title    string
	Album    string
	ImageUrl string
	TimeData string
}

type SongPresenceInformation struct {
	SongData      SongData
	SmallImageKey string
	SmallText     string
	StartTime     time.Time
	EndTime       time.Time
}

// Discord rich presence specification requires fields to be at least two characters long
func (songData *SongData) AppendNullCharacterToDataStrings() {
	songData.Artist = songData.Artist + "\x00"
	songData.Title = songData.Title + "\x00"

	// Not all YouTube music videos contain album metadata
	if len(songData.Album) == 1 {
		songData.Album = songData.Album + "\x00"
	}
}
