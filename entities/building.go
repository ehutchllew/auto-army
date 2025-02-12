package entities

import (
	"github.com/ehutchllew/autoarmy/components"
	"github.com/ehutchllew/autoarmy/constants"
)

type Building struct {
	components.Coordinates
	components.Dimensions
	components.LayerObject
	Capacity   uint8
	CapturedBy constants.PLAYER
	IsSpawn    bool
	Occupancy  uint8
}

func (b *Building) Type() constants.LayerObjectType {
	return b.Class
}
