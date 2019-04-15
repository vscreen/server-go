// +build !pi
// +build linux

package backend

import "fmt"

// <GLOBAL>
var Players = []string{MPV, VLC}

// <VLC>
const (
	vlcBin = "vlc"
)

var (
	vlcArgs []string
)

func init() {
	vlcArgs = []string{
		"--fullscreen",
		"--no-loop",
		"--no-osd",
		"--play-and-stop",
		"--intf=http",
		fmt.Sprintf("--http-password=%s", vlcPassword),
		fmt.Sprintf("--http-port=%d", vlcPort),
	}
}
