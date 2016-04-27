package core
import "math/rand"


type ElectionAlgorithm func(u map[string]*Upstream) (*Upstream, error)

type Elector interface {
	Elect(f ElectionAlgorithm) (*Upstream, error)
}

var RoundRobin = func (u map[string]*Upstream) (*Upstream, error) {
	i := int(float32(len(u)) * rand.Float32())
	for _, v := range u {
		if i == 0 {
			return v, nil
		} else {
			i--
		}
	}
	panic("impossible to be here")
}

var LeastConn = func (u map[string]*Upstream) (*Upstream, error) {
	// @todo here
	return nil, nil
}

var leastConn = func (u map[string]*Upstream) (*Upstream, error) {
	// @todo here
	return nil, nil
}

func (s *Server) Elect(f ElectionAlgorithm) (*Upstream, error) {
	return f(s.upstreams)
}


