package entities

import (
	"github.com/ehutchllew/autoarmy/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

type IEntity interface {
	Coords() (float64, float64)
	Img() *ebiten.Image
	TransCoords() (float64, float64)
	Type() constants.LayerObjectType
}
