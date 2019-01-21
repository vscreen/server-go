package player

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"path/filepath"
	"time"
)

type mpvCommand struct {
	Command []interface{} `json:"command"`
}

type Player struct {
	socketPath     string
	socketConn     net.Conn
	commandChannel chan mpvCommand
	eventHandlers  map[Property]EventHandler
}

func New() (*Player, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	socketPath := filepath.Join(tmpDir, socketName)
	cmd := exec.Command("mpv",
		"--really-quiet",
		"--idle=yes",
		"--fullscreen",
		fmt.Sprintf("--input-ipc-server=%s", socketPath),
		"--ytdl-format=best",
	)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	conn, err := tryDial(socketPath, 5, time.Second)
	if err != nil {
		return nil, err
	}

	p := Player{
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

func (p *Player) Close() {
	close(p.commandChannel)
	p.socketConn.Close()
}

// Handle sets handler to be a callback function when property changes
// Full list of valid properties:
// https://github.com/mpv-player/mpv/blob/master/DOCS/man/input.rst#property-list
func (p *Player) Handle(property Property, handler EventHandler) {
	p.eventHandlers[property] = handler
}

func (p *Player) Start() {
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

func (p *Player) send(cmd string, params ...interface{}) {
	p.commandChannel <- mpvCommand{
		Command: append([]interface{}{cmd}, params...),
	}
}

func (p *Player) Add(url string) {
	p.send("loadfile", url, "append-play")
}

func (p *Player) Play() {
	p.send("set_property", "pause", false)
}

func (p *Player) Pause() {
	p.send("set_property", "pause", true)
}
