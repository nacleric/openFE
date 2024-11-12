package core

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/math/f64"
)

/*
func reachableCells(mg *MGrid, pos PosXY, maxMoveDistance int) []PosXY {
	var directions = []PosXY{
		{0, -1}, // Up
		{0, 1},  // Down
		{-1, 0}, // Left
		{1, 0},  // Right
	}

	legalPositions := []PosXY{}
	rows := GRIDSIZE
	cols := GRIDSIZE

	queue := []PosXY{pos}
	visited := make([][]bool, rows)
	// fills visited with false 2d array, will have to change to account for objects
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	// starting position will be true
	visited[pos[1]][pos[0]] = true

	for len(queue) > 0 {
		// Deque the first cell
		current := queue[0]
		queue := queue[1:]
		col, row := current[0], current[1]

		for _, direction := range directions {
			newCol, newRow := col+direction[0], row+direction[1]

			// Check if the new cell is within bounds and not yet visited
			if newRow >= 0 && newRow < rows && newCol >= 0 && newCol < cols && !visited[newRow][newCol] {
				// Mark the new cell as visited and enqueue it
				visited[newRow][newCol] = true
				queue = append(queue, PosXY{newRow, newCol})
				legalPositions = append(legalPositions, PosXY{newCol, newRow})
			}
		}
	}

	return legalPositions
}
*/

func isValid(visited [][]bool, col, row int, col_len, row_len int) bool {
	if (row < 0 || col < 0 || row >= row_len || col >= col_len) {
		return false
	}

	if visited[row][col] {
		return false
	}

	return true
}


// Need to include maxMove distance
func reachableCells(mg *MGrid, pos PosXY, gridSize int) {
	fmt.Println("line 73")
	var directions = []PosXY{
		{0, -1}, // Up
		{0, 1},  // Down
		{-1, 0}, // Left
		{1, 0},  // Right
	}

	// legalPositions := []PosXY{}
	row_len := gridSize
	col_len := gridSize

	queue := []PosXY{pos}

	fmt.Println("line 86")
	visited := make([][]bool, row_len)
	for i := 0; i < row_len; i++ {
		visited[i] = make([]bool, col_len)
		for j := 0; j < col_len; j++ {
			visited[i][j] = false
		}
	}

	visited[pos[0]][pos[1]] = true

	fmt.Println("line 97")
	for len(queue) > 0 {
		current := queue[0]
		// dequeue first cell
		queue = queue[1:]

		col, row := current[0], current[1]
		for _, direction := range directions {
			adjacentCol, adjacentRow := col + direction[0], row + direction[1]
			// Check if cell is out of bound
			fmt.Println(direction,":",isValid(visited, col, row, col_len, row_len))
			if isValid(visited, col, row, col_len, row_len) {
				queue = append(queue, PosXY{col, row})
				visited[adjacentRow][adjacentCol] = true
				// legalPositions := append(legalPositions, Pos{col, row})
			}
		} 
	}
	fmt.Println(visited)
	fmt.Println("ran fine")
}

const emptyCell = -1

type GridCell struct {
	x0y0   f64.Vec2
	unitId int
}

// Will prob delete
func (gc *GridCell) ClearUnit() {
	gc.unitId = emptyCell
}

const notSelected = -1

type MGrid struct {
	turnState    TurnState
	grid         [][]GridCell
	pc           PlayerCursor
	Units        []*Unit
	selectedUnit int // UnitID, it is -1 if there is no selected unit
}

func (mg *MGrid) SearchUnit() {
}

func CreateMGrid(units []*Unit, gridSize int) MGrid {
	grid := make([][]GridCell, gridSize)
	for i := 0; i < gridSize; i++ {
		grid[i] = make([]GridCell, gridSize) // Initialize each row
		for j := 0; j < gridSize; j++ {
			grid[i][j] = GridCell{
				unitId: emptyCell,
			}
		}
	}

	for _, u := range units {
		pX := u.posXY[0]
		pY := u.posXY[1]
		grid[pY][pX].unitId = u.id
	}

	mgrid := MGrid{
		turnState:    SELECTUNIT,
		grid:         grid,
		pc:           PlayerCursor{PosXY{0, 0}, PosXY{0, 0}, color.RGBA{R: 0, G: 255, B: 0, A: 255}},
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

func (mg *MGrid) QueryCell(posXY PosXY) GridCell {
	pX := posXY[0]
	pY := posXY[1]
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

func (mg *MGrid) RenderCursor(screen *ebiten.Image, offsetX, offsetY float64) {
	f32cameraScale := float32(cameraScale)
	f32offsetX := float32(offsetX)
	f32offsetY := float32(offsetY)
	pX := mg.pc.posXY[0]
	pY := mg.pc.posXY[1]
	x0y0 := mg.grid[pY][pX].x0y0

	vector.StrokeRect(screen, float32(x0y0[0])+f32offsetX, float32(x0y0[1])+f32offsetY, 16*f32cameraScale, 16*f32cameraScale, 1, mg.pc.cursorColor, true)
}

func (mg *MGrid) RenderUnits(screen *ebiten.Image, offsetX, offsetY float64, count int) {
	for _, unit := range mg.Units {
		pX := unit.posXY[0]
		pY := unit.posXY[1]
		unit.rd.x0y0 = mg.grid[pY][pX].x0y0
		unit.IdleAnimation(screen, offsetX, offsetY, count)
	}
}

func (mg *MGrid) SetUnitPos(u *Unit, new_posXY PosXY) {
	new_pX := new_posXY[0]
	new_pY := new_posXY[1]
	mg.ClearGridCell(u.posXY[0], u.posXY[1])
	newGridCellPos := &mg.grid[new_pY][new_pX]

	// New grid location
	newGridCellPos.unitId = u.id
	u.rd.x0y0 = newGridCellPos.x0y0
	newPos := PosXY{new_pX, new_pY}
	u.posXY = newPos
}

func SetGridCellCoord(mg *MGrid, startingX0, startingY0 float64) {
	incX := float64(0)
	incY := float64(0)
	for row := range mg.grid {
		for col := range mg.grid[row] {
			x0 := startingX0 + incX
			y0 := startingY0 + incY
			mg.grid[row][col].x0y0 = f64.Vec2{x0, y0}
			if col < GRIDSIZE-1 {
				incX += 16 * cameraScale // No Idea why I needed to multiply this
			} else {
				incX = 0
				incY += 16 * cameraScale
			}
		}
	}
}

func RenderGrid(screen *ebiten.Image, mg *MGrid, offsetX, offsetY float64) {
	incX := float64(0)
	incY := float64(0)
	f32cameraScale := float32(cameraScale)
	for row := range mg.grid {
		for col := range mg.grid[row] {
			x0 := MapStartingX0 + offsetX + incX
			y0 := MapStartingY0 + offsetY + incY
			vector.StrokeRect(screen, float32(x0), float32(y0), 16*f32cameraScale, 16*f32cameraScale, 1, color.White, true)
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
	posXY       PosXY
	prevXY      PosXY
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
	pY := &pc.posXY[1]
	if *pY > 0 {
		pc.SetPrevCursor(pc.posXY)
		*pY -= 1
	}
}

func (pc *PlayerCursor) MoveCursorLeft() {
	pX := &pc.posXY[0]
	if *pX > 0 {
		pc.SetPrevCursor(pc.posXY)
		*pX -= 1
	}
}

func (pc *PlayerCursor) MoveCursorDown() {
	pY := &pc.posXY[1]
	if *pY < GRIDSIZE-1 {
		pc.SetPrevCursor(pc.posXY)
		*pY += 1
	}
}

func (pc *PlayerCursor) MoveCursorRight() {
	pX := &pc.posXY[0]
	if *pX < GRIDSIZE-1 {
		pc.SetPrevCursor(pc.posXY)
		*pX += 1
	}
}

// Might need this data for animations
func (pc *PlayerCursor) SetPrevCursor(posXY PosXY) {
	pc.prevXY = posXY
}
