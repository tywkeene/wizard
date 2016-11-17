package logging

import (
	"log"
	"os"
)

func OpenLog(path string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	log.Println("=====================================")
	log.Printf("Opened log (%s)\n", path)
}
