package entities

import (
	"github.com/ehutchllew/autoarmy/components"
	"github.com/ehutchllew/autoarmy/constants"
)

type Building struct {
	components.LayerObject
	Capacity   uint8
	CapturedBy constants.PLAYER
	IsSpawn    bool
	Length     uint8
}
