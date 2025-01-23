package lock

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

var LockExpire = 30

var LockName = "distributed_lock"
