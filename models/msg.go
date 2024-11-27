package models

import "time"

type Messege struct {
	Cmd   string
	Key   string
	Value string
	PX    time.Duration // Default is 0
}
