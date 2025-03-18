package scenes

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"log"
	"slices"
	"strconv"

	"github.com/ehutchllew/autoarmy/assets"
	"github.com/ehutchllew/autoarmy/cameras"
	"github.com/ehutchllew/autoarmy/components"
	"github.com/ehutchllew/autoarmy/constants"
	"github.com/ehutchllew/autoarmy/entities"
	"github.com/ehutchllew/autoarmy/services"
	"github.com/ehutchllew/autoarmy/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameScene struct {
	*services.Cursor
	camera          *cameras.Camera
	renderIndexObjs *RenderIndexObjects
	tileMapJson     *assets.TileMapJson
	tilesets        []assets.Tileset
	zIndexObjs      *ZIndexObjects
}

var (
	fontSource *text.GoTextFaceSource
	fontFace   *text.GoTextFace
)

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := ebiten.DrawImageOptions{}

	g.drawMap(screen, &opts)
	g.Cursor.Draw(screen)
}

func (g *GameScene) FirstLoad() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(assets.DepartMono_otf))
	if err != nil {
		fmt.Printf("Unable to generate font source: %v", err)
	}

	fontSource = s
	fontFace = &text.GoTextFace{
		Source: fontSource,
		Size:   16,
	}

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
	g.renderIndexObjs, g.zIndexObjs = g.firstLoadObjectState()
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
	for i := 0; i < len(g.renderIndexObjs.LayerZIndices)-1; i++ {
		objects := g.renderIndexObjs.Objects[uint8(i)]
		for _, o := range objects {
			opts.GeoM.Translate(o.TransCoords())
			screen.DrawImage(o.Img(), opts)
			opts.GeoM.Reset()
			renderBuildingBanner(o, screen, opts)
		}
	}
}

func (g *GameScene) firstLoadObjectState() (*RenderIndexObjects, *ZIndexObjects) {
	layerZIndices := make([]uint8, len(g.tileMapJson.Layers))

	var ri = &RenderIndexObjects{
		Objects: make(map[uint8][]entities.IEntity),
	}
	var zi = &ZIndexObjects{
		Objects: make(map[uint8]map[string]entities.IEntity),
	}

	for _, layer := range g.tileMapJson.Layers {
		z, err := strconv.ParseUint(layer.ZIndex, 10, 8)
		if err != nil {
			// NOTE: Potentially log fatal instead?
			fmt.Printf("Error parsing Layer Z-Index: %v", err)
		}
		currentZ := uint8(z)
		// We don't want dupes because we'll range over these values
		// to access each associated key in the `Objects` map for
		// rendering and/or mutating/updating.
		if !slices.Contains(layerZIndices, currentZ) {
			layerZIndices = append(layerZIndices, currentZ)
		}

		// Need to define tileset within this scope to ensure it resets for each layer
		var tileset assets.Tileset

		/*
		* TILES
		 */
		// idx equals index of the actual data in the slice
		// tileId equals the ID of the tile on its tileset file
		for idx, tileId := range layer.Data {
			// If there is no tile, then skip this iteration
			if tileId == 0 {
				continue
			}

			// Get associated tileset if tileset is nil
			if tileset == nil {
				// Loop backwards since we want to match the tile GID with
				// the highest possible `firstgid` of the tilesets
				for i := len(g.tilesets) - 1; i >= 0; i-- {
					t := g.tilesets[i]
					if tileId >= int(t.Gid()) {
						tileset = t
						break
					}
				}
			}

			// Get tile coordinates on tileset image
			x := idx % layer.Width
			y := idx / layer.Width

			// Multiply by the Tilesize our game uses to get coords
			x *= constants.Tilesize
			y *= constants.Tilesize

			fX := float64(x)
			fY := float64(y)

			img := tileset.Img(constants.ID(tileId))
			tile := &entities.Tile{
				Coordinates: components.Coordinates{
					X: fX,
					Y: fY,
				},
				Renderable: components.Renderable{
					Image: img,
				},
				Transformable: components.Transformable{
					Tx: fX,
					Ty: fY,
				},
			}

			ri.Objects[currentZ] = append(ri.Objects[currentZ], tile)
		}
		tileset = nil

		/*
		* OBJECTS
		 */
		// These zObjects go in the ZIndexObjects struct
		zObjects := make(map[string]entities.IEntity)
		// Work backwards to ensure render order when adding objects to `ri`
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
				// FIXME: #1 in `todo.txt` (convert to tile layer)
				fmt.Printf("Unable to unpack object :: Error: \n %v", err)
				continue
			}

			ri.Objects[currentZ] = append(ri.Objects[currentZ], object)

			x, y := object.Coords()
			zObjects[fmt.Sprintf("%.0f,%.0f", x, y)] = object
		}

		zi.Objects[currentZ] = zObjects

	}
	ri.LayerZIndices = layerZIndices
	zi.LayerZIndices = layerZIndices
	return ri, zi
}

// TODO: Think about eliminating `FirstLoad` and putting that logic here
func NewGameScene() *GameScene {
	return &GameScene{
		Cursor: services.NewCursorService("./assets/ui/cursor_0.png"),
	}
}

func assignObject(obj assets.TileMapObjectsJson, tileset assets.Tileset) (entities.IEntity, error) {
	coercedType := constants.LayerObjectType(obj.Type)

	switch coercedType {
	case constants.BUILDING:
		return assignBuilding(obj, tileset)
	case constants.CLIFF:
		return assignCliff(obj, tileset)
	case constants.STAIRS:
		return assignStairs(obj, tileset)
	}

	return nil, fmt.Errorf("Unsupported object type: (%v)", obj.Type)
}

func assignBuilding(obj assets.TileMapObjectsJson, tileset assets.Tileset) (*entities.Building, error) {
	objProps := make(map[string]any)

	for _, p := range obj.Properties {
		objProps[p.Name] = p.Value
	}

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

	img := tileset.Img(obj.Gid)

	tx, ty := obj.X, obj.Y
	ty -= float64(img.Bounds().Dy())

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
			Class: constants.LayerObjectType(obj.Type),
			Gid:   obj.Gid,
			Id:    obj.Id,
			Name:  constants.LayerObjectName(obj.Name),
		},
		Renderable: components.Renderable{
			Image: img,
		},
		Transformable: components.Transformable{
			Tx: tx,
			Ty: ty,
		},
		Capacity:   capacity,
		CapturedBy: constants.Player(capBy),
		IsSpawn:    isSpawn,
		Occupancy:  occ,
	}, nil
}

func assignCliff(obj assets.TileMapObjectsJson, tileset assets.Tileset) (*entities.Cliff, error) {
	img := tileset.Img(obj.Gid)

	tx, ty := obj.X, obj.Y
	ty -= float64(img.Bounds().Dy())

	return &entities.Cliff{
		Coordinates: components.Coordinates{
			X: obj.X,
			Y: obj.Y,
		},
		LayerObject: components.LayerObject{
			Class: constants.LayerObjectType(obj.Type),
			Gid:   obj.Gid,
			Id:    obj.Id,
			Name:  constants.LayerObjectName(obj.Name),
		},
		Renderable: components.Renderable{
			Image: img,
		},
		Transformable: components.Transformable{
			Tx: tx,
			Ty: ty,
		},
	}, nil
}

func assignStairs(obj assets.TileMapObjectsJson, tileset assets.Tileset) (*entities.Stairs, error) {
	objProps := make(map[string]any)

	for _, p := range obj.Properties {
		objProps[p.Name] = p.Value
	}

	ascend := utils.SafeConvertString(objProps["ascend"])
	if ascend == "" {
		return nil, errors.New("Stairs object is missing `ascend` property")
	}

	descend := utils.SafeConvertString(objProps["descend"])
	if descend == "" {
		return nil, errors.New("Stairs object is missing `descend` property")
	}

	img := tileset.Img(obj.Gid)

	tx, ty := obj.X, obj.Y
	ty -= float64(img.Bounds().Dy())

	return &entities.Stairs{
		Coordinates: components.Coordinates{
			X: obj.X,
			Y: obj.Y,
		},
		LayerObject: components.LayerObject{
			Class: constants.LayerObjectType(obj.Type),
			Gid:   obj.Gid,
			Id:    obj.Id,
			Name:  constants.LayerObjectName(obj.Name),
		},
		Renderable: components.Renderable{
			Image: img,
		},
		Transformable: components.Transformable{
			Tx: tx,
			Ty: ty,
		},
		Ascend:  constants.CardinalDirection(ascend),
		Descend: constants.CardinalDirection(descend),
	}, nil
}

func renderBuildingBanner(o entities.IEntity, screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// Check if building has occupancy & capacity, then display banner "O/C"
	if o.Type() == constants.BUILDING {
		coObj := o.(*entities.Building)
		if coObj.IsSpawn {
			tx, ty := o.TransCoords()
			scaleAmount := 0.80
			opts.GeoM.Scale(scaleAmount, scaleAmount)
			opts.GeoM.Translate(tx+float64(o.Img().Bounds().Dx())/2, ty)
			switch coObj.CapturedBy {
			case constants.BLUE:
				capBanner, _, err := ebitenutil.NewImageFromFile("./assets/ui/ribbon_blue.png")
				if err != nil {
					fmt.Printf("Unable to parse image: %v", err)
				}
				opts.GeoM.Translate(-float64(capBanner.Bounds().Dx())*scaleAmount/2, 0.0)
				screen.DrawImage(capBanner, opts)
			case constants.RED:
				capBanner, _, err := ebitenutil.NewImageFromFile("./assets/ui/ribbon_red.png")
				if err != nil {
					fmt.Printf("Unable to parse image: %v", err)
				}
				opts.GeoM.Translate(-float64(capBanner.Bounds().Dx())*scaleAmount/2, 0.0)
				screen.DrawImage(capBanner, opts)
			default:
				capBanner, _, err := ebitenutil.NewImageFromFile("./assets/ui/ribbon_gray.png")
				if err != nil {
					fmt.Printf("Unable to parse image: %v", err)
				}
				opts.GeoM.Translate(-float64(capBanner.Bounds().Dx())*scaleAmount/2, 0.0)
				screen.DrawImage(capBanner, opts)
			}

			textW, textH := text.Measure(fmt.Sprintf("%d/%d", coObj.Occupancy, coObj.Capacity), fontFace, 0)
			tOpts := &text.DrawOptions{}
			tOpts.GeoM.Translate(tx+float64(o.Img().Bounds().Dx())/2-textW/2, ty+(textH/4))
			tOpts.ColorScale.Scale(0, 0, 0, 1)
			text.Draw(screen, fmt.Sprintf("%d/%d", coObj.Occupancy, coObj.Capacity), fontFace, tOpts)
		}
		opts.GeoM.Reset()
	}
}

var _ Scene = (*GameScene)(nil)
