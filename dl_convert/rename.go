package dl_convert

import (
	"os"
	"strings"
	"log"
)

func rename(baseFile string) (string, string) {
	var tsFile = baseFile[:strings.LastIndex(baseFile, ".mp4")] + ".ts"

	renameErr := os.Rename(baseFile, tsFile)
	if renameErr != nil {
		log.Fatal("error rename to .ts: ", renameErr)
	}

	var mp4File = strings.Replace(tsFile, "ts", "mp4", 1)

	return tsFile, mp4File
}
