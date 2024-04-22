package bandl

import (
	"bufio"
	"log"
	"os"

	"github.com/anacrolix/torrent"
)

type DownloaderConfig struct {
	DownloadFile   string `yaml:"download_file"`
	DownloadFolder string `yaml:"download_folder"`
}

type OnDownloadFunc func()

type Downloader struct {
	client          *torrent.Client
	Config          *DownloaderConfig
	Trackers        []ITracker
	DownloadLogs    []DownloadLog
	downloadLogsMap map[string]*DownloadLog
	saved           bool
}

func (d *Downloader) AddTracker(tracker ITracker) {
	d.Trackers = append(d.Trackers, tracker)
}

func (d *Downloader) AddTorrentFromTrackResult(trackResult TrackResult) (*torrent.Torrent, error) {
	mi, err := GetMetaDataFromFileUrl(trackResult.Url)
	if err != nil {
		return nil, err
	}

	hash := mi.HashInfoBytes().HexString()
	downloadLog, downloaded := d.downloadLogsMap[hash]
	if downloaded {
		log.Println("Skip", hash, downloadLog.Name)
		return nil, nil
	}

	t, err := d.client.AddTorrent(mi)
	if err != nil {
		return nil, err
	}

	<-t.GotInfo()
	log.Println("Download", hash, t.Name())
	t.DownloadAll()
	go func() {
		defer t.Drop()
		<-t.Complete.On()
		log.Println("Download complete", hash, t.Name())
		d.DownloadLogs = append(d.DownloadLogs, DownloadLog{
			Hash:   hash,
			Name:   t.Name(),
			Source: trackResult.Source,
		})
		d.downloadLogsMap[hash] = &d.DownloadLogs[len(d.DownloadLogs)-1]
		d.saved = false
	}()
	return t, nil
}

func (d *Downloader) Track() {
	log.Println("Tracking...")
	for _, tracker := range d.Trackers {
		results := tracker.Track()
		log.Println("Found", len(results), "track results")
		for _, result := range results {
			d.AddTorrentFromTrackResult(result)
		}
	}
}

func (d *Downloader) Save() {
	log.Println("Save downloader...")
	if d.saved {
		return
	}
	f, _ := os.OpenFile(d.Config.DownloadFile, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	writer := bufio.NewWriter(f)
	for _, downloadLog := range d.DownloadLogs {
		log.Println("Save", downloadLog.Format())
		writer.WriteString(downloadLog.Format())
		writer.WriteRune('\n')
	}
	writer.Flush()
	d.saved = true
}

func (d *Downloader) Close() {
	log.Println("Close downloader...")
	d.client.Close()
	d.Save()
}

func NewDownloader(config *DownloaderConfig) *Downloader {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = config.DownloadFolder

	if _, err := os.Stat(config.DownloadFolder); err != nil {
		log.Printf("%s not found, create a new one.\n", config.DownloadFolder)
		os.Mkdir(config.DownloadFolder, 0755)
	}

	client, _ := torrent.NewClient(cfg)

	downloadLogs := ReadDownloadLogFile(config.DownloadFile)
	downloadLogMap := make(map[string]*DownloadLog, len(downloadLogs))
	for _, downloadLog := range downloadLogs {
		downloadLogMap[downloadLog.Hash] = &downloadLog
	}

	return &Downloader{
		client:          client,
		Trackers:        make([]ITracker, 0),
		Config:          config,
		DownloadLogs:    downloadLogs,
		downloadLogsMap: downloadLogMap,
		saved:           true,
	}
}
