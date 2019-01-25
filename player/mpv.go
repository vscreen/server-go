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
	playlist       []*VideoInfo
	infoCur        Info
	infoMutex      *sync.Mutex
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
		playlist:       make([]*VideoInfo, 0),
		infoCur:        Info{},
		infoMutex:      &sync.Mutex{},
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
	close(p.infoChannel)
}

// Handle sets handler to be a callback function when property changes
// Full list of valid properties:
// https://github.com/mpv-player/mpv/blob/master/DOCS/man/input.rst#property-list
func (p *MPVPlayer) handle(property Property, handler EventHandler) {
	p.eventHandlers[property] = handler
}

func (p *MPVPlayer) onPauseEvent(e Event) {
	p.updateInfo(func(info Info) Info {
		info.Playing = !e.Data.(bool)
		return info
	})
}

func (p *MPVPlayer) onTitleEvent(e Event) {
	if e.Data == nil {
		return
	}

	curVideo := p.playlist[0]
	p.playlist = p.playlist[1:]

	p.updateInfo(func(info Info) Info {
		info.Position = 0.0
		info.Thumbnail = curVideo.Thumbnail
		info.Title = curVideo.Title
		return info
	})
}

func (p *MPVPlayer) updateInfo(updater func(oldInfo Info) Info) {
	p.infoMutex.Lock()
	defer p.infoMutex.Unlock()

	newInfo := updater(p.infoCur)
	p.infoChannel <- newInfo
	p.infoCur = newInfo
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

	errorChan := make(chan string)
	eventChan := make(chan map[string]interface{})

	defer close(errorChan)
	defer close(eventChan)

	// Error handler
	go func() {
		// Handle error callback
		for reason := range errorChan {
			if reason == "success" {
				p.commandSuccess <- true
			} else {
				p.commandSuccess <- false
			}
		}
	}()

	// Event handler
	go func() {
		for response := range eventChan {
			handler := p.eventHandlers[Property(response["id"].(float64))]
			go handler(Event{
				Name: response["name"].(string),
				Data: response["data"],
			})
		}
	}()

	for {
		response := make(map[string]interface{})
		json.NewDecoder(p.socketConn).Decode(&response)

		if event, ok := response["event"]; ok && event == "property-change" {
			eventChan <- response
		} else if reason, ok := response["error"]; ok {
			errorChan <- reason.(string)
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
	if err := p.send("stop"); err != nil {
		return err
	}

	p.playlist = p.playlist[:0]
	p.updateInfo(func(info Info) Info {
		info.Position = 0.0
		info.Playing = false
		info.Title = ""
		info.Thumbnail = ""
		return info
	})
	return nil
}

// Next sets the player to the next video in the playlist
func (p *MPVPlayer) Next() error {
	return p.send("playlist-next", "force")
}

// Add adds the url to the playlist
func (p *MPVPlayer) Add(url string) error {
	info, err := extract(url)
	if err != nil {
		return err
	}

	p.playlist = append(p.playlist, info)
	return p.send("loadfile", info.URL, "append-play")
}

// Seek sets the current video to position (from 0 - 1.0)
func (p *MPVPlayer) Seek(position float64) error {
	if err := p.send("set_property", "percent-pos", position*100); err != nil {
		return err
	}

	p.updateInfo(func(info Info) Info {
		info.Position = position
		return info
	})
	return nil
}

// InfoListener returns a channel to get the most up to date info
func (p *MPVPlayer) InfoListener() <-chan Info {
	return p.infoChannel
}