package main

import (
	"github.com/ehutchllew/autoarmy/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	activeSceneId scenes.SceneId
	sceneMap      map[scenes.SceneId]scenes.Scene
}

func NewGame() *Game {
	activeSceneId := scenes.GameSceneId
	sceneMap := map[scenes.SceneId]scenes.Scene{
		scenes.GameSceneId: scenes.NewGameScene(),
	}
	sceneMap[activeSceneId].FirstLoad()

	return &Game{
		activeSceneId,
		sceneMap,
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneMap[g.activeSceneId].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
