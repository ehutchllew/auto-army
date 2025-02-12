package constants

const (
	Tilesize = 64
)

type ID uint16

type PLAYER string

const (
	BLUE   PLAYER = "BLUE"
	GREEN  PLAYER = "GREEN"
	NONE   PLAYER = "NONE"
	RED    PLAYER = "RED"
	YELLOW PLAYER = "YELLOW"
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
