package bangumi

import (
	"github.com/xLanStar/bandl"
)

const (
	bangumiSourceName = "bangumi"
)

type BangumiTrackSource struct {
	Tags string
}

func (s *BangumiTrackSource) FormatContent() string {
	return s.Tags
}

func (BangumiTrackSource) Name() string {
	return bangumiSourceName
}

func ParseBangumiTrackSource(sourceContent string) bandl.ITrackSource {
	return &BangumiTrackSource{
		Tags: sourceContent,
	}
}

func init() {
	bandl.RegisterTrackSource(bangumiSourceName, ParseBangumiTrackSource)
}
