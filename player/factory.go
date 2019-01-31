// +build !pi

package player

// New creates a player interface
func New() (*Player, error) {
	b, err := vlcNew()
	if err != nil {
		return nil, err
	}

	return new(b)
}
