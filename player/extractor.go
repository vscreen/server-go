package player

import (
	"bytes"
	"encoding/json"
	"os/exec"
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

	cmd := exec.Command("youtube-dl", "-fbest", "-j", url)
	cmd.Stdout = &buff
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(&buff).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
