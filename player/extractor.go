package player

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ytdlPath string
)

// VideoInfo is a data structure extracted info from the given url
type VideoInfo struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	Duration  uint64 `json:"duration"`
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
	fmt.Println("[extractor] updating youtube-dl")
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal("[extractor]", curDir)
	}

	ytdlPath = filepath.Join(curDir, "youtube-dl")

	// check youtube-dl
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

	fmt.Println("[extractor] updated youtube-dl")
}
