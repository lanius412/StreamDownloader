package api

import (
	"bufio"
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"StreamDownloader/dl"
)

//go:embed keys/key_twitch.txt
var f embed.FS

type Stream struct {
	Data []struct {
		StreamId string `json:"id"`
		Title    string `json:"title"`
	} `json:"data"`
}

func Twitch(channelName string) {
	log.Println("Search for " + channelName + " livestreaming on Twitch")

	clientId, token := load_env()
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
		log.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	s := Stream{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Println(err)
	}

	var liveUrl = "https://www.twitch.tv/" + channelName

	if len(s.Data) == 0 {
		log.Println(channelName + " has no live streaming")
	} else {
		log.Print(s.Data[0].Title)
		date := time.Now().Format("2006-01-02 15_04")
		streamId := s.Data[0].StreamId
		dl.LiveStream_dl(liveUrl, date+"-"+streamId)
	}
}

func load_env() (clientId string, token string) {
	file, err := f.Open("keys/key_twitch.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	keys := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	clientId = keys[0][strings.LastIndex(keys[0], " ")+1:]
	token = keys[1][strings.LastIndex(keys[1], " ")+1:]

	return clientId, token
}
