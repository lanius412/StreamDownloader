package dl_convert

import (
	"os"
	"strings"
	"log"
)

func search_file(newDir, fileId string) string {
	var baseFile string

	files, err := os.ReadDir(newDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), fileId) {
			baseFile = file.Name()
		}
	}

	return baseFile
}
