package extractor

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

// VideoInfo is a data structure extracted info from the given url
type VideoInfo struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	Duration  int64  `json:"duration"`
}

var (
	ytdlIn  io.WriteCloser
	ytdlOut *bufio.Scanner
)

// Init creates a connection to youtube-dl extension
func Init() {
	var err error

	log.Info("[extractor] cheking for a new update for youtube-dl")
	var updated bool

	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal("[extractor]", curDir)
	}

	ytdlPath := filepath.Join(curDir, "extractor.pyz")

	// check youtube-dl
	if info, err := os.Stat(ytdlPath); err == nil {
		modY, modM, modD := info.ModTime().Date()
		todY, todM, todD := time.Now().Date()
		if modY == todY || modM == todM || modD == todD {
			updated = true
		}
	}

	if !updated {
		resp, err := http.Get("https://api.github.com/repos/vscreen/extractor/releases/latest")
		if err != nil {
			log.Fatal("[extractor] failed to request the latest asset")
		}
		defer resp.Body.Close()

		data := make(map[string]interface{})
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatal("[extractor] failed to get the latest asset")
		}
		latestAsset := data["assets"].([]interface{})[0].(map[string]interface{})
		url := latestAsset["browser_download_url"].(string)

		resp, err = http.Get(url)
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

	// start youtube-dl service
	cmd := exec.Command(
		"python",
		ytdlPath,
	)

	ytdlIn, err = cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	ytdlOut = bufio.NewScanner(out)
}

func Extract(url string) (VideoInfo, error) {
	fmt.Fprintln(ytdlIn, url)
	if !ytdlOut.Scan() {
		return VideoInfo{}, errors.New("youtube-dl has stopped")
	}
	var info VideoInfo
	json.Unmarshal(ytdlOut.Bytes(), &info)
	return info, nil
}
