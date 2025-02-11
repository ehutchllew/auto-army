package scenes

import (
	"image"
	"image/color"
	"log"

	"github.com/ehutchllew/autoarmy/assets"
	"github.com/ehutchllew/autoarmy/cameras"
	"github.com/ehutchllew/autoarmy/constants"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameScene struct {
	camera      *cameras.Camera
	tileMapJson *assets.TileMapJson
	tilesets    []assets.Tileset
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := ebiten.DrawImageOptions{}

	g.drawMap(screen, &opts)
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

	g.camera = cameras.NewCamera(0.0, 0.0)
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

func (g *GameScene) drawMap(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	for _, layer := range g.tileMapJson.Layers {
		var tileset assets.Tileset
		for tileI, tileId := range layer.Data {
			if tileId == 0 {
				continue
			}

			// Get associated tileset if tileset is nil
			if tileset == nil {
				for i := len(g.tilesets) - 1; i >= 0; i-- {
					t := g.tilesets[i]
					if tileId >= int(t.Gid()) {
						tileset = t
						break
					}
				}
			}

			// Get tile index on tileset
			x := tileI % layer.Width
			y := tileI / layer.Width

			x *= constants.Tilesize
			y *= constants.Tilesize

			img := tileset.Img(constants.ID(tileId))

			opts.GeoM.Translate(float64(x), float64(y))

			var dT *assets.DynamicTileset
			if tileset.Type() == assets.DynamicType {
				dT = tileset.(*assets.DynamicTileset)
				/*
				* TODO: Figure out best way to attach a translation vector to the actual
				* tileset and not just the image.
				 */
				if dT.Collider == nil {
					collider := image.Rect(
						x,
						y-(img.Bounds().Dy()+constants.Tilesize),
						x-img.Bounds().Dx(),
						x-img.Bounds().Dy(),
					)
					colliderP := &collider
					dT.Collider = colliderP

				}
				opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()))+constants.Tilesize)
			}
			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(img, opts)

			if dT != nil {
				vector.StrokeRect(
					screen,
					float32(dT.Collider.Min.X),
					float32(dT.Collider.Min.Y),
					float32(dT.Collider.Dx()),
					float32(dT.Collider.Dy()),
					1.0,
					color.RGBA{255, 0, 0, 255},
					true,
				)
			}

			opts.GeoM.Reset()
		}

		tileset = nil
		for _, obj := range layer.Objects {
			if tileset == nil {
				for i := len(g.tilesets) - 1; i >= 0; i-- {
					t := g.tilesets[i]
					if obj.Gid >= t.Gid() {
						tileset = t
						break
					}
				}
			}

			// Assign object and its properties to a struct
		}
	}
}

func NewGameScene() *GameScene {
	return &GameScene{}
}

var _ Scene = (*GameScene)(nil)
