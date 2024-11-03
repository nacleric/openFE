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
	_x0         float32
	_y0         float32
	x0y0        f64.Vec2
	idleAnim    AnimationData
	spritesheet *ebiten.Image
}

type Unit struct {
	id        int
	pXHistory []int
	pYHistory []int
	pX        int
	pY        int
	rpg       RPG
	rd        RenderData
}

func CreateUnit(id int, spritesheet *ebiten.Image, rpg RPG, pX, pY int) Unit {
	idleAnimData := AnimationData{SpriteCell{0, 0, 16, 16}, 4, 16}

	GridCellStartingX0 := MapStartingX0 + float64(16*pX)
	GridCellStartingY0 := MapStartingY0 + float64(16*pY)

	rd := RenderData{
		x0y0:        f64.Vec2{GridCellStartingX0, GridCellStartingY0},
		idleAnim:    idleAnimData,
		spritesheet: spritesheet,
	}

	u := Unit{
		id:        id,
		pXHistory: []int{pX},
		pYHistory: []int{pY},
		pX:        pX,
		pY:        pY,
		rpg:       rpg,
		rd:        rd,
	}

	return u
}

func (u *Unit) pXAppendHistory(pX int) {
	u.pXHistory = append(u.pXHistory, pX)
}

func (u *Unit) pYAppendHistory(pY int) {
	u.pYHistory = append(u.pYHistory, pY)
}

func (u *Unit) IdleAnimation(screen *ebiten.Image, offsetX, offsetY float64, count int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(cameraScale), float64(cameraScale))
	x0 := u.rd.x0y0[0]
	y0 := u.rd.x0y0[1]
	op.GeoM.Translate(x0+offsetX, y0+offsetY)

	cellX := u.rd.idleAnim.sc.cellX
	cellY := u.rd.idleAnim.sc.cellY

	i := (count / u.rd.idleAnim.frameFrequency) % u.rd.idleAnim.frameCount
	sx, sy := u.rd.idleAnim.sc.GetCol(cellX)+i*u.rd.idleAnim.sc.frameWidth, u.rd.idleAnim.sc.GetRow(cellY)
	screen.DrawImage(u.rd.spritesheet.SubImage(image.Rect(sx, sy, sx+u.rd.idleAnim.sc.frameWidth, sy+u.rd.idleAnim.sc.frameHeight)).(*ebiten.Image), op)
}
