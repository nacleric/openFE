package main

import (
	_ "image/png"
	"log" // Adjust based on where these are defined

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/ldtkgo"

	// Import your internal package
	core "openFE/internal/core" // Use alias to avoid conflict
)

var (
	game        *core.Game
	ldtkProject *ldtkgo.Project
	floorSprite *ebiten.Image
	unitSprite  *ebiten.Image
)

func LoadSpritesheets() {
	var err error
	unitSprite, _, err = ebitenutil.NewImageFromFile("../../assets/demo/eliwood_map_idle.png")
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	LoadSpritesheets()
	var err error
	ldtkProject, err = ldtkgo.Open("../../assets/demo/demo.ldtk")
	if err != nil {
		panic("Map file doesn't exist")
	}

	floorSprite, _, err = ebitenutil.NewImageFromFile("../../assets/demo/floor.png")
	if err != nil {
		panic("Tilemap doesn't exist")
	}

	u := core.CreateUnit(0, unitSprite, core.NOBLE, 0, 1)
	i := core.CreateUnit(1, unitSprite, core.NOBLE, 1, 0)

	units := []core.Unit{u, i}
	unitPointers := make([]*core.Unit, len(units))
	for i := range units {
		unitPointers[i] = &units[i]
	}
	mgrid := core.CreateMGrid(unitPointers)

	game = &core.Game{
		Camera:  core.Camera{X: 0, Y: 0},
		MG:      mgrid,
		History: []core.MGrid{},
	}

	game.AppendHistory(game.MG)
	game.IncrementActionCounter()
}

func main() {
	ebiten.SetWindowSize(core.ScreenWidth*2, core.ScreenHeight*2)
	ebiten.SetWindowTitle("Platformer")
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
