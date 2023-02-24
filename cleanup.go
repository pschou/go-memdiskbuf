package memdiskbuf

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if Debug {
				log.Println("Caught signal", sig)
			}
			for file, _ := range tmpFile {
				os.Remove(file)
			}
		}
	}()
}

var (
	// Toggle debug
	Debug        = false
	tmpFile      = make(map[string]struct{})
	tmpFileMutex sync.Mutex
)

// Create an is-used mark
func unuse(f string) {
	tmpFileMutex.Lock()
	defer tmpFileMutex.Unlock()
	delete(tmpFile, f)
}

// Create an is-used mark
func use(f string) {
	tmpFileMutex.Lock()
	defer tmpFileMutex.Unlock()
	tmpFile[f] = struct{}{}
}
