// +build !pi
// +build darwin

package player

// <GLOBAL>
var Players = []string{MPV, VLC}

// <VLC>
const (
	vlcBin = "/Applications/VLC.app/Contents/MacOS/VLC"
)
