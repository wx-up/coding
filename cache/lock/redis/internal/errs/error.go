package errs

import "errors"

var (
	ErrFailedToPreemptLock = errors.New("redis-lock：抢锁失败")
	ErrLockNotHand         = errors.New("redis-lock：未持有锁")
)
