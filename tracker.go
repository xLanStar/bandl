package bandl

import (
	"log"
	"os"
)

type TrackingItem string

type ITracker interface {
	Init()
	addTracked(TrackResult)
	hasTracked(TrackResult) bool
	LoadTrackFile() []TrackingItem
	TrackItem(TrackingItem) []TrackResult
}

type Tracker struct {
	TrackFilePath string
	Tracked       map[string]bool
}

func (t *Tracker) Init() {
	t.Tracked = make(map[string]bool)

	if _, err := os.Stat(t.TrackFilePath); err != nil {
		os.Create(t.TrackFilePath)
	}
}

func (t *Tracker) addTracked(trackResult TrackResult) {
	t.Tracked[trackResult.Url] = true
}

func (t *Tracker) hasTracked(trackResult TrackResult) bool {
	_, ok := t.Tracked[trackResult.Url]
	return ok
}

func (t *Tracker) LoadTrackFile() []TrackingItem {
	log.Println("[Tracker] Load track file", t.TrackFilePath)
	return nil
}

func (t *Tracker) TrackItem(TrackingItem) []TrackResult {
	return nil
}
