package models

type AsyncMap struct {
	Map map[string]ChanContextPair
}

type ChanContextPair struct {
	Data   chan string
	Cancel chan struct{}
}
