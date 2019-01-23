// +build !pi

package player

// New creates a player interface
func New() (Player, error) {
	return MPVNew()
}
