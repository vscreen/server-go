// +build !pi

package player

func New() (Player, error) {
	return MPVNew()
}
