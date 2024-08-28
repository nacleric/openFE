package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/ldtkgo"
)

var (
	game        *Game
	ldtkProject *ldtkgo.Project
	floorSprite *ebiten.Image
	unitSprite  *ebiten.Image
)

const (
	GRIDSIZE int = 2
)

const (
	tileSize     float32 = 16
	screenWidth  int     = 256 * 2
	screenHeight int     = 128 * 2
)

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

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Q to quit", 0, 0)
	ebitenutil.DebugPrintAt(screen, "Arrow Keys to move Camera", 0, 16)
	ebitenutil.DebugPrintAt(screen, "Z/X ZoomIn/ZoomOut", 0, 32)
	ebitenutil.DebugPrintAt(screen, "C/V Undo/Redo", 0, 48)

	var cameraOffsetX float32
	var cameraOffsetY float32

	cameraOffsetX = g.camera.x0 * 16 * -1
	cameraOffsetY = g.camera.y0 * 16 * -1

	RenderGrid(screen, &g.mg, cameraOffsetX, cameraOffsetY)
	g.mg.RenderCursor(screen, cameraOffsetX, cameraOffsetY)
	g.mg.RenderUnits(screen, cameraOffsetX, cameraOffsetY)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		fmt.Println(g.mg.grid)
		g.AppendHistory(g.mg)
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

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Platformer")
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
