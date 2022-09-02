package gitdump

import (
	"context"
	"sync"
)

type HostingOptions struct {
	URL      string
	Username string
	Password string
}

type HostingRepo struct {
	SubDir   string
	URL      string
	Username string
	Password string
}

type Hosting interface {
	List(ctx context.Context, opts HostingOptions) (repos []HostingRepo, err error)
}

var (
	hostings                 = map[string]Hosting{}
	hostingsLock sync.Locker = &sync.Mutex{}
)

func SetHosting(vendor string, hosting Hosting) {
	hostingsLock.Lock()
	defer hostingsLock.Unlock()
	hostings[vendor] = hosting
}

func GetHosting(vendor string) Hosting {
	hostingsLock.Lock()
	defer hostingsLock.Unlock()
	return hostings[vendor]
}
