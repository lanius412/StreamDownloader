package dl

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

func LiveStream_dl(videoUrl string, streamId string) {
	log.Println("Download Start...")

	cDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	newDir := filepath.Join(cDir, "DownloadedVideos")
	os.Chdir(newDir)

	ytdlCmd := exec.Command("youtube-dl", "--hls-use-mpegts", videoUrl)
	if ytdlErr := ytdlCmd.Start(); ytdlErr != nil {
		log.Fatal("error youtube-dl: ", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	if s := <-sig; s != nil {
		log.Println("Receive SIGINT")
		log.Println("Download Stop")
		baseFile := search_file(newDir, streamId)
		tsFile := baseFile[:strings.LastIndex(baseFile, ".mp4")] + ".ts"
		renameErr := os.Rename(baseFile, tsFile)
		if renameErr != nil {
			log.Fatal("error rename to .ts: ", renameErr)
		}
		mp4File := strings.Replace(tsFile, "ts", "mp4", 1)
		ffmpeg(tsFile, mp4File)
		remove(tsFile)
	}

	if err = ytdlCmd.Wait(); err == nil {
		log.Println("Download Complete")
		tsFile := search_file(newDir, streamId)
		mp4File := strings.Replace(tsFile, "ts", "mp4", 1)
		ffmpeg(tsFile, mp4File)
		remove(tsFile)
	}

	log.Println("Convert Successfully")
}

func search_file(newDir string, streamId string) (baseFile string) {
	log.Println("Search for file")

	files, err := os.ReadDir(newDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), streamId) {
			baseFile = file.Name()
		}
	}
	return baseFile
}

func ffmpeg(tsFile string, mp4File string) {
	log.Println("Convert Start...")
	ffmpegCmd := exec.Command("ffmpeg", "-i", tsFile, mp4File)
	ffmpegErr := ffmpegCmd.Start()
	if ffmpegErr != nil {
		log.Fatal("error ffmpeg convert from .ts to .mp4", ffmpegErr)
	}
	ffmpegCmd.Wait()
}

func remove(tsFile string) {
	removeErr := os.Remove(tsFile)
	if removeErr != nil {
		log.Fatal("error remove .ts: ", removeErr)
	}
}
