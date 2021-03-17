package dl_convert

import (
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"strings"
	"fmt"
	"log"

	"StreamDownloader/checkLive"
)

var isRunEncode bool //true -> not encode

func IsRunEncode(tsFlag bool) {
	isRunEncode = tsFlag
}

func LiveStream_dl(liveUrl, streamId string) {
	fmt.Println("Download Start...")

	dlDir := set_directory(liveUrl)

	cTime := time.Now()
	ytdlCmd := exec.Command("youtube-dl", "--hls-use-mpegts", liveUrl)
	if ytdlErr := ytdlCmd.Start(); ytdlErr != nil {
		log.Fatal("error youtube-dl: ", ytdlErr)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if liveStatus := checkLive.IsLive(liveUrl); !liveStatus {
				waitErr := ytdlCmd.Wait()
				if waitErr != nil {
					log.Fatal("error youtube-dl wait :", waitErr)
				}
				fmt.Println("Stream End")

				goto GOTO_CONVERT
			}
		case <-sig:
			fmt.Println("Receive SIGINT")
			fmt.Println("Download Stop")

			goto GOTO_CONVERT
		}
	}

GOTO_CONVERT:
	fmt.Println("Search for Downloaded Video")
	var baseFile string
	for addMin := -1; addMin < 2; addMin++ {
		date := cTime.Add(time.Duration(addMin) * time.Minute).Format("2006-01-02 15_04")
		baseFile = search_file(dlDir, date+"-"+streamId)
		if baseFile != "" {
			fmt.Println("-> " + baseFile)
			break
		}

	}
	tsFile, mp4File := rename(baseFile)
	if !isRunEncode {
		ffmpeg(tsFile, mp4File)
		remove(tsFile)
		fmt.Println("Convert Successfully")
	}
	fmt.Println("Completed")
}

func set_directory(liveUrl string) string {
	pltform := [...] string{"youtube", "twitch"}

	cDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	
	var dlDir string
	for _, p := range pltform {
		if strings.Contains(liveUrl, p) {
			dlDir = filepath.Join(cDir, "DownloadedVideos", p)
		}
	}

	if err := os.MkdirAll(dlDir, 0777); err != nil {
		log.Fatal(err)
	}

	os.Chdir(dlDir)

	return dlDir
}
