package streamdeckvoicemeeter

import (
	"encoding/json"
	"log"
	"time"

	"github.com/cliffrowley/go-voicemeeter"

	"github.com/cliffrowley/go-streamdeck"
)

// ToggleUUID contains the UUID of the Toggle action.
const ToggleUUID = "com.github.cliffrowley.streamdeck-voicemeeter.toggle"

const toggleTickInterval = 50 * time.Millisecond

type toggleActionSettings struct {
	ParameterName string  `json:"parameterName"`
	OnValue       float32 `json:"onValue"`
	OffValue      float32 `json:"offValue"`
	CurValue      float32 `json:"curValue"`
}

func (s *toggleActionSettings) unmarshal(data []byte) *toggleActionSettings {
	json.Unmarshal(data, s)
	return s
}

func (s *toggleActionSettings) marshal() []byte {
	d, _ := json.Marshal(s)
	return d
}

// ToggleAction encapsulates the toggle action.
type ToggleAction struct {
	Client  *streamdeck.Client
	Context string

	settingsCache *toggleActionSettings
	stop          chan bool
}

func (a *ToggleAction) startTicker() {
	a.stop = make(chan bool)
	ticker := time.NewTicker(toggleTickInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				a.tick()
			case <-a.stop:
				break
			}
		}
	}()
}

func (a *ToggleAction) stopTicker() {
	a.stop <- true
}

func (a *ToggleAction) tick() {
	val, err := GetParameterFloat(a.settingsCache.ParameterName)
	if err != nil {
		log.Println("Error getting value: ", err)
	}

	if a.settingsCache.CurValue != val {
		a.settingsCache.CurValue = val

		a.Client.SetSettings(a.Context, a.settingsCache.marshal())

		if a.settingsCache.CurValue == a.settingsCache.OnValue {
			log.Printf("Setting state to 1 (%v)\n", a.settingsCache.CurValue)
			a.Client.SetState(a.Context, 1)
		} else {
			log.Printf("Setting state to 0 (%v)\n", a.settingsCache.CurValue)
			a.Client.SetState(a.Context, 0)
		}
	}
}

// WillAppear handles the willAppear event.
func (a *ToggleAction) WillAppear(e *streamdeck.WillAppearEvent) {
	s := (&toggleActionSettings{}).unmarshal(e.Payload.Settings)

	// TODO remove these temporary defaults
	if s.ParameterName == "" {
		s.ParameterName = "Strip[0].A1"
		s.OnValue = 1
		s.OffValue = 0
		s.CurValue = 0

		a.Client.SetSettings(e.Context, s.marshal())
	}

	a.settingsCache = s

	if a.settingsCache.ParameterName != "" {
		a.startTicker()
	}
}

// WillDisappear handles the willDisappear event.
func (a *ToggleAction) WillDisappear(e *streamdeck.WillDisappearEvent) {
	a.stopTicker()
}

// KeyUp handles the keyUp event.
func (a *ToggleAction) KeyUp(e *streamdeck.KeyUpEvent) {
	s := (&toggleActionSettings{}).unmarshal(e.Payload.Settings)

	if s.ParameterName != "" {
		if s.CurValue != s.OnValue {
			voicemeeter.SetParameterFloat(s.ParameterName, s.OnValue)
		} else {
			voicemeeter.SetParameterFloat(s.ParameterName, s.OffValue)
		}
	}
}
