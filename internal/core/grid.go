package core

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const emptyCell = -1

type GridCell struct {
	x0     float32
	y0     float32
	unitId int
}

// Will prob delete
func (gc *GridCell) ClearUnit() {
	gc.unitId = emptyCell
}

const notSelected = -1

type MGrid struct {
	turnState    TurnState
	grid         [GRIDSIZE][GRIDSIZE]GridCell
	pc           PlayerCursor
	Units        []*Unit
	selectedUnit int // UnitID, it is -1 if there is no selected unit
}

func (mg *MGrid) SearchUnit() {

}

func CreateMGrid(units []*Unit) MGrid {
	var grid [GRIDSIZE][GRIDSIZE]GridCell

	for i := 0; i < GRIDSIZE; i++ {
		for j := 0; j < GRIDSIZE; j++ {
			grid[i][j] = GridCell{
				unitId: emptyCell,
			}
		}
	}

	for _, u := range units {
		grid[u.pY][u.pX].unitId = u.id
	}

	mgrid := MGrid{
		turnState:    SELECTUNIT,
		grid:         grid,
		pc:           PlayerCursor{0, 0, 0, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255}},
		Units:        units,
		selectedUnit: notSelected,
	}

	SetGridCellCoord(&mgrid, MapStartingX0, MapStartingY0)

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

func (mg *MGrid) QueryUnit(pX, pY int) int {
	return mg.grid[pY][pX].unitId
}

func (mg *MGrid) SetSelectedUnit(id int) {
	mg.selectedUnit = id
}

func (mg *MGrid) ClearSelectedUnit() {
	mg.selectedUnit = notSelected
}

func (mg *MGrid) RenderCursor(screen *ebiten.Image, offsetX, offsetY float32) {
	pY := mg.pc.pY
	pX := mg.pc.pX

	x0 := mg.grid[pY][pX].x0
	y0 := mg.grid[pY][pX].y0

	vector.StrokeRect(screen, x0+offsetX, y0+offsetY, 16*cameraScale, 16*cameraScale, 1, mg.pc.cursorColor, true)
}

func (mg *MGrid) RenderUnits(screen *ebiten.Image, offsetX, offsetY float32, count int) {
	for _, unit := range mg.Units {
		unit.rd.x0 = mg.grid[unit.pY][unit.pX].x0
		unit.rd.y0 = mg.grid[unit.pY][unit.pX].y0
		unit.IdleAnimation(screen, offsetX, offsetY, count)
	}
}

func (mg *MGrid) SetUnitPos(u *Unit, new_pX, new_pY int) {
	mg.ClearGridCell(u.pX, u.pY)
	newGridCellPos := &mg.grid[new_pY][new_pX]

	// New grid location
	newGridCellPos.unitId = u.id
	fmt.Println(newGridCellPos.unitId)
	u.rd.x0 = newGridCellPos.x0
	u.rd.y0 = newGridCellPos.y0
	u.pX = new_pX
	u.pY = new_pY
}

func SetGridCellCoord(mg *MGrid, startingX0, startingY0 float32) {
	incX := float32(0)
	incY := float32(0)
	for row := range mg.grid {
		for col := range mg.grid[row] {
			x0 := startingX0 + incX
			y0 := startingY0 + incY
			mg.grid[row][col].x0 = x0
			mg.grid[row][col].y0 = y0
			if col < GRIDSIZE-1 {
				incX += 16 * cameraScale // No Idea why I needed to multiply this
			} else {
				incX = 0
				incY += 16 * cameraScale
			}
		}
	}
}

func RenderGrid(screen *ebiten.Image, mg *MGrid, offsetX, offsetY float32) {
	incX := float32(0)
	incY := float32(0)
	for row := range mg.grid {
		for col := range mg.grid[row] {
			x0 := MapStartingX0 + offsetX + incX
			y0 := MapStartingY0 + offsetY + incY
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

type TurnState int

const (
	SELECTUNIT TurnState = iota
	UNITACTION
	ACTIONOPTIONS // Unused for now
)

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
