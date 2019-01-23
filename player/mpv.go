package player

import (
	"encoding/json"
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

type mpvInfo struct{}

func (i *mpvInfo) Title() string {
	return ""
}

func (i *mpvInfo) Thumbnail() string {
	return ""
}

func (i *mpvInfo) Volume() float64 {
	return 0.0
}

func (i *mpvInfo) Position() float64 {
	return 0.0
}

func (i *mpvInfo) State() string {
	return ""
}

type Property uint8

const (
	_ Property = iota
	PropPause
)

func mapPropName(property Property) string {
	var name string
	switch property {
	case PropPause:
		name = "pause"
	}

	return name
}

type Event struct {
	Name string
	Data interface{}
}

type EventHandler func(Event)

type mpvCommand struct {
	Command []interface{} `json:"command"`
}

type MPVPlayer struct {
	socketPath     string
	socketConn     net.Conn
	commandChannel chan mpvCommand
	eventHandlers  map[Property]EventHandler
}

func MPVNew() (*MPVPlayer, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	socketPath := filepath.Join(tmpDir, socketName)
	cmd := exec.Command("mpv",
		"--idle=yes",
		"--fullscreen",
		fmt.Sprintf("--input-ipc-server=%s", socketPath),
		"--ytdl-format=best",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	conn, err := tryDial(socketPath, 5, time.Second)
	if err != nil {
		return nil, err
	}

	p := MPVPlayer{
		socketPath:     socketPath,
		socketConn:     conn,
		commandChannel: make(chan mpvCommand),
		eventHandlers:  make(map[Property]EventHandler),
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

func (p *MPVPlayer) Close() {
	close(p.commandChannel)
	p.socketConn.Close()
}

// Handle sets handler to be a callback function when property changes
// Full list of valid properties:
// https://github.com/mpv-player/mpv/blob/master/DOCS/man/input.rst#property-list
func (p *MPVPlayer) Handle(property Property, handler EventHandler) {
	p.eventHandlers[property] = handler
}

func (p *MPVPlayer) Start() {
	go func() {
		for command := range p.commandChannel {
			fmt.Println("got", command)
			json.NewEncoder(p.socketConn).Encode(&command)
		}
	}()

	// Disable events
	p.send("disable_event", "all")

	// Set observers
	for id := range p.eventHandlers {
		p.send("observe_property", id, mapPropName(id))
	}

	response := make(map[string]interface{})
	for {
		json.NewDecoder(p.socketConn).Decode(&response)

		fmt.Println("response", response)
		if event, ok := response["event"]; ok && event == "property-change" {
			handler := p.eventHandlers[Property(response["id"].(float64))]
			handler(Event{
				Name: response["name"].(string),
				Data: response["data"],
			})
		}

		for k := range response {
			delete(response, k)
		}
	}
}

func (p *MPVPlayer) send(cmd string, params ...interface{}) {
	p.commandChannel <- mpvCommand{
		Command: append([]interface{}{cmd}, params...),
	}
}

func (p *MPVPlayer) Play() error {
	p.send("set_property", "pause", false)
	return nil
}

func (p *MPVPlayer) Pause() error {
	p.send("set_property", "pause", true)
	return nil
}

func (p *MPVPlayer) Stop() error {
	p.send("stop")
	return nil
}

func (p *MPVPlayer) Next() error {
	p.send("playlist-next", "force")
	return nil
}

func (p *MPVPlayer) Add(url string) error {
	p.send("loadfile", url, "append-play")
	return nil
}

func (p *MPVPlayer) Seek(position float64) error {
	fmt.Println("SEEEK", position*100, position)
	p.send("set_property", "percent-pos", position*100)
	return nil
}

func (p *MPVPlayer) InfoListener() <-chan Info {
	// TODO!
	return nil
}
