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
	playlist   []*VideoInfo
	videoTimer *timer
	m          *sync.Mutex
}

func new(b backendPlayer) (*Player, error) {
	return &Player{
		State:      newState(),
		b:          b,
		playlist:   make([]*VideoInfo, 0),
		videoTimer: nil,
		m:          &sync.Mutex{},
	}, nil
}

func (p *Player) onFinish() {
	p.m.Lock()
	defer p.m.Unlock()

	if len(p.playlist) == 0 {
		p.State.dispatch(stateReset)

		p.videoTimer = nil
	} else {
		info := p.playlist[0]
		p.playlist = p.playlist[1:]

		// TODO! maybe todo something with error here
		p.b.set(info.URL)
		p.State.dispatch(func(curInfo *Info) {
			curInfo.Title = info.Title
			curInfo.Thumbnail = info.Thumbnail
			curInfo.Position = 0.0
			curInfo.Playing = true
		})
		p.videoTimer = newTimer(info.Duration, p.onFinish)
	}
}

// Play plays the current video. If there's no video,
// nothing will happen
func (p *Player) Play() error {
	p.m.Lock()
	defer p.m.Unlock()

	playing := false
	if p.videoTimer != nil {
		if err := p.b.play(); err != nil {
			return err
		}
		p.videoTimer.play()
		playing = true
	}

	p.State.dispatch(func(curInfo *Info) {
		curInfo.Playing = playing
	})
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

	p.State.dispatch(func(curInfo *Info) {
		curInfo.Position = p.videoTimer.pos()
		curInfo.Playing = false
	})

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
	p.State.dispatch(stateReset)
	p.videoTimer = nil

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

		p.State.dispatch(stateReset)
	} else {
		info := p.playlist[0]
		p.playlist = p.playlist[1:]

		if err := p.b.set(info.URL); err != nil {
			return err
		}
		p.videoTimer.stop()
		p.videoTimer = newTimer(info.Duration, p.onFinish)

		p.State.dispatch(func(curInfo *Info) {
			curInfo.Position = 0.0
			curInfo.Playing = true
			curInfo.Title = info.Title
			curInfo.Thumbnail = info.Thumbnail
		})
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

	p.State.dispatch(func(curInfo *Info) {
		curInfo.Position = 0.0
		curInfo.Playing = true
		curInfo.Title = info.Title
		curInfo.Thumbnail = info.Thumbnail
	})

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

	p.State.dispatch(func(curInfo *Info) {
		curInfo.Position = position
	})

	return nil
}

// Close cleans up resources
func (p *Player) Close() {
	p.b.close()
}
