package scenes

import (
	"github.com/ehutchllew/autoarmy/entities"
	"github.com/hajimehoshi/ebiten/v2"
)

type SceneId uint

const (
	GameSceneId SceneId = iota
	PauseSceneId
	StartSceneId
	ExitSceneId
)

type Scene interface {
	Draw(screen *ebiten.Image)
	FirstLoad()
	IsLoaded() bool
	OnEnter()
	OnExit()
	Update() SceneId
}

type RenderIndexObjects struct {
	LayerZIndices []uint8
	Objects       map[uint8][]entities.IEntity
}

type ZIndexObjects struct {
	LayerZIndices []uint8
	Objects       map[uint8]map[string]entities.IEntity
}
