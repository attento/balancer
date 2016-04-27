package core
import "sync"

type ConfigRepository interface {
	Get() Config
	Put(c Config)
}

type InMemoryConfigRepository struct {
	sync.RWMutex
	c Config
}

var InMemoryRepository = &InMemoryConfigRepository{}

func (r *InMemoryConfigRepository) Get() Config{
	r.RLock()
	defer r.RUnlock()
	return r.c
}

func (r *InMemoryConfigRepository) Put(c Config){
	r.Lock()
	r.c = c
	r.Unlock()
}

func (r *InMemoryConfigRepository) PutFilter(a Address, f Filter) {
	r.Lock()
	r.c.PutFilter(a, f)
	r.Unlock()
}

func (r *InMemoryConfigRepository) AddUpstream(a Address, u *Upstream) {
	r.Lock()
	r.c.AddUpstream(a, u)
	r.Unlock()
}

func (r *InMemoryConfigRepository) Refresh(s *Server) (*Server, bool) {
	cnf := r.Get()
	return cnf.Server(s.address)
}
