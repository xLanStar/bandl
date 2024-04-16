package bandl

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type ITrackSource interface {
	FormatContent() string
	Name() string
}
type DownloadLog struct {
	Hash   string
	Name   string
	Source ITrackSource
}

func (l *DownloadLog) Format() string {
	return l.Hash + "," + l.Name + "," + l.Source.Name() + ":" + l.Source.FormatContent()
}

var (
	trackSourceParses = make(map[string]func(sourceContent string) ITrackSource)
)

func RegisterTrackSource(name string, parse func(sourceContent string) ITrackSource) {
	log.Println("RegisterTrackSource", name)
	trackSourceParses[name] = parse
}

func ParseDownloadLog(line string) DownloadLog {
	strs := strings.Split(line, ",")
	source := strs[2]
	sourceTypeLength := strings.Index(source, ":")
	sourceType := source[:sourceTypeLength]
	if _, ok := trackSourceParses[sourceType]; !ok {
		log.Println("Unknown source type", sourceType)
		return DownloadLog{}
	}

	return DownloadLog{
		Hash:   strs[0],
		Name:   strs[1],
		Source: trackSourceParses[sourceType](source[sourceTypeLength+1:]),
	}
}

func ReadDownloadLogFile(path string) []DownloadLog {
	downloadLogs := make([]DownloadLog, 0)

	if _, err := os.Stat(path); err == nil {
		f, _ := os.OpenFile(path, os.O_RDWR, 0666)
		defer f.Close()

		reader := bufio.NewReader(f)
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			downloadLog := ParseDownloadLog(string(line))
			log.Println("Read", downloadLog)
			downloadLogs = append(downloadLogs, downloadLog)
		}
	}

	return downloadLogs
}
