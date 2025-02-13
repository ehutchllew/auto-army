package components

import (
	"github.com/ehutchllew/autoarmy/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

type Coordinates struct {
	X, Y float64
}

type Dimensions struct {
	Height, Width int
}

type LayerObject struct {
	Gid constants.ID
	// Id?
	Name  constants.LayerObjectName
	Class constants.LayerObjectType
}

type Renderable struct {
	Image *ebiten.Image
}

// TODO: Add a Transformable struct?
