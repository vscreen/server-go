package supervisor

import (
	"github.com/thejerf/suture"
)

type sutureAdaptor struct {
	Service
}

func (s *sutureAdaptor) Stop() {
	s.Stop()
}

// convert transforms Service to suture.Service
func convert(s Service) suture.Service {
	return &sutureAdaptor{s}
}
