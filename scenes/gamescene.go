package scenes

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ehutchllew/autoarmy/assets"
	"github.com/ehutchllew/autoarmy/cameras"
	"github.com/ehutchllew/autoarmy/components"
	"github.com/ehutchllew/autoarmy/constants"
	"github.com/ehutchllew/autoarmy/entities"
	"github.com/ehutchllew/autoarmy/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	camera      *cameras.Camera
	objects     map[string]entities.IEntity
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
	g.objects = make(map[string]entities.IEntity)
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

			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(img, opts)

			opts.GeoM.Reset()
		}

		tileset = nil
		for i := len(layer.Objects) - 1; i >= 0; i-- {
			obj := layer.Objects[i]
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
			object, err := assignObject(obj, tileset)
			if err != nil {
				continue
				log.Fatalf("Unable to unpack object :: Error: \n %+v", err)
			}

			x, y := object.Coords()
			g.objects[fmt.Sprintf("%f,%f", x, y)] = object

			// TODO: transform / assign transform? The image then draw on screen
			opts.GeoM.Translate(object.Coords())
			opts.GeoM.Translate(0.0, -(float64(object.Img().Bounds().Dy()))+constants.Tilesize)
			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(object.Img(), opts)

			opts.GeoM.Reset()

		}
	}
}

func NewGameScene() *GameScene {
	return &GameScene{}
}

func assignObject(obj assets.TileMapObjectsJson, tileset assets.Tileset) (entities.IEntity, error) {
	coercedType := constants.LayerObjectType(obj.Type)
	objProps := make(map[string]any)

	for _, p := range obj.Properties {
		objProps[p.Name] = p.Value
	}

	switch coercedType {
	case constants.BUILDING:
		capacity, err := utils.SafeConvertUint8(objProps["capacity"])
		if err != nil {
			return nil, err
		}

		capBy := utils.SafeConvertString(objProps["captured_by"])

		isSpawn := utils.SafeConvertBool(objProps["is_spawn"])

		occ, err := utils.SafeConvertUint8(objProps["occupancy"])
		if err != nil {
			return nil, err
		}

		return &entities.Building{
			Coordinates: components.Coordinates{
				X: obj.X,
				Y: obj.Y,
			},
			Dimensions: components.Dimensions{
				Height: obj.Height,
				Width:  obj.Width,
			},
			LayerObject: components.LayerObject{
				Class: coercedType,
				Gid:   obj.Gid,
				Id:    obj.Id,
				Name:  constants.LayerObjectName(obj.Name),
			},
			Capacity:   capacity,
			CapturedBy: constants.PLAYER(capBy),
			Renderable: components.Renderable{
				Image: tileset.Img(obj.Gid), // Looks like obj.Gid - tilset.Gid results in the actual PNG id
			},
			IsSpawn:   isSpawn,
			Occupancy: occ,
		}, nil
	}

	return nil, fmt.Errorf("Unsupported object type: (%v)", obj.Type)
}

var _ Scene = (*GameScene)(nil)
