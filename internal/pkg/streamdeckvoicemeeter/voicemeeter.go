package streamdeckvoicemeeter

import (
	"sync"
	"time"

	"github.com/cliffrowley/go-voicemeeter"
)

const dirtyInterval = 100 * time.Millisecond

var (
	dirtyLast time.Time
	dirtyLock sync.Mutex
	dirty     bool
)

// IsParametersDirty returns true if any parameters have been updated.
// Throttle and thread safe.
func IsParametersDirty() (bool, error) {
	dirtyLock.Lock()
	defer dirtyLock.Unlock()

	if time.Since(dirtyLast) >= dirtyInterval {
		d, err := voicemeeter.IsParametersDirty()
		if err != nil {
			return dirty, err
		}

		dirty = d
		dirtyLast = time.Now()
	}

	return dirty, nil
}

var (
	valuesCache = make(map[string]float32)
	valuesLock  sync.Mutex
)

// GetParameterFloat returns the float value of the given property.
// Thread safe.
func GetParameterFloat(name string) (float32, error) {
	dirty, err := IsParametersDirty()
	if err != nil {
		return 0, err
	}

	valuesLock.Lock()
	defer valuesLock.Unlock()

	if _, ok := valuesCache[name]; !ok || dirty {
		val, err := voicemeeter.GetParameterFloat(name)
		if err != nil {
			return 0, err
		}

		valuesCache[name] = val

		return val, nil
	}

	return valuesCache[name], nil
}

var (
	levelsLock  sync.Mutex
	levelsCache = make(map[voicemeeter.LevelType]map[uint32]float32, 0)
)

// GetLevel returns the level for the given type and channel.
func GetLevel(levelType voicemeeter.LevelType, channelIndex uint32) (float32, error) {
	levelsLock.Lock()
	defer levelsLock.Unlock()

	return voicemeeter.GetLevel(levelType, channelIndex)
}
