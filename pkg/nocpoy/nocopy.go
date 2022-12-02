package nocpoy

type noCopy struct{}

func (c *noCopy) Lock()   {}
func (c *noCopy) Unlock() {}

type Url struct {
	noCopy noCopy
}
