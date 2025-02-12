package entities

import "github.com/ehutchllew/autoarmy/constants"

type IEntity interface {
	Type() constants.LayerObjectType
}
