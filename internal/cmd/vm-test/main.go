package main

import (
	"log"
	"time"

	"github.com/cliffrowley/go-voicemeeter"
	"github.com/cliffrowley/streamdeck-voicemeeter/internal/pkg/streamdeckvoicemeeter"
)

func main() {
	log.Println("Initializing Voicemeeter..")
	err := voicemeeter.InitLibrary()
	if err != nil {
		log.Fatalln("Error initializing Voicemeeter: ", err)
	}

	log.Println("Logging into Voicemeeter..")
	err = voicemeeter.Login()
	if err != nil {
		log.Fatalln("Error logging into Voicemeeter: ", err)
	}

	defer func() {
		log.Println("Logging out of Voicemeeter..")
		err = voicemeeter.Logout()
		if err != nil {
			log.Fatalln("Error logging out of Voicemeeter: ", err)
		}

		log.Println("Cleaning up Voicemeeter..")
		err = voicemeeter.CleanupLibrary()
		if err != nil {
			log.Fatalln("Error cleaning up Voicemeeter: ", err)
		}
	}()

	for i := 0; i < 100; i++ {
		v, err := voicemeeter.GetParameterFloat("Strip[0].Gain")
		if err != nil {
			log.Fatalln("Error getting Voicemeeter parameter: ", err)
		}

		l, err := streamdeckvoicemeeter.GetLevel(voicemeeter.LevelInputPreFader, 0)
		if err != nil {
			log.Fatalln("Error getting Voicemeeter level: ", err)
		}

		log.Printf("[%v] -> %v=%v lvl=%v", i, "Strip[0].Gain", v, l)
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Done!")
}
