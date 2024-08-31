package main

import (
	_ "image/png"
	"log" // Adjust based on where these are defined

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/ldtkgo"

	// Import your internal package
	core "openFE/internal/game" // Use alias to avoid conflict
)

var (
	game        *core.Game
	ldtkProject *ldtkgo.Project
	floorSprite *ebiten.Image
	unitSprite  *ebiten.Image
)

func LoadSpritesheets() {
	var err error
	unitSprite, _, err = ebitenutil.NewImageFromFile("./assets/demo/eliwood_map_idle.png")
	if err != nil {
		log.Fatal(err)
	}
}

func init() {

	// Need to fix default instantiation for mg
	mgrid := core.CreateMGrid()
	game = &core.Game{
		Camera:  core.Camera{X: 0, Y: 0},
		MG:      mgrid,
		History: []core.MGrid{},
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
	u := core.CreateUnit(unitSprite, core.NOBLE, 0, 0)
	game.MG.Units = append(game.MG.Units, u)
}

func main() {
	ebiten.SetWindowSize(core.ScreenWidth*2, core.ScreenHeight*2)
	ebiten.SetWindowTitle("Platformer")
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
