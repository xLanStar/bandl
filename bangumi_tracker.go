package bandl

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/mmcdole/gofeed"
)

const (
	bangumiUrl        = "https://bangumi.moe/rss/tags/"
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

func ParseBangumiTrackSource(sourceContent string) ITrackSource {
	return &BangumiTrackSource{
		Tags: sourceContent,
	}
}

func init() {
	RegisterTrackSource(bangumiSourceName, ParseBangumiTrackSource)
}

type BangumiTrackerConfig struct {
	TrackFile string `yaml:"track_file"`
}

type BangumiTracker struct {
	Config     *BangumiTrackerConfig
	Downloaded map[string]bool
}

func (t *BangumiTracker) Init() {
	if _, err := os.Stat(t.Config.TrackFile); err != nil {
		f, _ := os.Create(t.Config.TrackFile)
		f.Close()
	}
}

func (t *BangumiTracker) Track() []TrackResult {
	results := make([]TrackResult, 0, 8)

	f, _ := os.OpenFile(t.Config.TrackFile, os.O_RDONLY, 0666)
	defer f.Close()

	reader := bufio.NewReader(f)
	tags := make([]string, 0, 256)
	fp := gofeed.NewParser()

	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line := strings.TrimSpace(string(bytes))
		if len(line) == 0 || line[0] == '#' {
			continue
		} else if line[0] == '+' {
			tags = append(tags, line[1:])
			continue
		} else if line[0] == '-' {
			tags = tags[:len(tags)-1]
			continue
		}

		log.Println("Track", strings.Join(tags, "+")+"+"+line)
		tagsStr := strings.Join(tags, "+") + "+" + line
		rssFeed, _ := fp.ParseURL(bangumiUrl + tagsStr)
		for _, item := range rssFeed.Items {
			log.Println("Found", item.Title)
			results = append(results, TrackResult{
				Url:    item.Enclosures[0].URL,
				Source: &BangumiTrackSource{Tags: tagsStr},
			})
		}
	}

	return results
}
