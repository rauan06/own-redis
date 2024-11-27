package models

type AsyncMap struct {
	Map map[string](chan string)
}
