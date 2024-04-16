package bandl

type ITracker interface {
	Init()
	Track() []TrackResult
}

type TrackResult struct {
	Url    string
	Source ITrackSource
}
