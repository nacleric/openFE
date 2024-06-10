package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/ldtkgo"
)

var (
	game        *Game
	ldtkProject *ldtkgo.Project
	floorSprite *ebiten.Image
	grid        [8][8]uint8
)

func renderGrid(screen *ebiten.Image, grid [8][8]uint8, offsetX, offsetY float64) {
	startingPosX := float64(screenWidth / 2)
	startingPosY := float64(screenHeight / 2)
	incPosition := float32(0)
	for row, _ := range grid {
		for col, _ := range grid[row] {
			fmt.Println(col)
			vector.StrokeRect(screen, float32(startingPosX+offsetX)+incPosition, float32(startingPosY+offsetY), 16, 16, 1, color.White, true)
			incPosition += 16
		}
	}
}

const (
	tileSize     = 16
	screenWidth  = 256
	screenHeight = 128

	// Map Dimension 16 by 8 tiles
	cameraWidth          = 16
	cameraHeight         = 8
	cameraScale  float64 = float64(screenWidth) / tileSize / float64(cameraWidth)
)

type Camera struct {
	pX float64
	pY float64
}

type Game struct {
	keys   []ebiten.Key
	camera Camera
	grid   [8][8]uint8
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		panic("Game quit change this later")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Q to quit", 0, 0)
	ebitenutil.DebugPrintAt(screen, "Arrow Keys to move Camera", 0, 16)

	var cameraOffsetX float64
	var cameraOffsetY float64

	cameraOffsetX = float64(g.camera.pX) * 16 * -1
	cameraOffsetY = float64(g.camera.pY) * 16 * -1

	renderGrid(screen, g.grid, cameraOffsetX, cameraOffsetY)

	/*
		for _, layer := range ldtkProject.Levels[0].Layers {
			for _, tile := range layer.Tiles {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.Position[0])+cameraOffsetX, float64(tile.Position[1])+cameraOffsetY)
				op.GeoM.Scale(cameraScale, cameraScale)
				screen.DrawImage(floorSprite.SubImage(image.Rect(tile.Src[0], tile.Src[1], tile.Src[0]+16, tile.Src[1]+16)).(*ebiten.Image), op)
			}
		}
	*/

	for _, keyPress := range g.keys {
		switch keyPress {
		case ebiten.KeyUp:
			g.camera.pY += -.25
		case ebiten.KeyDown:
			g.camera.pY += .25
		case ebiten.KeyLeft:
			g.camera.pX += -.25
		case ebiten.KeyRight:
			g.camera.pX += .25
		default:
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func init() {
	game = &Game{camera: Camera{0, 0}}
	var err error
	ldtkProject, err = ldtkgo.Open("assets/demo/demo.ldtk")
	if err != nil {
		panic("Map file doesn't exist")
	}

	floorSprite, _, err = ebitenutil.NewImageFromFile("./assets/demo/floor.png")
	if err != nil {
		panic("Tilemap doesn't exist")
	}

}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Platformer")
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
