package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/xLanStar/bandl"
	"github.com/xLanStar/bandl/internal/tracker/bangumi"
)

var (
	downloader *bandl.Downloader
)

func init() {
	downloaderConfig := &bandl.DownloaderConfig{
		DownloadFolder: "downloads",
		DownloadFile:   "downloads.txt",
	}
	bandl.InitConfig("downloader.yaml", downloaderConfig)
	downloader = bandl.NewDownloader(downloaderConfig)

	bangumiTracker := &bangumi.BangumiTracker{
		Tracker: bandl.Tracker{
			TrackFilePath: "tracks-bangumi.txt",
		},
	}
	bangumiTracker.Init()
	log.Println("test", bangumiTracker.LoadTrackFile())
	downloader.AddTracker(bangumiTracker)
}

func main() {
	defer downloader.Close()

	timer := time.NewTimer(0)
	go func() {
		for {
			timer.Reset(time.Hour)
			<-timer.C
			downloader.Track()
		}
	}()

	go func() {
		i := -1
		for i != 0 {
			i = -1
			fmt.Println("Press any key to trace...")
			fmt.Scanln(&i)
			downloader.Track()
		}
	}()

	var quit chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Exiting...")
}
