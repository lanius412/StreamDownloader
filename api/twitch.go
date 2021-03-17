package api

import (
	"net/http"
	"encoding/json"
	"os"
	"io"
	"fmt"
	"log"	

	"StreamDownloader/dl_convert"
	"StreamDownloader/env"
)

type Stream struct {
	Data []struct {
		StreamId string `json:"id"`
		Title    string `json:"title"`
	} `json:"data"`
}

func Twitch(channelName string) {
	fmt.Println("Search for " + channelName + " livestreaming on Twitch")

	_, clientId, token := env.Load_env()

	const url = "https://api.twitch.tv/helix/streams?user_login="

	client := &http.Client{}
	req, err := http.NewRequest("GET", url+channelName, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("client-id", clientId)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	s := Stream{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Fatal(err)
	}

	var liveUrl = "https://www.twitch.tv/" + channelName

	if len(s.Data) == 0 {
		fmt.Println(channelName + " has no live streaming")
		os.Exit(0)
	} else {
		fmt.Print("-> " + s.Data[0].Title)
		streamId := s.Data[0].StreamId
		dl_convert.LiveStream_dl(liveUrl, streamId)
	}
}
