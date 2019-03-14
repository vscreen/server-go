package player

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// backendPlayer is an abstraction of a video player. It contains
// some basic operations that should be available in most
// video player applications out there.
type backendPlayer interface {
	// These operations are thread safe
	play() error
	pause() error
	stop() error
	set(url string) error
	seek(position float64) error
	close()
}

// Player is an abstraction of a video player.
// the backend player is determined by the platform
type Player struct {
	State      *state
	b          backendPlayer
	playlist   []*videoInfo
	videoTimer *timer
	m          *sync.Mutex
}

func new(b backendPlayer) (*Player, error) {
	return &Player{
		State:      newState(),
		b:          b,
		playlist:   make([]*videoInfo, 0),
		videoTimer: nil,
		m:          &sync.Mutex{},
	}, nil
}

func (p *Player) onFinish() {
	p.m.Lock()
	defer p.m.Unlock()

	if len(p.playlist) == 0 {
		p.State.reset()
		p.videoTimer = nil
	} else {
		info := p.playlist[0]
		p.playlist = p.playlist[1:]

		// TODO! maybe todo something with error here
		p.b.set(info.URL)
		p.videoTimer = newTimer(info.Duration, p.onFinish)
		p.State.next(info.Title, info.Thumbnail)
	}
}

// Play plays the current video. If there's no video,
// nothing will happen
func (p *Player) Play() error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.videoTimer == nil {
		return nil
	}

	if err := p.b.play(); err != nil {
		return err
	}
	p.videoTimer.play()
	p.State.setPlaying(true)
	return nil
}

// Pause pauses the current video. If there's no video,
// nothing will happen
func (p *Player) Pause() error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.videoTimer == nil {
		return nil
	}

	if err := p.b.pause(); err != nil {
		return err
	}
	p.videoTimer.pause()
	p.State.setPlaying(false)

	return nil
}

// Stop stops the player and clear up the playlist
func (p *Player) Stop() error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.videoTimer == nil {
		return nil
	}

	if err := p.b.stop(); err != nil {
		return err
	}
	p.videoTimer.stop()

	p.playlist = p.playlist[:0]
	p.videoTimer = nil
	p.State.reset()

	return nil
}

// Next sets the player to the next video in the playlist
func (p *Player) Next() error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.videoTimer == nil {
		return nil
	}

	if len(p.playlist) == 0 {
		if err := p.b.stop(); err != nil {
			return err
		}

		p.videoTimer.stop()
		p.videoTimer = nil
		p.State.reset()
	} else {
		info := p.playlist[0]
		p.playlist = p.playlist[1:]

		if err := p.b.set(info.URL); err != nil {
			return err
		}

		p.videoTimer.stop()
		p.videoTimer = newTimer(info.Duration, p.onFinish)
		p.State.next(info.Title, info.Thumbnail)
	}

	return nil
}

// Add adds the url to the playlist
func (p *Player) Add(url string) error {
	p.m.Lock()
	defer p.m.Unlock()

	info, err := extract(url)
	if err != nil {
		return err
	}

	log.Info("[player]", info)

	if p.videoTimer != nil {
		p.playlist = append(p.playlist, info)
		return nil
	}

	if err := p.b.set(info.URL); err != nil {
		return err
	}

	p.videoTimer = newTimer(info.Duration, p.onFinish)
	p.State.next(info.Title, info.Thumbnail)
	return nil
}

// Seek sets the current video to position (from 0 - 1.0)
func (p *Player) Seek(position float64) error {
	p.m.Lock()
	defer p.m.Unlock()

	if err := p.b.seek(position); err != nil {
		return err
	}

	p.videoTimer.seek(position)
	p.State.seek(position)
	return nil
}

// Close cleans up resources
func (p *Player) Close() {
	p.b.close()
}
