package services

import (
	"fmt"

	"github.com/ehutchllew/autoarmy/components"
	"github.com/ehutchllew/autoarmy/constants"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Cursor struct {
	components.Renderable
	X int
	Y int
}

// NOTE: The update and draw for the mouse seems slightly laggy -- maybe need to try
// out `image.Point` updates versus trying to redraw the entire image.
func (cs *Cursor) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	// The image for the cursor is a bit weird, so requires a bit of a hacky solution to align it
	opts.GeoM.Translate(float64(cs.X)-constants.Tilesize/3, float64(cs.Y)-constants.Tilesize/3.5)
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
