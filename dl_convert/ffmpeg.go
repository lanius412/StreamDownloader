package dl_convert

import (
	"os/exec"
	"strconv"
	"strings"
	"fmt"
	"log"
)

func ffmpeg(tsFile, mp4File string) {
	fmt.Println("Convert Start...")
	ffmpegCmd := exec.Command("ffmpeg", "-i", tsFile, mp4File)
	stdout, _ := ffmpegCmd.StderrPipe()
	ffmpegErr := ffmpegCmd.Start()
	if ffmpegErr != nil {
		log.Fatal("error ffmpeg convert from .ts to .mp4", ffmpegErr)
	}

	buf := make([]byte, 8)
	for {
		_, err := stdout.Read(buf)
		if err != nil {
			break
		}
		allRes += string(buf)
		getRatio(allRes)
	}
	fmt.Println()
	ffmpegCmd.Wait()
}

var (
	duration = 0
	allRes   = ""
	lastPer  = -1
)

func getRatio(res string) {
	dIdx := strings.LastIndex(res, "Duration")
	if dIdx != -1 {
		if len(res) < dIdx+18 {
			return 
		}
		dur := res[dIdx+10:]
		if len(dur) > 8 {
			dur = dur[0:8]
			duration = durToSec(dur)
			fmt.Println("duration(seconds):", duration)

		}
		allRes = ""
	}
	if duration == 0 {
		return
	}
	tIdx := strings.LastIndex(res, "time=")
	if tIdx != -1 {
		if len(res) < tIdx+13 {
			return 
		}
		time := res[tIdx+5:]
		if len(time) > 8 {
			time = time[0:8]
			sec := durToSec(time)
			per := (sec * 100) / duration
			if lastPer != per {
				lastPer = per
				if per%2 == 0 {
					fmt.Print("#")
				}
			}
			allRes = ""
		}
	}
}

func durToSec(dur string) int {
	var sec int

	durAry := strings.Split(dur, ":")
	if len(durAry) != 3 {
		return sec
	}
	hours, _ := strconv.Atoi(durAry[0])
	sec = hours * (60 * 60)
	minutes, _ := strconv.Atoi(durAry[1])
	sec += minutes * 60
	seconds, _ := strconv.Atoi(durAry[2])
	sec += seconds
	
	return sec
}
