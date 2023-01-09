package orm

import "github.com/wx-up/coding/orm/internal/errs"

// ErrNoRows 用户不能直接使用 internal 包中的 ErrNoRows 错误，所以才会有在外部定义一个同名的错误
var ErrNoRows = errs.ErrNoRows
