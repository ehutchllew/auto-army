package entities

import (
	"github.com/ehutchllew/autoarmy/components"
	"github.com/ehutchllew/autoarmy/constants"
)

type Building struct {
	components.Coordinates
	components.Dimensions
	components.LayerObject
	components.Renderable
	components.Transformable
	Capacity   uint8
	CapturedBy constants.PLAYER
	IsSpawn    bool
	Occupancy  uint8
}

// func (b *Building) Coords() (float64, float64) {
// 	return b.X, b.Y
// }

// func (b *Building) Img() *ebiten.Image {
// 	return b.Image
// }

// func (b *Building) Type() constants.LayerObjectType {
// 	return b.Class
// }
