package core

var (
	// Map Dimension 16 by 8 tiles
	cameraWidth  float64 = 16
	cameraHeight float64 = 8
	cameraScale          = float64(ScreenWidth) / TileSize / cameraWidth
)

type Camera struct {
	X float32
	Y float32
}
