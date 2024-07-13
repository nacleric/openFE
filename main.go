package main

import (
	"fmt"
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
	GRIDSIZE int = 5
)

func SetGridCellCoord(mg *MGrid) {
	startingX0 := float32(screenWidth / 2)
	startingY0 := float32(screenHeight / 2)
	incX := float32(0)
	incY := float32(0)
	for row := range mg.grid {
		for col := range mg.grid[row] {
			x0 := startingX0 + incX
			y0 := startingY0 + incY
			mg.grid[row][col].x0 = x0
			mg.grid[row][col].y0 = y0
			if col < GRIDSIZE-1 {
				incX += 16 * cameraScale
			} else {
				incX = 0
				incY += 16 * cameraScale
			}
		}
	}
}

func RenderGrid(screen *ebiten.Image, mg *MGrid, offsetX, offsetY float32) {
	startingX0 := float32(screenWidth / 2)
	startingY0 := float32(screenHeight / 2)
	incX := float32(0)
	incY := float32(0)
	for row := range mg.grid {
		for col := range mg.grid[row] {
			x0 := startingX0 + offsetX + incX
			y0 := startingY0 + offsetY + incY
			mg.grid[row][col].isOccupied = false
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

func (sc *SpriteCell) GetRow(cellY int) int {
	return cellY * sc.frameHeight
}

func (sc *SpriteCell) GetCol(cellX int) int {
	return cellX * sc.frameWidth
}

type Job int

const (
	SMALLFOLK Job = iota
	NOBLE
)

type WeaponType int

const (
	BLUNT = iota
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

type AnimationData struct {
	sc             SpriteCell
	frameCount     int // Total number of columns for specific row
	frameFrequency int // How often frames transition
}

type RenderData struct {
	x0          float32
	y0          float32
	idleAnim    AnimationData
	spritesheet *ebiten.Image
}

type Unit struct {
	pX  int
	pY  int
	job Job
	rd  RenderData
}

func CreateUnit(spritesheet *ebiten.Image, j Job) Unit {
	idleAnimData := AnimationData{SpriteCell{0, 0, 16, 16}, 4, 16}
	rd := RenderData{
		x0:          0,
		y0:          0,
		idleAnim:    idleAnimData,
		spritesheet: spritesheet,
	}

	u := Unit{
		job: j,
		rd:  rd,
	}

	return u
}

func (u *Unit) IdleAnimation(screen *ebiten.Image, offsetX, offsetY float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(cameraScale), float64(cameraScale))
	op.GeoM.Translate(float64(u.rd.x0+offsetX), float64(u.rd.y0+offsetY))

	cellX := u.rd.idleAnim.sc.cellX
	cellY := u.rd.idleAnim.sc.cellY

	i := (game.count / u.rd.idleAnim.frameFrequency) % u.rd.idleAnim.frameCount
	sx, sy := u.rd.idleAnim.sc.GetCol(cellX)+i*u.rd.idleAnim.sc.frameWidth, u.rd.idleAnim.sc.GetRow(cellY)
	screen.DrawImage(u.rd.spritesheet.SubImage(image.Rect(sx, sy, sx+u.rd.idleAnim.sc.frameWidth, sy+u.rd.idleAnim.sc.frameHeight)).(*ebiten.Image), op)
}

type GridCell struct {
	x0         float32
	y0         float32
	unit       *Unit
	isOccupied bool
}

type MGrid struct {
	grid  [GRIDSIZE][GRIDSIZE]GridCell
	pc    PlayerCursor
	units []Unit
}

func (mg *MGrid) RenderCursor(screen *ebiten.Image, offsetX, offsetY float32) {
	pY := mg.pc.pY
	pX := mg.pc.pX

	x0 := mg.grid[pY][pX].x0
	y0 := mg.grid[pY][pX].y0
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	vector.StrokeRect(screen, x0+offsetX, y0+offsetY, 16*cameraScale, 16*cameraScale, 1, red, true)
}

func (mg *MGrid) SetUnitPos(u *Unit, new_pX, new_pY int) {
	mg.grid[new_pY][new_pX].unit = u
	mg.grid[new_pY][new_pX].unit.rd.x0 = mg.grid[new_pY][new_pX].x0
	mg.grid[new_pY][new_pX].unit.rd.y0 = mg.grid[new_pY][new_pX].y0
	mg.grid[new_pY][new_pX].unit.pX = new_pX
	mg.grid[new_pY][new_pX].unit.pY = new_pY
}

type Game struct {
	keys   []ebiten.Key
	camera Camera
	mg     MGrid
	count  int
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.count++
	SetGridCellCoord(&g.mg)

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
		g.mg.pc.MoveCursorUp()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.mg.pc.MoveCursorLeft()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.mg.pc.MoveCursorDown()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.mg.pc.MoveCursorRight()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		fmt.Println("unimplemented")
		pX := g.mg.pc.pX
		pY := g.mg.pc.pY
		g.mg.SetUnitPos(&g.mg.units[0], pX, pY)
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

	RenderGrid(screen, &g.mg, cameraOffsetX, cameraOffsetY)
	g.mg.RenderCursor(screen, cameraOffsetX, cameraOffsetY)
	g.mg.units[0].IdleAnimation(screen, cameraOffsetX, cameraOffsetY)

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

type JobStats map[Job]JStats

func init() {
	jobs := JobStats{
		SMALLFOLK: {
			aSpeed:   1,
			movement: 3,
			mounted:  false,
		},
		NOBLE: {
			aSpeed:   1,
			movement: 4,
			mounted:  false,
		},
	}

	// game = &Game{camera: Camera{0, 0}, pc: PlayerCursor{0, 0, 0, 0}}
	// Need to fix default instantiation for mg
	game = &Game{
		camera: Camera{0, 0},
		mg: MGrid{
			pc: PlayerCursor{0, 0, 0, 0},
		},
	}

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
	u := CreateUnit(unitSprite, NOBLE)
	fmt.Println(jobs[u.job])
	game.mg.units = append(game.mg.units, u)
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Platformer")
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
