package player

import (
	"github.com/vscreen/server-go/player/extractor"
)

// New creates a player interface
func New(player string) (*Player, error) {
	var b backendPlayer
	var err error

	extractor.Init()

	switch player {
	case MPV:
		b, err = mpvNew()
	case OMX:
		omxInit()
		b, err = omxNew()
	default: // vlc is default
		b, err = vlcNew()
	}

	if err != nil {
		return nil, err
	}

	return new(b)
}
