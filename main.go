package main

import (
	"StreamDownloader/api"
	"flag"
)

func main() {
	var (
		youtubeFlag = flag.String("yt", "", "channelName(download youtube stream)")
		twitchFlag  = flag.String("ttv", "", "channelName(download twitch stream")
	)
	flag.Parse()

	switch {
	case *youtubeFlag != "":
		api.Youtube(*youtubeFlag)
	case *twitchFlag != "":
		api.Twitch(*twitchFlag)
	}

}
