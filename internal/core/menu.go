package core

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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
	//fadeInFrames int // Fade-in duration
	//frameCount   int // Tracks number of elapsed frames for fade-in
}

func (m *ActionMenu) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		m.Selected = (m.Selected + 1) % len(m.MenuOptions)
		fmt.Println("in here, arrowRight clicked")
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		m.Selected = (m.Selected - 1 + len(m.MenuOptions)) % len(m.MenuOptions)
		fmt.Println("in here, arrowLeft clicked")

	}

	/*
		if m.frameCount < m.fadeInFrames {
			m.frameCount++
		}
	*/
}

func (m *ActionMenu) Reset() {
	m.Selected = 0
}

func (m *ActionMenu) Draw(screen *ebiten.Image, x0y0 f64.Vec2, offsetX, offsetY float64) {
	f32cameraScale := float32(CAMERASCALE)
	f32offsetX := float32(offsetX)
	f32offsetY := float32(offsetY)
	padX := float32(6 * f32cameraScale)
	padY := float32(20 * f32cameraScale)
	gap := float32(1 * f32cameraScale)
	padXgap := float32(8*f32cameraScale) + gap

	startLocationX := float32(x0y0[X]) + f32offsetX - padX
	color := color.RGBA{R: 25, G: 0, B: 255, A: 5}
	vector.DrawFilledRect(screen, startLocationX, float32(x0y0[Y])+f32offsetY+padY, 8*f32cameraScale, 8*f32cameraScale, color, true)
	vector.DrawFilledRect(screen, startLocationX+padXgap, float32(x0y0[Y])+f32offsetY+padY, 8*f32cameraScale, 8*f32cameraScale, color, true)
	vector.DrawFilledRect(screen, startLocationX+(padXgap*2), float32(x0y0[Y])+f32offsetY+padY, 8*f32cameraScale, 8*f32cameraScale, color, true)

}

/*
func (m *ActionMenu) Draw(screen *ebiten.Image, x0y0 f64.Vec2) {
	alpha := 1.0
	if m.frameCount < m.fadeInFrames {
		alpha = float64(m.frameCount) / float64(m.fadeInFrames)
	}

	for i, option := range m.menuOptions {
		clr := color.RGBA{255, 255, 255, uint8(255 * alpha)}
		if i == m.selected {
			clr = color.RGBA{255, 200, 0, uint8(255 * alpha)}
		}

		x, y := 100, 100+i*30 // Calculate position
		text.Draw(screen, option, face, x, y, clr)

	}
}
*/
