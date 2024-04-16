package bandl

import (
	"net/http"

	"github.com/anacrolix/torrent/metainfo"
)

func GetMetaDataFromFileUrl(url string) (*metainfo.MetaInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	mi, err := metainfo.Load(resp.Body)
	if err != nil {
		return nil, err
	}

	return mi, nil
}
