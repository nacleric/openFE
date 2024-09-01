package core

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Keys          []ebiten.Key
	Camera        Camera
	MG            MGrid
	Count         int
	History       []MGrid
	ActionCounter int
}

func (g *Game) AppendHistory(mg MGrid) {
	g.History = append(g.History, mg)
}

func (g *Game) incrementActionCounter() {
	if g.ActionCounter < len(g.History)-1 {
		g.ActionCounter += 1
	}
}

func (g *Game) deincrementActionCounter() {
	if g.ActionCounter > 0 {
		g.ActionCounter -= 1
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Q to quit", 0, 0)
	ebitenutil.DebugPrintAt(screen, "Arrow Keys to move Camera", 0, 16)
	ebitenutil.DebugPrintAt(screen, "Z/X ZoomIn/ZoomOut", 0, 32)
	ebitenutil.DebugPrintAt(screen, "C/V Undo/Redo", 0, 48)

	var cameraOffsetX float32
	var cameraOffsetY float32

	cameraOffsetX = g.Camera.X * 16 * -1
	cameraOffsetY = g.Camera.Y * 16 * -1

	RenderGrid(screen, &g.MG, cameraOffsetX, cameraOffsetY)
	g.MG.RenderCursor(screen, cameraOffsetX, cameraOffsetY)
	g.MG.RenderUnits(screen, cameraOffsetX, cameraOffsetY, g.Count)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		fmt.Println(g.MG.grid)
		g.AppendHistory(g.MG)
		cursor_pX := g.MG.pc.pX
		cursor_pY := g.MG.pc.pY
		cell := g.MG.QueryCell(cursor_pX, cursor_pY)
		u := cell.unit

		if g.MG.turnState == SELECTUNIT {
			if u != nil {
				g.MG.SetSelectedUnit(u)
				g.MG.pc.SetColor(BLUE)
				g.MG.SetState(UNITACTION)
			} else {
				fmt.Println("No unit found at the selected position")
			}
		} else if g.MG.turnState == UNITACTION {
			if g.MG.selectedUnit != nil {
				if g.MG.selectedUnit.pX == cursor_pX && g.MG.selectedUnit.pY == cursor_pY {
					fmt.Println("clicked tile is on the same tile as selected unit, wasting action")
					g.MG.ClearSelectedUnit()
					g.MG.pc.SetColor(GREEN)
					g.MG.SetState(SELECTUNIT)
				} else {
					// Ensure the units slice is not empty before accessing it
					if len(g.MG.Units) > 0 {
						g.MG.SetUnitPos(&g.MG.Units[0], cursor_pX, cursor_pY)
						g.MG.ClearSelectedUnit()
						g.ActionCounter += 1
						g.MG.pc.SetColor(GREEN)
						g.MG.SetState(SELECTUNIT)
					} else {
						fmt.Println("No units available to move")
						g.MG.ClearSelectedUnit()
					}
				}
			}
		}

	}

	for _, keyPress := range g.Keys {
		switch keyPress {
		case ebiten.KeyUp:
			g.Camera.Y += -.25
		case ebiten.KeyDown:
			g.Camera.Y += .25
		case ebiten.KeyLeft:
			g.Camera.X += -.25
		case ebiten.KeyRight:
			g.Camera.X += .25
		default:
		}
	}

}

func (g *Game) Update() error {
	g.Keys = inpututil.AppendPressedKeys(g.Keys[:0])
	g.Count++
	SetGridCellCoord(&g.MG)
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
		g.MG.pc.MoveCursorUp()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.MG.pc.MoveCursorLeft()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.MG.pc.MoveCursorDown()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.MG.pc.MoveCursorRight()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.deincrementActionCounter()
		g.MG = g.History[g.ActionCounter]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		g.incrementActionCounter()
		g.MG = g.History[g.ActionCounter]
	}

	return nil
}
