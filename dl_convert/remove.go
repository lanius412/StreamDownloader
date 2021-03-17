package dl_convert

import (
	"os"
	"log"
)

func remove(tsFile string) {
	removeErr := os.Remove(tsFile)
	if removeErr != nil {
		log.Fatal("error remove .ts: ", removeErr)
	}
}
