package player

import (
	"github.com/vscreen/server-go/player/extractor"
)

// State represents the current state of the video player
type State struct {
	Title     string
	Thumbnail string
	Volume    float64
	Position  float64
	Playing   bool
	Playlist  []extractor.VideoInfo
	timer     *timer
}

func (s *State) reset() {
	t := s.timer // use old timer instead of creating a new one
	t.stop()
	*s = State{timer: t}
}

func (s *State) next() {
	playlist := s.Playlist
	if len(playlist) == 0 {
		s.reset()
		return
	}

	next := playlist[0]
	s.Playlist = playlist[1:]
	s.Position = 0.0
	s.Playing = true
	s.Title = next.Title
	s.Thumbnail = next.Thumbnail
	s.timer.new(next.Duration)
}

func (s *State) seek(position float64) {
	s.Position = position
	s.timer.seek(position)
}

func (s *State) play() {
	s.timer.play()
	s.Playing = true
}

func (s *State) pause() {
	s.timer.pause()
	s.Playing = false
	s.Position = s.timer.pos()
}
