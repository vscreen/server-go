package player

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ytdlPath string
)

// VideoInfo is a data structure extracted info from the given url
type VideoInfo struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	Duration  int64  `json:"duration"`
}

func extract(url string) (*VideoInfo, error) {
	var buff bytes.Buffer
	var info VideoInfo

	cmd := exec.Command("python", ytdlPath, "-fbest", "-j", url)
	cmd.Stdout = &buff
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(&buff).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func init() {
	log.Info("[extractor] cheking for a new update for youtube-dl")
	var updated bool

	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal("[extractor]", curDir)
	}

	ytdlPath = filepath.Join(curDir, "youtube-dl")

	// check youtube-dl
	if info, err := os.Stat(ytdlPath); err == nil {
		modY, modM, modD := info.ModTime().Date()
		todY, todM, todD := time.Now().Date()
		if modY == todY || modM == todM || modD == todD {
			updated = true
		}
	}

	if !updated {
		resp, err := http.Get("https://yt-dl.org/downloads/latest/youtube-dl")
		if err != nil {
			log.Fatal("[extractor] failed to update youtube-dl")
		}
		defer resp.Body.Close()

		f, err := os.Create(ytdlPath)
		if err != nil {
			log.Fatal("[extractor]", err)
		}
		defer f.Close()

		_, err = io.CopyBuffer(f, resp.Body, nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("[extractor] updated youtube-dl")
	} else {
		log.Info("[extractor] youtube-dl is up to date already")
	}

}
