package game

var (
	// Map Dimension 16 by 8 tiles
	cameraWidth  float32 = 16
	cameraHeight float32 = 8
	cameraScale          = float32(ScreenWidth) / TileSize / cameraWidth
)

type Camera struct {
	X float32
	Y float32
}
