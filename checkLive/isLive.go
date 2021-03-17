package checkLive

import (
	"net/http"
	"encoding/json"
	"io"
	"strings"
	"log"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"StreamDownloader/env"
)

type Tmp struct {
	Data []struct {
	} `json:"data"`
}

func IsLive(liveUrl string) bool {
	var isLive bool
	
	developerKey, clientId, token := env.Load_env()

	switch {
	case strings.Contains(liveUrl, "youtube"):

		client := &http.Client{
			Transport: &transport.APIKey{Key: developerKey},
		}
		service, err := youtube.New(client)
		if err != nil {
			log.Fatal("error Creating New Youtube Client: %w", err)
		}
		streamId := strings.Replace(liveUrl, "https://www.youtube.com/watch?v=", "", 1)
		isLiveCall := service.Videos.List([]string{"snippet"}).Id(streamId)
		isLiveResp, err := isLiveCall.Do()
		if err != nil {
			log.Fatal("error Receive Response from Youtube API(LiveStatus): ", err)
		}
		liveStatus := isLiveResp.Items[0].Snippet.LiveBroadcastContent //live, upcoming or none
		if liveStatus == "none" {
			isLive = false
		} else if liveStatus == "live" {
			isLive = true
		}
		break

	case strings.Contains(liveUrl, "twitch"):
		const apiUrl = "https://api.twitch.tv/helix/streams?user_login="

		channelName := strings.Replace(liveUrl, "https://www.twitch.tv/", "", 1)
		client := &http.Client{}
		req, err := http.NewRequest("GET", apiUrl+channelName, nil)
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("client-id", clientId)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		tmp := Tmp{}
		err = json.Unmarshal(body, &tmp)
		if err != nil {
			log.Println(err)
		}
		if len(tmp.Data) == 0 {
			isLive = false
		} else {
			isLive = true
		}
		break

	}

	return isLive
}
