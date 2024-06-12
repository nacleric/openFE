package main

import (
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

var (
	game        *Game
	ldtkProject *ldtkgo.Project
	floorSprite *ebiten.Image
	unitSprite  *ebiten.Image
)

const (
	GRIDSIZE = 8
)

func renderGrid(screen *ebiten.Image, grid [GRIDSIZE][GRIDSIZE]uint8, offsetX, offsetY float32) {
	startingPosX := float32(screenWidth / 2)
	startingPosY := float32(screenHeight / 2)
	incX := float32(0)
	incY := float32(0)
	for row := range grid {
		for col := range grid[row] {
			vector.StrokeRect(screen, startingPosX+offsetX+incX, startingPosY+offsetY+incY, 16*cameraScale, 16*cameraScale, 1, color.White, true)
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

	// Map Dimension 16 by 8 tiles
	cameraWidth  = 16
	cameraHeight = 8
	cameraScale  = float32(screenWidth) / tileSize / cameraWidth
)

type Camera struct {
	pX float32
	pY float32
}

type Game struct {
	keys   []ebiten.Key
	camera Camera
	grid   [GRIDSIZE][GRIDSIZE]uint8
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		panic("Game quit change this later")
	}
	return nil
}

/*
func (p *Player) IdleAnimation(screen *ebiten.Image, offSetX, offSetY float64) {
	op := &ebiten.DrawImageOptions{}
	if p.direction == "LEFT" {
		op.GeoM.Scale(cameraScale*-1, cameraScale)
		op.GeoM.Translate(float64(p.idleAnim.sc.frameWidth)*cameraScale, 0)
	} else if p.direction == "RIGHT" {
		op.GeoM.Scale(cameraScale, cameraScale)
	}

	op.GeoM.Translate(p.posX+offSetX, p.posY+offSetY)

	cellX := p.idleAnim.sc.cellX
	cellY := p.idleAnim.sc.cellY

	i := (game.count / p.idleAnim.frameFrequency) % p.idleAnim.frameCount
	sx, sy := p.idleAnim.sc.getCol(cellX)+i*p.idleAnim.sc.frameWidth, p.idleAnim.sc.getRow(cellY)
	screen.DrawImage(p.spritesheet.SubImage(image.Rect(sx, sy, sx+p.idleAnim.sc.frameWidth, sy+p.idleAnim.sc.frameHeight)).(*ebiten.Image), op)
}
*/
/*
type Unit struct {
	gridCoord
}
*/
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Q to quit", 0, 0)
	ebitenutil.DebugPrintAt(screen, "Arrow Keys to move Camera", 0, 16)

	var cameraOffsetX float32
	var cameraOffsetY float32

	cameraOffsetX = g.camera.pX * 16 * -1
	cameraOffsetY = g.camera.pY * 16 * -1

	renderGrid(screen, g.grid, cameraOffsetX, cameraOffsetY)
	screen.DrawImage(floorSprite.SubImage(image.Rect(tile.Src[0], tile.Src[1], tile.Src[0]+16, tile.Src[1]+16)).(*ebiten.Image), op)
	/*
		for _, layer := range ldtkProject.Levels[0].Layers {
			for _, tile := range layer.Tiles {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.Position[0])+cameraOffsetX, float64(tile.Position[1])+cameraOffsetY)
				op.GeoM.Scale(cameraScale, cameraScale)
				screen.DrawImage(floorSprite.SubImage(image.Rect(tile.Src[0], tile.Src[1], tile.Src[0]+16, tile.Src[1]+16)).(*ebiten.Image), op)
			}
		}
	*/

	for _, keyPress := range g.keys {
		switch keyPress {
		case ebiten.KeyUp:
			g.camera.pY += -.25
		case ebiten.KeyDown:
			g.camera.pY += .25
		case ebiten.KeyLeft:
			g.camera.pX += -.25
		case ebiten.KeyRight:
			g.camera.pX += .25
		default:
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func LoadSpritesheets() {
	var err error
	unitSprite, _, err = ebitenutil.NewImageFromFile("./assets/Characters/chicken_sprites.png")
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	game = &Game{camera: Camera{0, 0}}
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

}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Platformer")
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
