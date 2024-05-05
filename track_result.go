package bandl

import "fmt"

type TrackResult struct {
	Url    string
	Source ITrackSource
}

func (r TrackResult) String() string {
	return fmt.Sprintf("TrackResult{Source: \033[33m%s\033[0m, Url: \033[36m%s\033[0m}", r.Source.Name(), r.Url)
}
