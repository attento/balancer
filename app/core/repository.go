package core

type ConfigRepository interface {
	Get() Config
	Put(c Config)

	Server(a Address) (*Server, bool)
	NewServer(a Address) bool
	RemoveServer(a Address)

	PutFilter(a Address, f Filter)
	AddUpstream(a Address, u *Upstream)
	SetUpstreams(a Address, us []*Upstream)
}
