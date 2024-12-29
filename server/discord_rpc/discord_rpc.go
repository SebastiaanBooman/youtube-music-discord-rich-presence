package discordrpc

import (
	"YoutubeMusicRichPresence/song_data_types"
	"github.com/hugolgst/rich-go/client"
)

func Login() error {
	err := client.Login("1260982067676577886")
	if err != nil {
		return err
	}
	return nil
}

func Logout() error {
	client.Logout()
	return nil
}

func SetSongActivity(songDataInformation songdatatypes.SongPresenceInformation) {
	err := client.SetActivity(client.Activity{
		State:      songDataInformation.SongData.Artist,
		Details:    songDataInformation.SongData.Title,
		LargeImage: songDataInformation.SongData.ImageUrl,
		LargeText:  songDataInformation.SongData.Album,
		SmallImage: songDataInformation.SmallImageKey,
		SmallText:  songDataInformation.SmallText,
		Timestamps: &client.Timestamps{
			Start: &songDataInformation.StartTime,
			End:   &songDataInformation.EndTime,
		},
	})

	if err != nil {
		panic(err)
	}
}
