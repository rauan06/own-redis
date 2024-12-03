package models

import "context"

type AsyncMap struct {
	Map map[string]*ChanContextPair
}

type ChanContextPair struct {
	Data       chan string
	Context    context.Context
	CancelFunc context.CancelFunc
}
