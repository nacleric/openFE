package main

var (
	// Map Dimension 16 by 8 tiles
	cameraWidth  float32 = 16
	cameraHeight float32 = 8
	cameraScale          = float32(screenWidth) / tileSize / cameraWidth
)

type Camera struct {
	x0 float32
	y0 float32
}
