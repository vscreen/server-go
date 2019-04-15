package backend

// New creates a player interface
func New(player string) (Player, error) {
	var b Player
	var err error

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

	return b, nil
}
