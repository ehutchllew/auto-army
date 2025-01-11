package scenes

import (
	"image/color"

	"github.com/ehutchllew/autoarmy/cameras"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	camera     *cameras.Camera
	tileMapImg *ebiten.Image
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := ebiten.DrawImageOptions{}

}

func (g *GameScene) FirstLoad() {
}

func (g *GameScene) IsLoaded() bool {
	panic("unimplemented")
}

func (g *GameScene) OnEnter() {
	panic("unimplemented")
}

func (g *GameScene) OnExit() {
	panic("unimplemented")
}

func (g *GameScene) Update() SceneId {
	panic("unimplemented")
}

func NewGameScene() *GameScene {
	return &GameScene{}
}

var _ Scene = (*GameScene)(nil)
