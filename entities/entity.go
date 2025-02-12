package entities

import "github.com/ehutchllew/autoarmy/constants"

type IEntity interface {
	Coords() (float64, float64)
	Type() constants.LayerObjectType
}
