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
	CapturedBy constants.Player
	IsSpawn    bool
	Occupancy  uint8
}
