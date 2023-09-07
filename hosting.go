package gitdump

import (
	"context"
	"errors"
	"sync"
)

type HostingOptions struct {
	URL      string
	Username string
	Password string
}

func (ho HostingOptions) MustURL() {
	if ho.URL == "" {
		panic(errors.New("missing argument 'URL'"))
	}
}

func (ho HostingOptions) MustUsername() {
	if ho.Username == "" {
		panic(errors.New("missing argument 'Username'"))
	}
}

func (ho HostingOptions) MustPassword() {
	if ho.Password == "" {
		panic(errors.New("missing argument 'Password'"))
	}
}

type HostingRepo struct {
	SubDir   string
	URL      string
	Username string
	Password string
}

type Hosting interface {
	List(ctx context.Context, opts HostingOptions) (out []HostingRepo, err error)
}

var (
	hostingRegistry                 = map[string]Hosting{}
	hostingRegistryLock sync.Locker = &sync.Mutex{}
)

func SetHosting(vendor string, hosting Hosting) {
	hostingRegistryLock.Lock()
	defer hostingRegistryLock.Unlock()
	hostingRegistry[vendor] = hosting
}

func GetHosting(vendor string) Hosting {
	hostingRegistryLock.Lock()
	defer hostingRegistryLock.Unlock()
	return hostingRegistry[vendor]
}
