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

func Add(x, y int) int {
	return x + y
}

var (
	game        *Game
	ldtkProject *ldtkgo.Project
	floorSprite *ebiten.Image
	unitSprite  *ebiten.Image
)

const (
	GRIDSIZE int = 2
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
	x0 float32
	y0 float32
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
	BLUNT WeaponType = iota
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

func CreateUnit(spritesheet *ebiten.Image, j Job, pX, pY int) Unit {
	idleAnimData := AnimationData{SpriteCell{0, 0, 16, 16}, 4, 16}
	rd := RenderData{
		x0:          0,
		y0:          0,
		idleAnim:    idleAnimData,
		spritesheet: spritesheet,
	}

	u := Unit{
		pX:  pX,
		pY:  pY,
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
	x0   float32
	y0   float32
	unit *Unit
}

// Will prob delete
func (gc *GridCell) ClearUnit() {
	gc.unit = nil
}

type TurnState int

const (
	SELECTUNIT TurnState = iota
	UNITACTION
	ACTIONOPTIONS // Unused for now
)

type MGrid struct {
	turnState    TurnState
	grid         [GRIDSIZE][GRIDSIZE]GridCell
	pc           PlayerCursor
	units        []Unit
	selectedUnit *Unit
}

func CreateMGrid() MGrid {
	var grid [GRIDSIZE][GRIDSIZE]GridCell

	for i := 0; i < GRIDSIZE; i++ {
		for j := 0; j < GRIDSIZE; j++ {
			grid[i][j] = GridCell{
				unit: nil,
			}
		}
	}

	units := []Unit{}

	mgrid := MGrid{
		turnState:    SELECTUNIT,
		grid:         grid,
		pc:           PlayerCursor{0, 0, 0, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255}},
		units:        units,
		selectedUnit: nil,
	}

	return mgrid
}

func (mg *MGrid) ClearGridCell(pX, pY int) {
	mg.grid[pY][pX].ClearUnit()
}

func (mg *MGrid) SetState(ts TurnState) {
	mg.turnState = ts
}

func (mg *MGrid) QueryCell(pX, pY int) GridCell {
	return mg.grid[pY][pX]
}

func (mg *MGrid) QueryUnit(pX, pY int) *Unit {
	return mg.grid[pY][pX].unit
}

func (mg *MGrid) SetSelectedUnit(u *Unit) {
	mg.selectedUnit = u
}

func (mg *MGrid) ClearSelectedUnit() {
	mg.selectedUnit = nil
}

func (mg *MGrid) RenderCursor(screen *ebiten.Image, offsetX, offsetY float32) {
	pY := mg.pc.pY
	pX := mg.pc.pX

	x0 := mg.grid[pY][pX].x0
	y0 := mg.grid[pY][pX].y0

	vector.StrokeRect(screen, x0+offsetX, y0+offsetY, 16*cameraScale, 16*cameraScale, 1, mg.pc.cursorColor, true)
}

func (mg *MGrid) RenderUnits(screen *ebiten.Image, offsetX, offsetY float32) {
	for _, unit := range mg.units {
		mg.grid[unit.pY][unit.pX].unit = &unit
		mg.grid[unit.pY][unit.pX].unit.rd.x0 = mg.grid[unit.pY][unit.pX].x0
		mg.grid[unit.pY][unit.pX].unit.rd.y0 = mg.grid[unit.pY][unit.pX].y0
		unit.IdleAnimation(screen, offsetX, offsetY)
	}
}

func (mg *MGrid) SetUnitPos(u *Unit, new_pX, new_pY int) {
	mg.ClearGridCell(u.pX, u.pY)
	mg.grid[new_pY][new_pX].unit = u
	mg.grid[new_pY][new_pX].unit.rd.x0 = mg.grid[new_pY][new_pX].x0
	mg.grid[new_pY][new_pX].unit.rd.y0 = mg.grid[new_pY][new_pX].y0
	mg.grid[new_pY][new_pX].unit.pX = new_pX
	mg.grid[new_pY][new_pX].unit.pY = new_pY
	mg.SetSelectedUnit(u)
}

type Game struct {
	keys          []ebiten.Key
	camera        Camera
	mg            MGrid
	count         int
	history       []MGrid
	actionCounter int
}

func (g *Game) AppendHistory(mg MGrid) {
	g.history = append(g.history, mg)
}

func (g *Game) incrementActionCounter() {
	if g.actionCounter < len(g.history)-1 {
		g.actionCounter += 1
	}
}

func (g *Game) deincrementActionCounter() {
	if g.actionCounter > 0 {
		g.actionCounter -= 1
	}
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.count++
	SetGridCellCoord(&g.mg)
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		panic("Game quit change this later")
	}

	// zoom in
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		cameraScale /= .5
	}

	// zoom out
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

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.deincrementActionCounter()
		g.mg = g.history[g.actionCounter]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		g.incrementActionCounter()
		g.mg = g.history[g.actionCounter]
	}

	return nil
}

// Cursor for grid only
type PlayerCursor struct {
	pX          int
	pY          int
	prevX       int
	prevY       int
	cursorColor color.Color
}

type RGB int

// All colors here (might be bad)
const (
	RED RGB = iota
	GREEN
	BLUE
)

func (pc *PlayerCursor) SetColor(rgb RGB) {
	var newColor color.Color
	// might just have all the color options here
	if rgb == RED {
		newColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	} else if rgb == GREEN {
		newColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	} else if rgb == BLUE {
		newColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	}
	pc.cursorColor = newColor
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

	cameraOffsetX = g.camera.x0 * 16 * -1
	cameraOffsetY = g.camera.y0 * 16 * -1

	RenderGrid(screen, &g.mg, cameraOffsetX, cameraOffsetY)
	g.mg.RenderCursor(screen, cameraOffsetX, cameraOffsetY)
	g.mg.RenderUnits(screen, cameraOffsetX, cameraOffsetY)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		fmt.Println(g.mg.grid)

		cursor_pX := g.mg.pc.pX
		cursor_pY := g.mg.pc.pY
		cell := g.mg.QueryCell(cursor_pX, cursor_pY)
		u := cell.unit

		if g.mg.turnState == SELECTUNIT {
			if u != nil {
				g.mg.SetSelectedUnit(u)
				g.mg.pc.SetColor(BLUE)
				g.mg.SetState(UNITACTION)
			} else {
				fmt.Println("No unit found at the selected position")
			}
		} else if g.mg.turnState == UNITACTION {
			if g.mg.selectedUnit != nil {
				if g.mg.selectedUnit.pX == cursor_pX && g.mg.selectedUnit.pY == cursor_pY {
					fmt.Println("clicked tile is on the same tile as selected unit, wasting action")
					g.mg.ClearSelectedUnit()
					g.mg.pc.SetColor(GREEN)
					g.mg.SetState(SELECTUNIT)
				} else {
					// Ensure the units slice is not empty before accessing it
					if len(g.mg.units) > 0 {
						g.mg.SetUnitPos(&g.mg.units[0], cursor_pX, cursor_pY)
						g.mg.ClearSelectedUnit()
						g.AppendHistory(g.mg)
						g.actionCounter += 1
						g.mg.pc.SetColor(GREEN)
						g.mg.SetState(SELECTUNIT)
					} else {
						fmt.Println("No units available to move")
						g.mg.ClearSelectedUnit()
					}
				}
			}
		}

	}

	for _, keyPress := range g.keys {
		switch keyPress {
		case ebiten.KeyUp:
			g.camera.y0 += -.25
		case ebiten.KeyDown:
			g.camera.y0 += .25
		case ebiten.KeyLeft:
			g.camera.x0 += -.25
		case ebiten.KeyRight:
			g.camera.x0 += .25
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

	fmt.Println(jobs)

	// Need to fix default instantiation for mg
	mgrid := CreateMGrid()
	game = &Game{
		camera:  Camera{0, 0},
		mg:      mgrid,
		history: []MGrid{},
	}
	game.AppendHistory(mgrid)

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
	u := CreateUnit(unitSprite, NOBLE, 0, 0)
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
