package nocpoy

import "sync"

type noCopy struct{}

var _ sync.Locker = (*noCopy)(nil)

func (c *noCopy) Lock()   {}
func (c *noCopy) Unlock() {}

type Url struct {
	noCopy noCopy
}
