package bangumi

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/xLanStar/bandl"
)

const (
	bangumiUrl = "https://bangumi.moe/rss/tags/"
)

var (
	fp = gofeed.NewParser()
)

type BangumiTracker struct {
	bandl.Tracker
}

func (t *BangumiTracker) Init() {
	t.Tracker.Init()
}

func (t *BangumiTracker) LoadTrackFile() []bandl.TrackingItem {
	f, _ := os.OpenFile(t.TrackFilePath, os.O_RDONLY, 0666)
	defer f.Close()

	reader := bufio.NewReader(f)
	currentTags := make([]string, 0, 8)
	result := make([]bandl.TrackingItem, 0, 32)

	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line := strings.TrimSpace(string(bytes))
		if len(line) == 0 || line[0] == '#' {
			continue
		} else if line[0] == '+' {
			currentTags = append(currentTags, line[1:])
			continue
		} else if line[0] == '-' {
			currentTags = currentTags[:len(currentTags)-1]
			continue
		}

		tagsStr := strings.Join(currentTags, "+") + "+" + line
		result = append(result, bandl.TrackingItem(tagsStr))
	}

	return result
}

func (t *BangumiTracker) TrackItem(trackingItem bandl.TrackingItem) []bandl.TrackResult {
	rssFeed, _ := fp.ParseURL(bangumiUrl + string(trackingItem))
	result := make([]bandl.TrackResult, 0, len(rssFeed.Items))
	for _, item := range rssFeed.Items {
		log.Println("Found", item.Title)
		result = append(result, bandl.TrackResult{
			Url:    item.Enclosures[0].URL,
			Source: &BangumiTrackSource{Tags: string(trackingItem)},
		})
	}
	return result
}
