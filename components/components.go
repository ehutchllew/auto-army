package components

import (
	"github.com/ehutchllew/autoarmy/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

type Coordinates struct {
	X, Y float64
}

func (c *Coordinates) Coords() (float64, float64) {
	return c.X, c.Y
}

type Dimensions struct {
	Height, Width int
}

type LayerObject struct {
	Gid   constants.ID
	Id    constants.ID // So far for debugging only
	Name  constants.LayerObjectName
	Class constants.LayerRenderableType
}

func (lo *LayerObject) Type() constants.LayerRenderableType {
	return lo.Class
}

type Renderable struct {
	Image *ebiten.Image
}

func (r *Renderable) Img() *ebiten.Image {
	return r.Image
}

type Transformable struct {
	Tx, Ty float64
}

func (t *Transformable) TransCoords() (float64, float64) {
	return t.Tx, t.Ty
}
