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

type LayerObjectType string

const (
	MAIN_BASE LayerObjectType = "MainBase"
	STAIRS    LayerObjectType = "Stairs"
	TOWER     LayerObjectType = "Tower"
)
