package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Dx  float32
	Dy  float32
	Img *ebiten.Image
	X   float64
	Y   float64
}
