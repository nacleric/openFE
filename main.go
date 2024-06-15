package main

import (
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
	unitSprite  *ebiten.Image
)

const (
	GRIDSIZE = 8
)

func renderGrid(screen *ebiten.Image, grid [GRIDSIZE][GRIDSIZE]uint8, offsetX, offsetY float32) {
	startingPosX := float32(screenWidth / 2)
	startingPosY := float32(screenHeight / 2)
	incX := float32(0)
	incY := float32(0)
	for row := range grid {
		for col := range grid[row] {
			vector.StrokeRect(screen, startingPosX+offsetX+incX, startingPosY+offsetY+incY, 16*cameraScale, 16*cameraScale, 1, color.White, true)
			if col < GRIDSIZE-1 {
				incX += 16 * cameraScale
			} else {
				incX = 0
				incY += 16 * cameraScale
			}
		}
	}
}

const (
	tileSize     float32 = 16
	screenWidth  int     = 256 * 2
	screenHeight int     = 128 * 2

	// Map Dimension 16 by 8 tiles
	cameraWidth  = 16
	cameraHeight = 8
	cameraScale  = float32(screenWidth) / tileSize / cameraWidth
)

type Camera struct {
	pX float32
	pY float32
}

type Game struct {
	keys   []ebiten.Key
	camera Camera
	grid   [GRIDSIZE][GRIDSIZE]uint8
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		panic("Game quit change this later")
	}
	return nil
}

type SpriteCell struct {
	cellX       int // column of spritesheet Ex: 0 is first col 16 is 2nd col
	cellY       int // row of spritesheet Ex: 16 is 2nd row
	frameWidth  int // Size of Sprite frame (most likely 16x16)
	frameHeight int
}

func (sc *SpriteCell) getRow(cellY int) int {
	return cellY * sc.frameHeight
}

func (sc *SpriteCell) getCol(cellX int) int {
	return cellX * sc.frameWidth
}

type AnimationData struct {
	sc             SpriteCell
	frameCount     int // Total number of columns for specific row
	frameFrequency int // How often frames transition
}

type Unit struct {
	posX     float32
	posY     float32
	width    int
	height   int
	idleAnim AnimationData
}

func CreateUnit(spritesheet *ebiten.Image) Unit {
	idleAnimData := AnimationData{SpriteCell{0, 0, 16, 16}, 2, 64}

	u := Unit{
		posX:     0,
		posY:     0,
		width:    16,
		height:   16,
		idleAnim: idleAnimData,
	}

	return u
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Q to quit", 0, 0)
	ebitenutil.DebugPrintAt(screen, "Arrow Keys to move Camera", 0, 16)

	var cameraOffsetX float32
	var cameraOffsetY float32

	cameraOffsetX = g.camera.pX * 16 * -1
	cameraOffsetY = g.camera.pY * 16 * -1

	renderGrid(screen, g.grid, cameraOffsetX, cameraOffsetY)

	/*
		cellX := p.idleAnim.sc.cellX
		cellY := p.idleAnim.sc.cellY

		i := (game.count / p.idleAnim.frameFrequency) % p.idleAnim.frameCount
		sx, sy := p.idleAnim.sc.getCol(cellX)+i*p.idleAnim.sc.frameWidth, p.idleAnim.sc.getRow(cellY)
		screen.DrawImage(p.spritesheet.SubImage(image.Rect(sx, sy, sx+p.idleAnim.sc.frameWidth, sy+p.idleAnim.sc.frameHeight)).(*ebiten.Image), op)
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

func LoadSpritesheets() {
	var err error
	unitSprite, _, err = ebitenutil.NewImageFromFile("./assets/demo/eliwood_map_idle.png")
	if err != nil {
		log.Fatal(err)
	}
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

	LoadSpritesheets()
	CreateUnit(unitSprite)

}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Platformer")
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
