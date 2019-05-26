package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/zserge/lorca"
)

//go:generate go run gen.go

// Player is a thin wrapper of JS calls
type Player struct {
	lorca.UI
	socket net.Listener
}

// NewPlayer creates a new player ui
func NewPlayer() (*Player, error) {
	ui, err := lorca.New("", "", 100, 100)
	if err != nil {
		return nil, err
	}

	b, err := ui.Bounds()
	if err != nil {
		return nil, err
	}
	b.WindowState = lorca.WindowStateFullscreen
	ui.SetBounds(b)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	return &Player{
		UI:     ui,
		socket: ln,
	}, nil
}

// Wait blocks execution until the ui closes
func (p *Player) Wait() {
	<-p.UI.Done()
	p.socket.Close()
}

// Close kills the ui
func (p *Player) Close() {
	p.UI.Close()
}

func (p *Player) call(name string, args ...interface{}) {
	var sb strings.Builder

	sb.WriteString("window.on.")
	sb.WriteString(name)
	sb.WriteString("(")

	argsArr := make([]string, 0, len(args))
	for _, arg := range args {
		val := fmt.Sprint(arg)
		if _, ok := arg.(string); ok {
			val = "\"" + val + "\""
		}
		argsArr = append(argsArr, val)
	}

	sb.WriteString(strings.Join(argsArr, ","))
	sb.WriteString(")")
	p.UI.Eval(sb.String())
}

// Play sends play command to the player
func (p *Player) Play() {
	p.call("play")
}

// Pause sends pause command to the player
func (p *Player) Pause() {
	p.call("pause")
}

// Stop sends stop command to the player
func (p *Player) Stop() {
	p.call("stop")
}

// Next sends next command to the player
func (p *Player) Next() {
	p.call("next")
}

// Add sends add command to the player
func (p *Player) Add(url string) {
	p.call("add", url)
}

// Seek sends seek command to the player
func (p *Player) Seek(position float64) {
	p.call("seek", position)
}
