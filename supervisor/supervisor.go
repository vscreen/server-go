package supervisor

import (
	"github.com/thejerf/suture"
)

// From creates a new supervisor tree with name as the service name
func From(name string, services ...Service) *suture.Supervisor {
	s := suture.NewSimple(name)
	from(s, services...)
	return s
}

func from(s *suture.Supervisor, services ...Service) {
	for _, service := range services {
		super, ok := service.(Supervisor)
		if !ok {
			s.Add(convert(service))
			continue
		}

		newSuper := suture.NewSimple(super.Name())
		from(newSuper, super.Children()...)
		s.Add(newSuper)
		s.Add(convert(super))
	}
}
