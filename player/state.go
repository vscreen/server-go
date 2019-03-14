package player

import "errors"

// Info represents the current state of the video player
type Info struct {
	Title     string
	Thumbnail string
	Volume    float64
	Position  float64
	Playing   bool
}

type action func(curInfo *Info)

func stateReset(curInfo *Info) {
	*curInfo = Info{}
}

// state maintains internal video's state by using the actor model.
type state struct {
	c              chan action
	curInfo        Info
	subscriberChan chan Info
	subscribed     bool
}

func newState() *state {
	buffSize := 64

	s := state{
		c:              make(chan action, buffSize),
		subscriberChan: make(chan Info),
	}
	go s.actorLoop(s.c)
	return &s
}

func (s *state) actorLoop(mailbox <-chan action) {
	for {
		action := <-mailbox
		action(&s.curInfo)

		// If there's a slow receiver, curInfo will just be destroyed and
		// continue looping.
		select {
		case s.subscriberChan <- s.curInfo:
		default:
		}
	}
}

func (s *state) dispatch(act action) {
	s.c <- act
}

// Subscribe lets the subscriber knows the newest info.
// If there's already a subscriber, Subscribe will return
// a nil channel and error will be set.
func (s *state) Subscribe() (<-chan Info, error) {
	if s.subscribed {
		return nil, errors.New("can't have more than 1 subscribers")
	}

	return s.subscriberChan, nil
}
