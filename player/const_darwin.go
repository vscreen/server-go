// +build !pi
// +build darwin

package player

import "fmt"

// <GLOBAL>
var Players = []string{MPV, VLC}

// <VLC>
const (
	vlcBin = "/Applications/VLC.app/Contents/MacOS/VLC"
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
		"--intf=macosx",
		"--extraintf=http",
		fmt.Sprintf("--http-password=%s", vlcPassword),
		fmt.Sprintf("--http-port=%d", vlcPort),
	}
}
