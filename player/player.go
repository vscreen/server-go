// +build !pi

package player

import "github.com/vscreen/server-go/player/mpv"

func New() (Player, error) {
	return mpv.New()
}
