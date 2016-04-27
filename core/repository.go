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

func  (r *InMemoryConfigRepository) Init() {
	r.Lock()
	if r.c == nil {
		r.c = make(map[Address]*Server)
	}
	r.Unlock()
}

func (r *InMemoryConfigRepository) Get() Config{
	r.Init()
	r.RLock()
	defer r.RUnlock()
	return r.c
}

func (r *InMemoryConfigRepository) Put(c Config){
	r.Init()
	r.Lock()
	r.c = c
	r.Unlock()
}

func (r *InMemoryConfigRepository) PutFilter(a Address, f Filter) {
	r.Init()
	r.Lock()
	r.c.PutFilter(a, f)
	r.Unlock()
}

func (r *InMemoryConfigRepository) NewServer(a Address) {
	r.Init()
	r.Lock()
	r.c.NewServer(a)
	r.Unlock()
}

func (r *InMemoryConfigRepository) RemoveServer(a Address) {
	r.Init()
	r.Lock()
	r.c.RemoveServer(a)
	r.Unlock()
}

func (r *InMemoryConfigRepository) AddUpstream(a Address, u *Upstream) {
	r.Init()
	r.Lock()
	r.c.AddUpstream(a, u)
	r.Unlock()
}

func (r *InMemoryConfigRepository) SetUpstreams(a Address, us map[string]*Upstream) {
	r.Init()
	r.Lock()
	r.c.SetUpstreams(a, us)
	r.Unlock()
}

func (r *InMemoryConfigRepository) Refresh(s *Server) (*Server, bool) {
	cnf := r.Get()
	return cnf.Server(s.address)
}
