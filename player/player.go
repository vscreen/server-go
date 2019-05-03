package player

import (
	"errors"
	"fmt"

	"github.com/vscreen/server-go/player/backend"
	"github.com/vscreen/server-go/player/extractor"
	svr "github.com/vscreen/server-go/supervisor"
)

type msg func(b backend.Player, i *State)
type action func(b backend.Player, i *State) error

// Player is an abstraction of a video player.
// the backend player is determined by the platform
type Player struct {
	mailbox        chan msg
	subscribed     bool
	subscriberChan chan State
	backend        backend.Player
	children       []svr.Service
}

// New creates Player abstraction
func New(player string) (*Player, error) {
	extractor.Init()
	b, err := backend.New(player)
	if err != nil {
		return nil, err
	}

	buffSize := 64
	mailbox := make(chan msg, buffSize)
	subscriberChan := make(chan State)
	p := Player{
		mailbox:        mailbox,
		subscriberChan: subscriberChan,
		backend:        b,
		children:       nil,
	}
	return &p, nil
}

// Subscribe lets the subscriber knows the newest info.
// If there's already a subscriber, Subscribe will return
// a nil channel and error will be set.
func (p *Player) Subscribe() (<-chan State, error) {
	if p.subscribed {
		return nil, errors.New("can't have more than 1 subscribers")
	}

	p.subscribed = true
	return p.subscriberChan, nil
}

// Name describes supervisor's name
func (p *Player) Name() string {
	return "player"
}

// Children lists all childrens to be supervised
func (p *Player) Children() []svr.Service {
	return p.children
}

// Serve start the server and satisfied VscreenService
func (p *Player) Serve() {
	var s State
	s.timer = newTimer(p.onFinish)
	var mailbox <-chan msg = p.mailbox
	var subscriberChan chan<- State = p.subscriberChan
	b := p.backend

	for {
		action := <-mailbox
		action(b, &s)

		// If there's a slow receiver, curInfo will just be destroyed and
		// continue looping.
		select {
		case subscriberChan <- s:
		default:
		}
	}
}

func (p *Player) send(a action) error {
	e := make(chan error)
	p.mailbox <- func(b backend.Player, s *State) {
		defer close(e)
		e <- a(b, s)
	}
	return <-e
}

func (p *Player) onFinish() {
	p.send(func(b backend.Player, s *State) error {
		playlist := s.Playlist
		fmt.Println("onFinish")
		if len(playlist) == 0 {
			s.reset()
			return nil
		}

		next := playlist[0]
		if err := b.Set(next.URL); err != nil {
			return err
		}
		s.next()
		return nil
	})
}

// Play plays the current video. If there's no video,
// nothing will happen
func (p *Player) Play() error {
	return p.send(func(b backend.Player, s *State) error {
		if err := b.Play(); err != nil {
			return err
		}

		s.play()
		return nil
	})
}

// Pause pauses the current video. If there's no video,
// nothing will happen
func (p *Player) Pause() error {
	return p.send(func(b backend.Player, s *State) error {
		if err := b.Pause(); err != nil {
			return err
		}

		s.pause()
		return nil
	})
}

// Stop stops the player and clear up the playlist
func (p *Player) Stop() error {
	return p.send(func(b backend.Player, s *State) error {
		if err := b.Stop(); err != nil {
			return err
		}

		s.reset()
		return nil
	})
}

// Next sets the player to the next video in the playlist
func (p *Player) Next() error {
	return p.send(func(b backend.Player, s *State) error {
		if len(s.Playlist) == 0 {
			if err := b.Stop(); err != nil {
				return err
			}
		} else {
			next := s.Playlist[0].URL
			if err := b.Set(next); err != nil {
				return err
			}
		}

		s.next()
		return nil
	})
}

// Add adds the url to the playlist
func (p *Player) Add(url string) error {
	return p.send(func(b backend.Player, s *State) error {
		info, err := extractor.Extract(url)
		if err != nil {
			return err
		}

		s.Playlist = append(s.Playlist, info)

		if s.Playing || s.Title != "" {
			return nil
		}

		if err := b.Set(info.URL); err != nil {
			return err
		}

		s.next()
		return nil
	})
}

// Seek sets the current video to position (from 0 - 1.0)
func (p *Player) Seek(position float64) error {
	return p.send(func(b backend.Player, s *State) error {
		if err := b.Seek(position); err != nil {
			return err
		}

		s.seek(position)
		return nil
	})
}

// Close cleans up resources
func (p *Player) Close() {
	p.send(func(b backend.Player, s *State) error {
		b.Close()
		return nil
	})
}
