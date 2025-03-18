package services

import (
	"fmt"

	"github.com/ehutchllew/autoarmy/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Cursor struct {
	components.Renderable
	X int
	Y int
}

func (cs *Cursor) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(cs.X), float64(cs.Y))
	screen.DrawImage(cs.Renderable.Image, &opts)
	opts.GeoM.Reset()
}

func NewCursorService(imgPath string) *Cursor {
	c, _, err := ebitenutil.NewImageFromFile(imgPath)
	if err != nil {
		fmt.Printf("Unable to parse image: %v\n", err)
	}

	x, y := ebiten.CursorPosition()

	return &Cursor{
		Renderable: components.Renderable{
			Image: c,
		},
		X: x,
		Y: y,
	}
}

func (cs *Cursor) Position() (int, int) {
	return cs.X, cs.Y
}

func (cs *Cursor) Update() {
	cs.X, cs.Y = ebiten.CursorPosition()
}
