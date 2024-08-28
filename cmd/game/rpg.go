package main

func Add(x, y int) int {
	return x + y
}

type Job int

const (
	SMALLFOLK Job = iota
	NOBLE
)

type WeaponType int

const (
	BLUNT WeaponType = iota
	PIERCE
	SLICE
	POSITIONAL
)

type BStats struct {
	bSpeed int
	str    int
}

type JStats struct {
	aSpeed   int
	movement int
	mounted  bool
}
