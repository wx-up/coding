package pattern

import (
	"context"

	"golang.org/x/sync/singleflight"
)

type ReadThroughSingleFlight struct {
	ReadThroughCache
}

func NewReadThroughSingleFlight(loadFunc func(ctx context.Context, key string) (any, error)) *ReadThroughSingleFlight {
	g := singleflight.Group{}
	res := &ReadThroughSingleFlight{
		ReadThroughCache: ReadThroughCache{
			Loader: LoadFunc(func(ctx context.Context, key string) (any, error) {
				defer func() {
					g.Forget(key)
				}()
				val, err, _ := g.Do(key, func() (interface{}, error) {
					return loadFunc(ctx, key)
				})
				return val, err
			}),
		},
	}
	return res
}
