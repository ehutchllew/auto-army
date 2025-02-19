package entities

import "github.com/ehutchllew/autoarmy/components"

type Cliff struct {
	components.Coordinates
	components.LayerObject
	components.Renderable
	components.Transformable
}
