package constants

import "time"

const (
	TimeoutQuick    = 2 * time.Second
	TimeoutRegular  = 5 * time.Second
	TimeoutGenerous = 10 * time.Second
)
