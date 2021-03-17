package api

import (
	"net/http"
	"os"
	"fmt"
	"log"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"StreamDownloader/dl_convert"
	"StreamDownloader/env"
)

func Youtube(channelName string) {
	fmt.Println("Search for " + channelName + " livestreaming on Youtube")

	developerKey, _, _ := env.Load_env()

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

	searchLiveCall := service.Search.List([]string{"id", "snippet"}).ChannelId(channelId).Type("video").EventType("Live").MaxResults(1)
	searchLiveResp, err := searchLiveCall.Do()
	if err != nil {
		log.Fatal("error Receive Response from Youtube API(LiveUrl): %w", err)
	}

	if searchLiveResp.PageInfo.TotalResults == 0 {
		fmt.Println(channelName + " has no live streaming")
		os.Exit(0)
	} else {
		var streamId = searchLiveResp.Items[0].Id.VideoId
		var liveUrl = "https://www.youtube.com/watch?v=" + streamId
		fmt.Println("-> " + searchLiveResp.Items[0].Snippet.Title)
		dl_convert.LiveStream_dl(liveUrl, streamId)
	}
}
