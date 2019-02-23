package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cliffrowley/go-voicemeeter"
	"github.com/cliffrowley/streamdeck-voicemeeter/internal/pkg/streamdeckvoicemeeter"
)

var logfile *os.File

func enableLogging() {
	temp := os.Getenv("TEMP")
	if temp != "" {
		fh, err := os.OpenFile(filepath.Join(temp, "streamdeck-voicemeeter.txt"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0x755)
		if err != nil {
			log.Println("Failed enabling file logging")
		}
		log.SetOutput(fh)
		log.Println("Enabled file logging")
		logfile = fh
	}
}

func disableLogging() {
	logfile.Close()
}

func main() {
	// enableLogging()

	err := voicemeeter.InitLibrary()
	if err != nil {
		log.Fatalln("Error initializing Voicemeeter: ", err)
	}

	err = voicemeeter.Login()
	if err != nil {
		log.Fatalln("Error logging in to Voicemeeter: ", err)
	}

	defer func() {
		voicemeeter.Logout()
		voicemeeter.CleanupLibrary()
		// disableLogging()
	}()

	err = streamdeckvoicemeeter.Run()
	if err != nil {
		log.Fatalln("Error running plugin: ", err)
	}
}
