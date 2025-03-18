package entities

import (
	"github.com/ehutchllew/autoarmy/components"
	"github.com/ehutchllew/autoarmy/constants"
)

type Tile struct {
	components.Coordinates
	components.Renderable
	components.Transformable
}

func (t *Tile) Type() constants.LayerObjectType {
	return constants.TILE
}
