package streamdeckvoicemeeter

import (
	"github.com/cliffrowley/go-streamdeck"
)

type contextMap map[string]interface{}
type actionContextMap map[string]contextMap

func (m actionContextMap) add(an string, cn string, a interface{}) {
	if _, ok := m[an]; !ok {
		m[an] = make(map[string]interface{}, 0)
	}
	m[an][cn] = a
}

func (m actionContextMap) find(an string, cn string) (interface{}, bool) {
	if _, ok := m[an]; ok {
		if _, ok := m[an][cn]; ok {
			return m[an][cn], true
		}
	}
	return nil, false
}

func (m actionContextMap) remove(an string, cn string) {
	if _, ok := m.find(an, cn); ok {
		delete(m[an], cn)
		if len(m[an]) == 0 {
			delete(m, an)
		}
	}
}

var (
	actions = make(actionContextMap, 0)
)

// Run runs the plugin.
func Run() error {
	client, err := streamdeck.Connect()
	if err != nil {
		return err
	}

	client.HandleWillAppearFunc(func(e *streamdeck.WillAppearEvent) {
		switch e.Action {
		case ToggleUUID:
			a := &ToggleAction{Client: client, Context: e.Context}
			actions.add(e.Action, e.Context, a)
			a.WillAppear(e)
		}
	})

	client.HandleWillDisappearFunc(func(e *streamdeck.WillDisappearEvent) {
		if a, ok := actions.find(e.Action, e.Context); ok {
			if ca, ok := a.(streamdeck.WillDisappearHandler); ok {
				ca.WillDisappear(e)
			}
		}
	})

	client.HandleKeyUpFunc(func(e *streamdeck.KeyUpEvent) {
		if a, ok := actions.find(e.Action, e.Context); ok {
			if ca, ok := a.(streamdeck.KeyUpHandler); ok {
				ca.KeyUp(e)
			}
		}
	})

	return client.Run()
}
