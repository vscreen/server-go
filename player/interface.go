package player

// Info represents the current state of the video player
type Info struct {
	Title     string
	Thumbnail string
	Volume    float64
	Position  float64
	Playing   bool
}

// Player is an abstraction of a video player. It contains
// some basic operations that should be available in most
// video player applications out there.
type Player interface {
	Play() error
	Pause() error
	Stop() error
	Next() error
	Add(url string) error
	Seek(position float64) error

	// Start starts the video player. Some initializations should be
	// done here. This should be the main loop for the video. Meaning,
	// Start blocks the main program flow.
	Start()
	InfoListener() <-chan Info
	Close()
}
