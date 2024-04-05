package executor

import "errors"

var (
	ErrNotScheduled = errors.New("this runnable isn't in scheduled state")
)
