package core

var (
	// Map Dimension 16 by 8 tiles
	cameraWidth  float64 = 16
	cameraHeight float64 = 8
	// CAMERASCALE          = float64(ScreenWidth) / TileSize / cameraWidth
	CAMERASCALE = float64(2)
)

type Camera struct {
	X float64
	Y float64
}
