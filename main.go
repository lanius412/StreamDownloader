package main

import (
	"flag"

	"StreamDownloader/dl_convert"
	"StreamDownloader/api"
)

func main() {
	var (
		youtubeFlag = flag.String("yt", "", "channelName(download youtube stream)")
		twitchFlag  = flag.String("ttv", "", "channelName(download twitch stream")
		tsFlag = flag.Bool("ts", false, "output .ts(do not encode to mp4)")
	)
	flag.Parse()

	dl_convert.IsRunEncode(*tsFlag)

	switch {
	case *youtubeFlag != "":
		api.Youtube(*youtubeFlag)
	case *twitchFlag != "":
		api.Twitch(*twitchFlag)
	}

}
