package core

var (
	// Map Dimension 16 by 8 tiles
	cameraWidth  float64 = 16
	cameraHeight float64 = 8
	// cameraScale          = float64(ScreenWidth) / TileSize / cameraWidth
	cameraScale          = float64(4)

)

type Camera struct {
	X float64
	Y float64
}
