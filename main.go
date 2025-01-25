package main

import (
	_ "image/png"
	"log" // Adjust based on where these are defined

	"github.com/hajimehoshi/ebiten/v2"

	// Import your internal package
	core "openFE/internal/core" // Use alias to avoid conflict
)

var (
	game *core.Game
)

func init() {
	core.LoadSpritesheets()

	unitInfo := core.RPG{Job: core.NOBLE, Movement: 2}
	u := core.CreateUnit(0, core.UnitSprite, unitInfo, core.PosXY{0, 1})
	i := core.CreateUnit(1, core.UnitSprite, unitInfo, core.PosXY{1, 0})

	units := []core.Unit{u, i}
	unitPointers := make([]*core.Unit, len(units))
	for i := range units {
		unitPointers[i] = &units[i]
	}
	mgrid := core.CreateMGrid(unitPointers, core.CursorSprite, core.LdtkProject)

	actionMenu := core.ActionMenu{MenuOptions: []string{"foobar"}, Selected: 0}
	menuManager := core.MenuManager{ActionMenu: actionMenu}

	game = &core.Game{
		Camera:      core.Camera{X: 0, Y: 0},
		MG:          mgrid,
		History:     []core.MGrid{},
		MenuManager: menuManager,
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
