package errs

import (
	"errors"
	"fmt"
)

var (
	ErrCacheIsCompletely = errors.New("cache：缓存已经满了")
	ErrSetKeyFail        = errors.New("cache：设置键值对失败")
)

func NewErrKeyNotFound(key string) error {
	return fmt.Errorf("cache：找不到 key %s", key)
}
