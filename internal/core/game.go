package core

import (
	"fmt"
	"image"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func DebugMessages(screen *ebiten.Image, mg *MGrid) {
	ebitenutil.DebugPrintAt(screen, "Q to quit", 0, 0)
	ebitenutil.DebugPrintAt(screen, "Arrow Keys to move Camera", 0, 16)
	ebitenutil.DebugPrintAt(screen, "Z/X ZoomIn/ZoomOut", 0, 32)
	ebitenutil.DebugPrintAt(screen, "C/V Undo/Redo", 0, 48)
	pc_str := fmt.Sprintf("cursor: [%d, %d]", mg.pc.posXY[0], mg.pc.posXY[1])
	ebitenutil.DebugPrintAt(screen, pc_str, 0, 64)
	cameraScale := fmt.Sprintf("CameraScale: [%d]", cameraScale)
	ebitenutil.DebugPrintAt(screen, cameraScale, 0, 80)
}

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

func (g *Game) IncrementActionCounter() {
	if g.ActionCounter < len(g.History)-1 {
		g.ActionCounter += 1
	}
}

func (g *Game) DeincrementActionCounter() {
	if g.ActionCounter > 0 {
		g.ActionCounter -= 1
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	var cameraOffsetX float64
	var cameraOffsetY float64

	cameraOffsetX = g.Camera.X * 16 * -1
	cameraOffsetY = g.Camera.Y * 16 * -1

	for _, layer := range LdtkProject.Levels[0].Layers {
		for _, tile := range layer.Tiles {
			x0 := float64(tile.Position[0])
			y0 := float64(tile.Position[1])
			op := &ebiten.DrawImageOptions{}

			// No idea why I needed to divide by camerascale
			// in order to fix zoomin zoomout when I didn't need to do that for unitsprite
			op.GeoM.Translate(float64(x0+cameraOffsetX/cameraScale), float64(y0+cameraOffsetY/cameraScale))
			op.GeoM.Scale(float64(cameraScale), float64(cameraScale))
			screen.DrawImage(FloorSprite.SubImage(image.Rect(tile.Src[0], tile.Src[1], tile.Src[0]+16, tile.Src[1]+16)).(*ebiten.Image), op)
		}
	}

	RenderGrid(screen, &g.MG, cameraOffsetX, cameraOffsetY)
	if g.MG.turnState == UNITMOVEMENT {
		g.MG.RenderLegalPositions(screen, cameraOffsetX, cameraOffsetY, g.Count)
	}
	g.MG.RenderCursor(screen, cameraOffsetX, cameraOffsetY, g.Count)
	g.MG.RenderUnits(screen, cameraOffsetX, cameraOffsetY, g.Count)
	DebugMessages(screen, &g.MG)
}

func (g *Game) Update() error {
	g.Keys = inpututil.AppendPressedKeys(g.Keys[:0])
	g.Count++
	SetGridCellCoord(&g.MG, MapStartingX0, MapStartingY0)
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
		fmt.Println("Undo is pressed")
		g.DeincrementActionCounter()
		g.MG = g.History[g.ActionCounter]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		fmt.Println("Redo is pressed")
		g.IncrementActionCounter()
		g.MG = g.History[g.ActionCounter]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		fmt.Println("debugger triggered")
	}

	enterPressed := inpututil.IsKeyJustPressed(ebiten.KeyEnter)

	if g.MG.turnState == SELECTUNIT && enterPressed {
		cursor_posXY := g.MG.pc.posXY
		cell := g.MG.QueryCell(cursor_posXY)
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if cell.unitId != notSelected {
				g.MG.SetSelectedUnit(cell.unitId)
				g.MG.pc.SetColor(BLUE)
				legalPositions := reachableCells(&g.MG, cursor_posXY, GRIDSIZE, 2)
				g.MG.legalPositions = legalPositions
				g.MG.SetState(UNITMOVEMENT)
			} else {
				fmt.Println("No unit found at the selected position")
			}
		}

		enterPressed = false
	}

	if g.MG.turnState == UNITMOVEMENT && enterPressed {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			cursor_posXY := g.MG.pc.posXY
			// Note: might be removed
			cursor_posX := cursor_posXY[X]
			cursor_posY := cursor_posXY[Y]
			// --
			selectedUnitId := g.MG.selectedUnit
			selectedUnit := g.MG.Units[selectedUnitId]
			if selectedUnit.posXY[X] == cursor_posX && selectedUnit.posXY[Y] == cursor_posY {
				fmt.Println("clicked tile is on the same tile as selected unit, wasting action")
				g.MG.ClearSelectedUnit()
				g.MG.pc.SetColor(GREEN)
				g.MG.SetState(SELECTUNIT)
			} else if slices.Contains(g.MG.legalPositions, cursor_posXY) {
				fmt.Println("legalMove")
				g.MG.SetUnitPos(selectedUnit, cursor_posXY)
				g.MG.pc.SetColor(GREEN)
				g.MG.SetState(SELECTUNIT)
				g.AppendHistory(g.MG)
				g.MG.Units[selectedUnit.id].posXYAppendHistory(cursor_posXY)
				g.MG.ClearSelectedUnit()
				g.ActionCounter += 1
				// g.MG.SetState(UNITACTIONS)
			} else {
				fmt.Println("not legalMove")
			}
		}

		enterPressed = false
	}

	if g.MG.turnState == UNITACTIONS {
		fmt.Println("select actions for player")
	}

	// Will need to be refactored. Put state above keypress
	/*
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			cursor_posXY := g.MG.pc.posXY
			// Note: might be removed
			cursor_posX := cursor_posXY[X]
			cursor_posY := cursor_posXY[Y]
			cell := g.MG.QueryCell(cursor_posXY)
			if g.MG.turnState == SELECTUNIT {
				if cell.unitId != notSelected {
					g.MG.SetSelectedUnit(cell.unitId)
					g.MG.pc.SetColor(BLUE)
					g.MG.SetState(UNITMOVEMENT)
					legalPositions := reachableCells(&g.MG, cursor_posXY, GRIDSIZE, 2)
					g.MG.legalPositions = legalPositions

				} else {
					fmt.Println("No unit found at the selected position")
				}
			} else if g.MG.turnState == UNITMOVEMENT {
				selectedUnitId := g.MG.selectedUnit
				selectedUnit := g.MG.Units[selectedUnitId]
				if selectedUnit.posXY[X] == cursor_posX && selectedUnit.posXY[Y] == cursor_posY {
					fmt.Println("clicked tile is on the same tile as selected unit, wasting action")
					g.MG.ClearSelectedUnit()
					g.MG.pc.SetColor(GREEN)
					g.MG.SetState(SELECTUNIT)
				} else if slices.Contains(g.MG.legalPositions, cursor_posXY) {
					fmt.Println("legalMove")
					g.MG.SetUnitPos(selectedUnit, cursor_posXY)
					g.MG.pc.SetColor(GREEN)
					g.MG.SetState(SELECTUNIT)
					g.AppendHistory(g.MG)
					g.MG.Units[selectedUnit.id].posXYAppendHistory(cursor_posXY)
					g.MG.ClearSelectedUnit()
					g.ActionCounter += 1
					// g.MG.SetState(UNITACTIONS)
				} else {
					fmt.Println("not legalMove")
				}
			} else if g.MG.turnState == UNITACTIONS {
				fmt.Println("select actions for player")
			}
		}
	*/

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

	return nil
}
