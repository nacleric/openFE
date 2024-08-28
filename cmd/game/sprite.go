package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
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
	x0          float32
	y0          float32
	idleAnim    AnimationData
	spritesheet *ebiten.Image
}

type Unit struct {
	pX  int
	pY  int
	job Job
	rd  RenderData
}

func CreateUnit(spritesheet *ebiten.Image, j Job, pX, pY int) Unit {
	idleAnimData := AnimationData{SpriteCell{0, 0, 16, 16}, 4, 16}
	rd := RenderData{
		x0:          0,
		y0:          0,
		idleAnim:    idleAnimData,
		spritesheet: spritesheet,
	}

	u := Unit{
		pX:  pX,
		pY:  pY,
		job: j,
		rd:  rd,
	}

	return u
}

func (u *Unit) IdleAnimation(screen *ebiten.Image, offsetX, offsetY float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(cameraScale), float64(cameraScale))
	op.GeoM.Translate(float64(u.rd.x0+offsetX), float64(u.rd.y0+offsetY))

	cellX := u.rd.idleAnim.sc.cellX
	cellY := u.rd.idleAnim.sc.cellY

	i := (game.count / u.rd.idleAnim.frameFrequency) % u.rd.idleAnim.frameCount
	sx, sy := u.rd.idleAnim.sc.GetCol(cellX)+i*u.rd.idleAnim.sc.frameWidth, u.rd.idleAnim.sc.GetRow(cellY)
	screen.DrawImage(u.rd.spritesheet.SubImage(image.Rect(sx, sy, sx+u.rd.idleAnim.sc.frameWidth, sy+u.rd.idleAnim.sc.frameHeight)).(*ebiten.Image), op)
}