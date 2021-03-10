package api

import (
	_ "embed"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"StreamDownloader/dl"
)

//go:embed keys/key_youtube.txt
var developerKey string

func Youtube(channelName string) {
	log.Println("Search for " + channelName + " livestreaming on Youtube")

	developerKey = developerKey[strings.LastIndex(developerKey, " ")+1:]

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatal("error Creating New Youtube Client: %w", err)
	}

	channelCall := service.Search.List([]string{"id"}).Type("channel").Q(channelName).MaxResults(1)
	channelResp, err := channelCall.Do()
	if err != nil {
		log.Fatal("error Receive Response from Youtube API(ChannelId): %w", err)
	}

	var channelId = channelResp.Items[0].Id.ChannelId

	liveCall := service.Search.List([]string{"id", "snippet"}).ChannelId(channelId).Type("video").EventType("Live").MaxResults(1)
	liveResp, err := liveCall.Do()
	if err != nil {
		log.Fatal("error Receive Response from Youtube API(LiveUrl): %w", err)
	}
	streamId := liveResp.Items[0].Id.VideoId

	var liveUrl = "https://www.youtube.com/watch?v=" + streamId

	snippet := liveResp.Items[0].Snippet
	if isLive := snippet.LiveBroadcastContent; isLive == "live" {
		log.Println(snippet.Title)
		date := time.Now().Format("2006-01-02 15_04")
		dl.LiveStream_dl(liveUrl, date+"-"+streamId)
	} else {
		log.Println(channelName + " has no live streaming")
	}

}
