package main

import (
	"image"
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
	GRIDSIZE = 3
)

func setGridCellCoord(grid *[GRIDSIZE][GRIDSIZE]GridCell) {
	startingX0 := float32(screenWidth / 2)
	startingY0 := float32(screenHeight / 2)
	incX := float32(0)
	incY := float32(0)
	for row := range grid {
		for col := range grid[row] {
			x0 := startingX0 + incX
			y0 := startingY0 + incY
			grid[row][col].x0 = x0
			grid[row][col].y0 = y0
			if col < GRIDSIZE-1 {
				incX += 16 * cameraScale
			} else {
				incX = 0
				incY += 16 * cameraScale
			}
		}
	}
}

func renderGrid(screen *ebiten.Image, grid *[GRIDSIZE][GRIDSIZE]GridCell, offsetX, offsetY float32) {
	startingX0 := float32(screenWidth / 2)
	startingY0 := float32(screenHeight / 2)
	incX := float32(0)
	incY := float32(0)
	for row := range grid {
		for col := range grid[row] {
			x0 := startingX0 + offsetX + incX
			y0 := startingY0 + offsetY + incY
			grid[row][col].isOccupied = false
			vector.StrokeRect(screen, x0, y0, 16*cameraScale, 16*cameraScale, 1, color.White, true)
			if col < GRIDSIZE-1 {
				incX += 16 * cameraScale
			} else {
				incX = 0
				incY += 16 * cameraScale
			}
		}
	}
}

func (pc *PlayerCursor) renderCursor(screen *ebiten.Image, grid *[GRIDSIZE][GRIDSIZE]GridCell, offsetX, offsetY float32) {
	x0 := grid[pc.pY][pc.pX].x0
	y0 := grid[pc.pY][pc.pX].y0
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	vector.StrokeRect(screen, x0+offsetX, y0+offsetY, 16*cameraScale, 16*cameraScale, 1, red, true)
}

const (
	tileSize     float32 = 16
	screenWidth  int     = 256 * 2
	screenHeight int     = 128 * 2
)

var (
	// Map Dimension 16 by 8 tiles
	cameraWidth  float32 = 16
	cameraHeight float32 = 8
	cameraScale          = float32(screenWidth) / tileSize / cameraWidth
)

type Camera struct {
	pX float32
	pY float32
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
	x0          float32
	y0          float32
	idleAnim    AnimationData
	spritesheet *ebiten.Image
}

func CreateUnit(spritesheet *ebiten.Image) Unit {
	idleAnimData := AnimationData{SpriteCell{0, 0, 16, 16}, 4, 16}

	u := Unit{
		x0:          0,
		y0:          0,
		idleAnim:    idleAnimData,
		spritesheet: spritesheet,
	}

	return u
}

func (u *Unit) IdleAnimation(screen *ebiten.Image, offsetX, offsetY float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(cameraScale), float64(cameraScale))
	op.GeoM.Translate(float64(u.x0+offsetX), float64(u.y0+offsetY))

	cellX := u.idleAnim.sc.cellX
	cellY := u.idleAnim.sc.cellY

	i := (game.count / u.idleAnim.frameFrequency) % u.idleAnim.frameCount
	sx, sy := u.idleAnim.sc.getCol(cellX)+i*u.idleAnim.sc.frameWidth, u.idleAnim.sc.getRow(cellY)
	screen.DrawImage(u.spritesheet.SubImage(image.Rect(sx, sy, sx+u.idleAnim.sc.frameWidth, sy+u.idleAnim.sc.frameHeight)).(*ebiten.Image), op)
}

type GridCell struct {
	x0         float32
	y0         float32
	unit       *Unit
	isOccupied bool
}

type Game struct {
	keys   []ebiten.Key
	camera Camera
	grid   [GRIDSIZE][GRIDSIZE]GridCell
	count  int
	units  []Unit
	pc     PlayerCursor
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.count++
	setGridCellCoord(&g.grid)

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		panic("Game quit change this later")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		cameraScale /= .5
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		cameraScale *= .5
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.pc.MoveCursorUp()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.pc.MoveCursorLeft()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.pc.MoveCursorDown()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.pc.MoveCursorRight()
	}

	return nil
}

// Cursor for grid only
type PlayerCursor struct {
	pX    int
	pY    int
	prevX int
	prevY int
}

func (pc *PlayerCursor) MoveCursorUp() {
	if pc.pY > 0 {
		pc.SetPrevCursor(pc.pX, pc.pY)
		pc.pY -= 1
	}
}

func (pc *PlayerCursor) MoveCursorLeft() {
	if pc.pX > 0 {
		pc.SetPrevCursor(pc.pX, pc.pY)
		pc.pX -= 1
	}
}

func (pc *PlayerCursor) MoveCursorDown() {
	if pc.pY < GRIDSIZE-1 {
		pc.SetPrevCursor(pc.pX, pc.pY)
		pc.pY += 1
	}
}

func (pc *PlayerCursor) MoveCursorRight() {
	if pc.pX < GRIDSIZE-1 {
		pc.SetPrevCursor(pc.pX, pc.pY)
		pc.pX += 1
	}
}

// Might need this data for animations
func (pc *PlayerCursor) SetPrevCursor(pX, pY int) {
	pc.prevX = pX
	pc.prevX = pY
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Q to quit", 0, 0)
	ebitenutil.DebugPrintAt(screen, "Arrow Keys to move Camera", 0, 16)
	ebitenutil.DebugPrintAt(screen, "Z/X ZoomIn/ZoomOut", 0, 32)

	var cameraOffsetX float32
	var cameraOffsetY float32

	cameraOffsetX = g.camera.pX * 16 * -1
	cameraOffsetY = g.camera.pY * 16 * -1

	renderGrid(screen, &g.grid, cameraOffsetX, cameraOffsetY)
	g.pc.renderCursor(screen, &g.grid, cameraOffsetX, cameraOffsetY)

	g.grid[0][0].unit = &g.units[0]
	g.grid[0][0].unit.x0 = g.grid[0][0].x0
	g.grid[0][0].unit.y0 = g.grid[0][0].y0
	g.grid[0][0].unit.IdleAnimation(screen, cameraOffsetX, cameraOffsetY)

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
	game = &Game{camera: Camera{0, 0}, pc: PlayerCursor{0, 0, 0, 0}}
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
	u := CreateUnit(unitSprite)
	game.units = append(game.units, u)
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Platformer")
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
