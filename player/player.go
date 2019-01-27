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
	add(url string) error
	seek(position float64) error
	close()
}

// Info represents the current state of the video player
type Info struct {
	Title     string
	Thumbnail string
	Volume    float64
	Position  float64
	Playing   bool
}

// Player is an abstraction of a video player.
// the backend player is determined by the platform
type Player struct {
	b          backendPlayer
	playlist   []*VideoInfo
	infoCur    Info
	infoChan   chan Info
	videoTimer *timer
	m          *sync.Mutex
}

func new(b backendPlayer) (*Player, error) {
	return &Player{
		b:          b,
		playlist:   make([]*VideoInfo, 0),
		infoChan:   make(chan Info),
		infoCur:    Info{},
		videoTimer: nil,
		m:          &sync.Mutex{},
	}, nil
}

func (p *Player) onFinish() {
	p.m.Lock()
	defer p.m.Unlock()

	if len(p.playlist) == 0 {
		p.infoCur.Title = ""
		p.infoCur.Thumbnail = ""
		p.infoCur.Position = 0.0
		p.infoCur.Playing = false
		p.videoTimer = nil
	} else {
		info := p.playlist[0]
		p.playlist = p.playlist[1:]

		// TODO! maybe todo something with error here
		p.b.add(info.URL)
		p.infoCur.Title = info.Title
		p.infoCur.Thumbnail = info.Thumbnail
		p.infoCur.Position = 0.0
		p.infoCur.Playing = true
		p.videoTimer = newTimer(info.Duration, p.onFinish)
	}

	p.infoChan <- p.infoCur
}

// Play plays the current video. If there's no video,
// nothing will happen
func (p *Player) Play() error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.videoTimer == nil {
		p.infoCur.Playing = false
	} else {
		if err := p.b.play(); err != nil {
			return err
		}
		p.videoTimer.play()
		p.infoCur.Playing = true
	}

	p.infoChan <- p.infoCur
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

	p.infoCur.Playing = false
	p.infoCur.Position = p.videoTimer.pos()
	p.infoChan <- p.infoCur

	return nil
}

// Stop stops the player and clear up the playlist
func (p *Player) Stop() error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.videoTimer == nil {
		return nil
	}

	// stop == seek(1.0)
	if err := p.b.seek(1.0); err != nil {
		return err
	}
	p.videoTimer.stop()

	p.playlist = p.playlist[:0]
	p.infoCur.Position = 0.0
	p.infoCur.Playing = false
	p.infoCur.Title = ""
	p.infoCur.Thumbnail = ""
	p.videoTimer = nil
	p.infoChan <- p.infoCur

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
		if err := p.b.seek(1.0); err != nil {
			return err
		}

		p.videoTimer.stop()
		p.videoTimer = nil

		p.infoCur.Position = 0.0
		p.infoCur.Playing = false
		p.infoCur.Title = ""
		p.infoCur.Thumbnail = ""
	} else {
		info := p.playlist[0]
		p.playlist = p.playlist[1:]

		if err := p.b.add(info.URL); err != nil {
			return err
		}
		p.videoTimer.stop()
		p.videoTimer = newTimer(info.Duration, p.onFinish)

		p.infoCur.Title = info.Title
		p.infoCur.Thumbnail = info.Thumbnail
		p.infoCur.Position = 0.0
		p.infoCur.Playing = true
	}

	p.infoChan <- p.infoCur

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

	if err := p.b.add(url); err != nil {
		return err
	}
	p.videoTimer = newTimer(info.Duration, p.onFinish)

	p.infoCur.Title = info.Title
	p.infoCur.Thumbnail = info.Thumbnail
	p.infoCur.Position = 0.0
	p.infoCur.Playing = true
	p.infoChan <- p.infoCur

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

	p.infoCur.Position = position
	p.infoChan <- p.infoCur

	return nil
}

// InfoListener returns a channel to get the most up to date info
func (p *Player) InfoListener() <-chan Info {
	return p.infoChan
}

// Close cleans up resources
func (p *Player) Close() {
	p.b.close()
	close(p.infoChan)
}
