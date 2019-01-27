package player

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	socketName = "vscreen.sock"
)

type mpvCommand struct {
	Command []interface{} `json:"command"`
}

// mpvPlayer is an abstraction on top of mpv
type mpvPlayer struct {
	socketPath string
	socketConn net.Conn
}

// mpvNew creates mpvPlayer instance
func mpvNew() (*mpvPlayer, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	socketPath := filepath.Join(tmpDir, socketName)
	cmd := exec.Command("mpv",
		"--msg-level=all=error",
		"--idle=yes",
		"--fullscreen",
		fmt.Sprintf("--input-ipc-server=%s", socketPath),
		"--ytdl-format=best",
	)

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	conn, err := tryDial(socketPath, 5, time.Second)
	if err != nil {
		return nil, err
	}

	p := mpvPlayer{
		socketPath: socketPath,
		socketConn: conn,
	}

	// Disable events
	if err := p.send("disable_event", "all"); err != nil {
		return nil, err
	}

	return &p, nil
}

func tryDial(socketPath string, trials int, d time.Duration) (conn net.Conn, err error) {
	for i := 0; i < trials; i++ {
		conn, err = net.Dial("unix", socketPath)
		if err == nil {
			return conn, nil
		}
		time.Sleep(d)
	}
	return nil, err
}

// Close cleans up mpvPlayer's resources
func (p *mpvPlayer) close() {
	p.socketConn.Close()
}

func (p *mpvPlayer) send(cmd string, params ...interface{}) error {
	command := mpvCommand{
		Command: append([]interface{}{cmd}, params...),
	}
	if err := json.NewEncoder(p.socketConn).Encode(&command); err != nil {
		return err
	}

	response := make(map[string]interface{})
	json.NewDecoder(p.socketConn).Decode(&response)

	reason := response["error"].(string)

	if reason != "success" {
		return errors.New("operation failed")
	}
	return nil
}

func (p *mpvPlayer) play() error {
	return p.send("set_property", "pause", false)
}

func (p *mpvPlayer) pause() error {
	return p.send("set_property", "pause", true)
}

func (p *mpvPlayer) add(url string) error {
	return p.send("loadfile", url, "replace")
}

func (p *mpvPlayer) seek(position float64) error {
	return p.send("set_property", "percent-pos", position*100)
}
