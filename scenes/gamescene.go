package scenes

import (
	"image/color"
	"log"

	"github.com/ehutchllew/autoarmy/assets"
	"github.com/ehutchllew/autoarmy/cameras"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	camera      *cameras.Camera
	tileMapJson *assets.TileMapJson
	tilesets    []*assets.Tileset
}

// TODO: now with tileMapJson and tilesets we can draw this
func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	// opts := ebiten.DrawImageOptions{}
	//
	// g.drawMap(screen, &opts)

}

func (g *GameScene) FirstLoad() {
	tileMapJson, err := assets.NewTileMapJson("./assets/maps/map1.json")
	if err != nil {
		log.Fatalf("Unable to load Tilemap JSON: %v", err)
	}

	tilesets, err := tileMapJson.GenTilesets()
	if err != nil {
		log.Fatalf("Unable to load tilesets: %v", err)
	}

	g.tileMapJson = tileMapJson
	g.tilesets = tilesets
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

// func (g *GameScene) drawMap(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
// 	for layerI, layer := range g.tileMapJson.Layers {
// 		for tileI, tile := range layer.Data {
// 			if tileI == 0 {
// 				continue
// 			}
//
// 			// Get tile index on tileset
// 			x := tileI % layer.Width
// 			y := tileI / layer.Width
//
// 			x *= constants.Tilesize
// 			y *= constants.Tilesize
//
// 			tileset := g.tileMapJson.Tilesets[layerI]
//
// 			opts.GeoM.Translate(float64(x), float64(y))
// 			opts.GeoM.Translate(0.0, -(float64(tileset)))
// 		}
// 	}
// }

func NewGameScene() *GameScene {
	return &GameScene{}
}

var _ Scene = (*GameScene)(nil)
