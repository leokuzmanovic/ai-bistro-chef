package constants

import "time"

const (
	TimeoutQuick    = 2 * time.Second
	TimeoutRegular  = 10 * time.Second
	TimeoutGenerous = 20 * time.Second
)
