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
	"sync"
	"time"
)

const (
	socketName = "vscreen.sock"
)

// Property is an id type to map mpv's properties
type Property uint8

const (
	_ Property = iota
	// PropPause is mpv's property
	PropPause
	// PropTitle is mpv's property
	PropTitle
)

func mapPropName(property Property) string {
	var name string
	switch property {
	case PropPause:
		name = "pause"
	case PropTitle:
		name = "media-title"
	}

	return name
}

// Event is a data structure to hold responses from mpv when state changes
type Event struct {
	Name string
	Data interface{}
}

// EventHandler is a signature to handle Event
type EventHandler func(Event)

type mpvCommand struct {
	Command []interface{} `json:"command"`
}

// MPVPlayer is an abstraction on top of mpv
type MPVPlayer struct {
	socketPath     string
	socketConn     net.Conn
	commandMutex   *sync.Mutex
	commandSuccess chan bool
	eventHandlers  map[Property]EventHandler
	infoChannel    chan Info
}

// MPVNew creates MPVPlayer instance
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
		commandMutex:   &sync.Mutex{},
		commandSuccess: make(chan bool),
		eventHandlers:  make(map[Property]EventHandler),
		infoChannel:    make(chan Info),
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

// Close cleans up MPVPlayer's resources
func (p *MPVPlayer) Close() {
	p.socketConn.Close()
}

// Handle sets handler to be a callback function when property changes
// Full list of valid properties:
// https://github.com/mpv-player/mpv/blob/master/DOCS/man/input.rst#property-list
func (p *MPVPlayer) handle(property Property, handler EventHandler) {
	p.eventHandlers[property] = handler
}

func (p *MPVPlayer) onPauseEvent(e Event) {
	// TODO!
}

func (p *MPVPlayer) onTitleEvent(e Event) {
	// TODO!
}

func (p *MPVPlayer) onTumbnailEvent(e Event) {
	// TODO!
}

// Start initializes player and block until done
func (p *MPVPlayer) Start() {
	// Set up property handlers for state management
	p.handle(PropPause, p.onPauseEvent)
	p.handle(PropTitle, p.onTitleEvent)

	// Disable events
	go p.send("disable_event", "all")

	// Set observers
	for id := range p.eventHandlers {
		go p.send("observe_property", id, mapPropName(id))
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
		} else if errorReason, ok := response["error"]; ok {
			// Handle error callback
			if errorReason == "success" {
				p.commandSuccess <- true
			} else {
				p.commandSuccess <- false
			}
		}

		for k := range response {
			delete(response, k)
		}
	}
}

func (p *MPVPlayer) send(cmd string, params ...interface{}) error {
	p.commandMutex.Lock()
	defer p.commandMutex.Unlock()

	command := mpvCommand{
		Command: append([]interface{}{cmd}, params...),
	}
	if err := json.NewEncoder(p.socketConn).Encode(&command); err != nil {
		return err
	}

	ok := <-p.commandSuccess
	if !ok {
		return errors.New("operation failed")
	}
	return nil
}

// Play plays the current video. If there's no video,
// nothing will happen
func (p *MPVPlayer) Play() error {
	return p.send("set_property", "pause", false)
}

// Pause pauses the current video. If there's no video,
// nothing will happen
func (p *MPVPlayer) Pause() error {
	return p.send("set_property", "pause", true)
}

// Stop stops the player and clear up the playlist
func (p *MPVPlayer) Stop() error {
	return p.send("stop")
}

// Next sets the player to the next video in the playlist
func (p *MPVPlayer) Next() error {
	return p.send("playlist-next", "force")
}

// Add adds the url to the playlist
func (p *MPVPlayer) Add(url string) error {
	return p.send("loadfile", url, "append-play")
}

// Seek sets the current video to position (from 0 - 1.0)
func (p *MPVPlayer) Seek(position float64) error {
	return p.send("set_property", "percent-pos", position*100)
}

// InfoListener returns a channel to get the most up to date info
func (p *MPVPlayer) InfoListener() <-chan Info {
	// TODO!
	return nil
}
