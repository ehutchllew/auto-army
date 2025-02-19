package entities

import (
	"github.com/ehutchllew/autoarmy/components"
	"github.com/ehutchllew/autoarmy/constants"
)

type Stairs struct {
	components.Coordinates
	components.LayerObject
	components.Renderable
	components.Transformable
	Ascend  constants.CardinalDirection
	Descend constants.CardinalDirection
}
