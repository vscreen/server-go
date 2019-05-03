package supervisor

// Service is an abstraction of mini services in vscreen
type Service interface {
	Serve()
	Close()
}

type Supervisor interface {
	Service
	Name() string
	Children() []Service
}
