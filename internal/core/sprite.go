package core

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

type SpriteCell struct {
	cellX       int // column of spritesheet Ex: 0 is first col 16 is 2nd col
	cellY       int // row of spritesheet Ex: 16 is 2nd row
	frameWidth  int // Size of Sprite frame (most likely 16x16)
	frameHeight int
}

func (sc *SpriteCell) GetRow(cellY int) int {
	return cellY * sc.frameHeight
}

func (sc *SpriteCell) GetCol(cellX int) int {
	return cellX * sc.frameWidth
}

type AnimationData struct {
	sc             SpriteCell
	frameCount     int // Total number of columns for specific row
	frameFrequency int // How often frames transition
}

type RenderData struct {
	x0y0        f64.Vec2
	ad          AnimationData
	spritesheet *ebiten.Image
}

type PosXY [2]int

type Unit struct {
	id           int
	posXYHistory []PosXY
	posXY        PosXY
	rpg          RPG
	rd           RenderData
}

func CreateUnit(id int, spritesheet *ebiten.Image, rpg RPG, posXY PosXY) Unit {
	idleAnimData := AnimationData{SpriteCell{0, 0, 16, 16}, 4, 16}

	GridCellStartingX0 := MapStartingX0 + float64(16*posXY[X])
	GridCellStartingY0 := MapStartingY0 + float64(16*posXY[Y])

	rd := RenderData{
		x0y0:        f64.Vec2{GridCellStartingX0, GridCellStartingY0},
		ad:          idleAnimData,
		spritesheet: spritesheet,
	}

	u := Unit{
		id:           id,
		posXYHistory: []PosXY{posXY},
		posXY:        posXY,
		rpg:          rpg,
		rd:           rd,
	}

	return u
}

func (u *Unit) posXYAppendHistory(posXY PosXY) {
	u.posXYHistory = append(u.posXYHistory, posXY)
}

func (u *Unit) IdleAnimation(screen *ebiten.Image, offsetX, offsetY float64, count int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(CAMERASCALE), float64(CAMERASCALE))
	// Note: might move render calculation to where it's being called
	x0 := u.rd.x0y0[X]
	y0 := u.rd.x0y0[Y]
	op.GeoM.Translate(x0+offsetX, y0+offsetY)

	cellX := u.rd.ad.sc.cellX
	cellY := u.rd.ad.sc.cellY

	i := (count / u.rd.ad.frameFrequency) % u.rd.ad.frameCount
	sx, sy := u.rd.ad.sc.GetCol(cellX)+i*u.rd.ad.sc.frameWidth, u.rd.ad.sc.GetRow(cellY)
	screen.DrawImage(u.rd.spritesheet.SubImage(image.Rect(sx, sy, sx+u.rd.ad.sc.frameWidth, sy+u.rd.ad.sc.frameHeight)).(*ebiten.Image), op)
}
