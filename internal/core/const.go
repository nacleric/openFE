package core

const (
	GRIDSIZE int = 5
)

const (
	TileSize float64 = 16
	// ScreenWidth  int     = 256 * 2
	// ScreenHeight int     = 128 * 2
	ScreenWidth  int = 800
	ScreenHeight int = 600
)

// Map Stuff
// const (
// 	MapStartingX0 float32 = float32(ScreenWidth / 2)
// 	MapStartingY0         = float32(ScreenHeight / 2)
// )

const (
	MapStartingX0 float64 = 0
	MapStartingY0         = 0
)

type TurnState int

const (
	SELECTUNIT TurnState = iota
	UNITMOVEMENT
	UNITACTIONS // Unused for now
)

const (
	X = 0
	Y = 1
)

