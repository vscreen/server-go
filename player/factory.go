// +build !pi

package player

const (
	// MPV (https://mpv.io) is a free, open source, and cross-platform media player
	MPV = "mpv"
	// VLC (https://www.videolan.org/vlc/index.html) is a free and open source
	// cross-platform multimedia player and framework that plays most multimedia
	// files as well as DVDs, Audio CDs, VCDs, and various streaming protocols.
	VLC = "vlc"
)

var Players = []string{MPV, VLC}

// New creates a player interface
func New(player string) (*Player, error) {
	var b backendPlayer
	var err error

	switch player {
	case MPV:
		b, err = mpvNew()
	default: // vlc is default
		b, err = vlcNew()
	}

	if err != nil {
		return nil, err
	}

	return new(b)
}
