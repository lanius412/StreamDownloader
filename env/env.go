package env

import (
	"embed"
	"bufio"
	"strings"
	"log"
	
)

//go:embed keys/key_youtube.txt
var developerKey string

//go:embed keys/key_twitch.txt
var f embed.FS

func Load_env() (string, string, string) {
	var youtubeDeveloperKey = developerKey[strings.LastIndex(developerKey, " ")+1:]

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

	var twitchClientId = keys[0][strings.LastIndex(keys[0], " ")+1:]
	var twitchToken = keys[1][strings.LastIndex(keys[1], " ")+1:]

	return youtubeDeveloperKey, twitchClientId, twitchToken
}
