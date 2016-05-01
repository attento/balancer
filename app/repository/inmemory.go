package repository

import (
	"sync"
	"github.com/attento/balancer/app/core"
)


type InMemoryConfigRepository struct {
	sync.RWMutex
	c core.Config
}

func NewInMemoryConfigRepository() *InMemoryConfigRepository {
	return &InMemoryConfigRepository{}
}

func (r *InMemoryConfigRepository) Init() {
	r.RLock()
	defer r.RUnlock()
	if r.c.Servers() == nil {
		r.c = core.Create()
	}
}

func (r *InMemoryConfigRepository) Get() core.Config{
	r.Init()
	r.RLock()
	defer r.RUnlock()
	return r.c
}

func (r *InMemoryConfigRepository) Put(c core.Config){
	r.Init()
	r.Lock()
	r.c = c
	r.Unlock()
}

func (r *InMemoryConfigRepository) PutFilter(a core.Address, f core.Filter) {
	r.Init()
	r.Lock()
	r.c.PutFilter(a, f)
	r.Unlock()
}

func (r *InMemoryConfigRepository) NewServer(a core.Address) bool {
	r.Init()
	r.Init()
	r.Lock()
	defer r.Unlock()
	return r.c.NewServer(a)
}

func (r *InMemoryConfigRepository) RemoveServer(a core.Address) {
	r.Init()
	r.Lock()
	r.c.RemoveServer(a)
	r.Unlock()
}

func (r *InMemoryConfigRepository) AddUpstream(a core.Address, u *core.Upstream) {
	r.Init()
	r.Lock()
	r.c.AddUpstream(a, u)
	r.Unlock()
}

func (r *InMemoryConfigRepository) SetUpstreams(a core.Address, us []*core.Upstream) {
	r.Init()

	r.Lock()
	r.c.SetUpstreams(a, core.ConvertConfigUpstreamToMap(us))
	r.Unlock()
}

func (r *InMemoryConfigRepository) Server(a core.Address) (*core.Server, bool) {
	r.Init()
	cnf := r.Get()
	return cnf.Server(a)
}
