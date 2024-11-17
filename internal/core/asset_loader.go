package core

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/ldtkgo"
)

var (
	LdtkProject *ldtkgo.Project
	FloorSprite *ebiten.Image
	UnitSprite  *ebiten.Image
	CursorSprite *ebiten.Image
)

func LoadSpritesheets() {
	var err error
	UnitSprite, _, err = ebitenutil.NewImageFromFile("../../assets/demo/eliwood_map_idle.png")
	if err != nil {
		log.Fatal(err)
	}

	// 20x20 offset by 2x2 during render 40 wide for both sprites
	CursorSprite, _, err = ebitenutil.NewImageFromFile("../../assets/demo/cursor.png")
	if err != nil {
		log.Fatal(err)
	}

	LdtkProject, err = ldtkgo.Open("../../assets/demo/8x8.ldtk")
	if err != nil {
		panic("Map file doesn't exist")
	}

	FloorSprite, _, err = ebitenutil.NewImageFromFile("../../assets/demo/experiment.png")
	if err != nil {
		panic("Tilemap doesn't exist")
	}
}
