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
