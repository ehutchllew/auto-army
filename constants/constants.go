package constants

const (
	Tilesize = 64
)

type ID uint16

type CardinalDirection string

const (
	NORTH CardinalDirection = "NORTH"
	EAST  CardinalDirection = "EAST"
	SOUTH CardinalDirection = "SOUTH"
	WEST  CardinalDirection = "WEST"
)

type Player string

const (
	BLUE   Player = "BLUE"
	GREEN  Player = "GREEN"
	NONE   Player = "NONE"
	RED    Player = "RED"
	YELLOW Player = "YELLOW"
)

type LayerObjectName string

const (
	MAIN_BASE LayerObjectName = "MainBase"
	TOWER     LayerObjectName = "Tower"
)

// These should correlate to entities
type LayerObjectType string

const (
	BUILDING LayerObjectType = "Building"
	STAIRS   LayerObjectType = "Stairs"
)
