package backend

// Player is an abstraction of a video player. It contains
// some basic operations that should be available in most
// video player applications out there.
type Player interface {
	// These operations are thread safe
	Play() error
	Pause() error
	Stop() error
	Set(url string) error
	Seek(position float64) error
	Close()
}
