package components

import "github.com/ehutchllew/autoarmy/constants"

type Coordinates struct {
	X, Y float64
}

type Dimensions struct {
	Height, Width uint16
}

type LayerObjectProperty struct {
	Name       string
	ProperType string
	Type       string
	Value      string
}

type LayerObject struct {
	Coordinates
	Dimensions
	Gid constants.ID
	// Id?
	Properties []LayerObjectProperty
	Type       constants.LayerObjectType
}
