package converter

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
	"sync"
)

const (
	lastOpenedKey = "lastOpened"
	maxLastOpened = 5
)

type Application struct {
	app    fyne.App
	window fyne.Window

	mu         sync.Mutex
	lastOpened []string
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Run() error {
	a.app = app.NewWithID(BundleId)

	a.loadLastOpened()

	NewGreetingWindow(a).Show()
	a.app.Run()

	return nil
}

func (a *Application) getApp() fyne.App {
	return a.app
}

func (a *Application) loadLastOpened() {
	a.mu.Lock()
	defer a.mu.Unlock()

	rawJson := a.app.Preferences().StringWithFallback(lastOpenedKey, "[]")
	if err := json.Unmarshal([]byte(rawJson), &a.lastOpened); err != nil {
		log.Println("failed to load last opened files:", err)
	}
}

func (a *Application) appendLastOpened(path string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	//check if we already have this path
	for _, p := range a.lastOpened {
		if p == path {
			return
		}
	}

	a.lastOpened = append(a.lastOpened, path)
	if len(a.lastOpened) > maxLastOpened {
		a.lastOpened = a.lastOpened[1:]
	}

	rawJson, err := json.Marshal(a.lastOpened)
	if err != nil {
		log.Println("failed to save last opened files:", err)
		return
	}

	a.app.Preferences().SetString(lastOpenedKey, string(rawJson))
}

func (a *Application) getLastOpened() []string {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.lastOpened
}
