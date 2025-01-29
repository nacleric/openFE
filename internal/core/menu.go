package core

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/math/f64"
)

type MenuManager struct {
	// MenuStack  []int // Enums to keep track of which menu's are on top of each other
	ActionMenu ActionMenu
}

type ActionMenu struct {
	MenuOptions []string
	Selected    int // Index of selected option
	rds         []RenderData
	//fadeInFrames int // Fade-in duration
	//frameCount   int // Tracks number of elapsed frames for fade-in
}

func CreateActionMenu(spritesheet *ebiten.Image) ActionMenu {
	idleAnimData0 := AnimationData{SpriteCell{0, 0, 16, 16}, 5, 16}
	idleAnimData1 := AnimationData{SpriteCell{0, 1, 16, 16}, 5, 16}
	idleAnimData2 := AnimationData{SpriteCell{0, 2, 16, 16}, 5, 16}

	icon0_rd := RenderData{
		ad:          idleAnimData0,
		spritesheet: spritesheet,
	}
	icon1_rd := RenderData{
		ad:          idleAnimData1,
		spritesheet: spritesheet,
	}

	icon2_rd := RenderData{
		ad:          idleAnimData2,
		spritesheet: spritesheet,
	}

	rds := []RenderData{icon0_rd, icon1_rd, icon2_rd}

	actionMenu := ActionMenu{MenuOptions: []string{"attack", "items", "skip"}, Selected: 0, rds: rds}
	return actionMenu
}

func (m *ActionMenu) Update() {
	isArrowLeftPressed := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft)
	if isArrowLeftPressed {
		if m.Selected != 0 {
			m.Selected = (m.Selected - 1 + len(m.MenuOptions)) % len(m.MenuOptions)
		}
		isArrowLeftPressed = false
	}

	isArrowRightPressed := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight)
	if isArrowRightPressed {
		if m.Selected != len(m.MenuOptions)-1 {
			m.Selected = (m.Selected + 1) % len(m.MenuOptions)
		}
		isArrowRightPressed = false
	}
}

func (m *ActionMenu) DrawMenu(screen *ebiten.Image, x0y0 f64.Vec2, offsetX, offsetY float64, count int) {
	f32cameraScale := float32(CAMERASCALE)
	f32offsetX := float32(offsetX)
	f32offsetY := float32(offsetY)
	padX := float32(6 * f32cameraScale)
	padY := float32(20 * f32cameraScale)
	gap := float32(1 * f32cameraScale)
	padXgap := float32(8*f32cameraScale) + gap

	startLocationX := float32(x0y0[X]) + f32offsetX - padX
	color := color.RGBA{R: 25, G: 0, B: 255, A: 5}

	square0 := []float32{startLocationX, float32(x0y0[Y]) + f32offsetY + padY}
	square1 := []float32{startLocationX + padXgap, float32(x0y0[Y]) + f32offsetY + padY}
	square2 := []float32{startLocationX + (padXgap * 2), float32(x0y0[Y]) + f32offsetY + padY}
	vector.DrawFilledRect(screen, square0[X], square0[Y], 8*f32cameraScale, 8*f32cameraScale, color, true)
	vector.DrawFilledRect(screen, square1[X], square1[Y], 8*f32cameraScale, 8*f32cameraScale, color, true)
	vector.DrawFilledRect(screen, square2[X], square2[Y], 8*f32cameraScale, 8*f32cameraScale, color, true)

	m.IdleAnimation(screen, count, square0[X], square0[Y], 0)
	m.IdleAnimation(screen, count, square1[X], square1[Y], 1)
	m.IdleAnimation(screen, count, square2[X], square2[Y], 2)
}

func (m *ActionMenu) IdleAnimation(screen *ebiten.Image, count int, x0, y0 float32, index int) {
	op := &ebiten.DrawImageOptions{}

	cellX := m.rds[index].ad.sc.cellX
	cellY := m.rds[index].ad.sc.cellY

	if m.Selected == index {
		op.GeoM.Scale(float64(CAMERASCALE/1.75), float64(CAMERASCALE/1.75))
		op.GeoM.Translate(float64(x0)-CAMERASCALE, float64(y0)-CAMERASCALE)

		i := (count / m.rds[index].ad.frameFrequency) % m.rds[index].ad.frameCount
		sx, sy := m.rds[index].ad.sc.GetCol(cellX)+i*m.rds[index].ad.sc.frameWidth, m.rds[index].ad.sc.GetRow(cellY)
		screen.DrawImage(m.rds[index].spritesheet.SubImage(image.Rect(sx, sy, sx+m.rds[index].ad.sc.frameWidth, sy+m.rds[index].ad.sc.frameHeight)).(*ebiten.Image), op)
	} else {
		op.GeoM.Scale(float64(CAMERASCALE/2), float64(CAMERASCALE/2))
		op.GeoM.Translate(float64(x0), float64(y0))

		i := 0
		sx, sy := m.rds[index].ad.sc.GetCol(cellX)+i*m.rds[index].ad.sc.frameWidth, m.rds[index].ad.sc.GetRow(cellY)
		screen.DrawImage(m.rds[index].spritesheet.SubImage(image.Rect(sx, sy, sx+m.rds[index].ad.sc.frameWidth, sy+m.rds[index].ad.sc.frameHeight)).(*ebiten.Image), op)
	}
}
