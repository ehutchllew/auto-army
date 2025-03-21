package assets

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/ehutchllew/autoarmy/constants"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TilesetType uint8

const (
	UniformType TilesetType = iota
	DynamicType
)

type Tileset interface {
	Gid() constants.ID
	Img(id constants.ID) *ebiten.Image
	Type() TilesetType
}

type UniformTileset struct {
	columns uint8
	gid     constants.ID
	img     *ebiten.Image
}

type UniformTilesetJson struct {
	Columns uint8  `json:"columns"`
	Path    string `json:"image"`
}

type DynamicTileset struct {
	gid  constants.ID
	imgs []*ebiten.Image
}

type DynamicTilesetTile struct {
	id          constants.ID
	image       string
	imageHeight uint16
	imageWidth  uint16
}

type DynamicTilesetJson struct {
	Tiles []*DynamicTilesetTile `json:"tiles"`
}

type GenericTileset[T UniformTilesetJson | DynamicTilesetJson] struct {
	Data *T
}

func (u *UniformTileset) Gid() constants.ID {
	return u.gid
}

func (u *UniformTileset) Img(id constants.ID) *ebiten.Image {
	// The ID that gets passed in is the global ID used by the `map*.json`. To get the real ID of the actual image from its associated tileset we need to subtract the "firstGid" of the tilemap tileset from the passed in global ID.

	realId := id - u.gid
	srcX := int((uint16(realId) % uint16(u.columns)) * constants.Tilesize)
	srcY := int((uint16(realId) / uint16(u.columns)) * constants.Tilesize)

	return u.img.SubImage(
		image.Rect(srcX, srcY, srcX+constants.Tilesize, srcY+constants.Tilesize),
	).(*ebiten.Image)
}

func (u *UniformTileset) Type() TilesetType {
	return UniformType
}

func (d *DynamicTileset) Gid() constants.ID {
	return d.gid
}

func (d *DynamicTileset) Img(id constants.ID) *ebiten.Image {
	realId := id - d.gid

	return d.imgs[realId]
}

func (d *DynamicTileset) Type() TilesetType {
	return DynamicType
}

func NewTileset(tp string, gid constants.ID) (Tileset, error) {
	content, err := os.ReadFile(tp)
	if err != nil {
		return nil, fmt.Errorf("Error reading file at path: (%s) -- Error: %w", tp, err)
	}

	var rawMsg map[string]any
	if err := json.Unmarshal(content, &rawMsg); err != nil {
		return nil, fmt.Errorf("Error unmarshalling tileset: %w", err)
	}

	var tileset Tileset
	isValidTileset := false
	if v, ok := rawMsg["columns"]; ok && v != 0.0 {
		isValidTileset = true

		uT := &GenericTileset[UniformTilesetJson]{
			Data: &UniformTilesetJson{
				Columns: uint8(rawMsg["columns"].(float64)),
				Path:    rawMsg["image"].(string),
			},
		}

		imgPath := filepath.Clean(uT.Data.Path)
		imgPath = strings.ReplaceAll(imgPath, "\\", "/")
		imgPath = strings.TrimPrefix(imgPath, "../")
		imgPath = filepath.Join("assets/", imgPath)
		img, _, err := ebitenutil.NewImageFromFile(imgPath)
		if err != nil {
			return nil, fmt.Errorf("UniformTileset: Unable to create image from file at path: (%s) -- Error: %w", imgPath, err)
		}

		tileset = &UniformTileset{
			columns: uint8(rawMsg["columns"].(float64)),
			gid:     gid,
			img:     img,
		}
	}

	if tiles, ok := rawMsg["tiles"]; ok {
		isValidTileset = true

		tilesSlice := tiles.([]interface{})

		convertedTiles := make([]*DynamicTilesetTile, len(tilesSlice))

		for i, t := range tilesSlice {
			tileMap := t.(map[string]interface{})
			convertedTiles[i] = &DynamicTilesetTile{
				image: tileMap["image"].(string),
			}
		}

		dT := &GenericTileset[DynamicTilesetJson]{
			Data: &DynamicTilesetJson{
				Tiles: convertedTiles,
			},
		}

		imgs := make([]*ebiten.Image, 0)
		for _, tile := range dT.Data.Tiles {
			imgPath := filepath.Clean(tile.image)
			imgPath = strings.ReplaceAll(imgPath, "\\", "/")
			imgPath = strings.TrimPrefix(imgPath, "../")
			imgPath = filepath.Join("assets/", imgPath)
			img, _, err := ebitenutil.NewImageFromFile(imgPath)
			if err != nil {
				return nil, fmt.Errorf("DynamicTileset: Unable to create image from file at path: (%s) -- Error: %w", imgPath, err)
			}

			imgs = append(imgs, img)
		}

		tileset = &DynamicTileset{
			gid:  gid,
			imgs: imgs,
		}
	}

	if !isValidTileset {
		return nil, fmt.Errorf("Parsed JSON file at path: (%s) is not a valid tileset", tp)
	}

	return tileset, nil
}
