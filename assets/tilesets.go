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
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id constants.ID) *ebiten.Image
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

type DynamicTilesJson struct {
	Height uint16       `json:"imageheight"`
	Id     constants.ID `json:"id"`
	Path   string       `json:"image"`
	Width  uint16       `json:"imageWidth"`
}

type DynamicTilesetJson struct {
	Tiles []*DynamicTilesJson `json:"tiles"`
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

func (d *DynamicTileset) Img(id constants.ID) *ebiten.Image {
	realId := id - d.gid
	return d.imgs[realId]
}

func NewTileset(tp string, gid constants.ID) (*Tileset, error) {
	content, err := os.ReadFile(tp)
	if err != nil {
		return nil, fmt.Errorf("Error reading file at path: (%s) -- Error: %w", tp, err)
	}

	var rawMsg map[string]interface{}
	if err := json.Unmarshal(content, &rawMsg); err != nil {
		return nil, fmt.Errorf("Error unmarshalling tileset: %w", err)
	}

	var tileset Tileset
	if rawMsg["columns"] != 0.0 {
		imgPath := filepath.Clean(rawMsg["image"].(string))
		imgPath = strings.ReplaceAll(imgPath, "\\", "/")
		imgPath = strings.TrimPrefix(imgPath, "../")
		imgPath = filepath.Join("assets/", imgPath)
		img, _, err := ebitenutil.NewImageFromFile(imgPath)
		if err != nil {
			return nil, fmt.Errorf("Unable to create image from file at path: (%s) -- Error: %w", imgPath, err)
		}

		tileset = &UniformTileset{
			columns: uint8(rawMsg["columns"].(float64)),
			gid:     gid,
			img:     img,
		}
	} else {
		// TODO: Do DynamicTilesetJson
	}

	return &tileset, nil
}
