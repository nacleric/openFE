package core

func Add(x, y int) int {
	return x + y
}

type Job int

const (
	HOPLITE Job = iota
	GAMBLER
	NOBLE
)

type RPG struct {
	Job      Job
	Movement int
}
