package assets

import (
	"encoding/json"
	"fmt"
	"image"
	"os"

	"github.com/ehutchllew/autoarmy/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

type Tileset interface {
	Img(id uint8) *ebiten.Image
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

// TODO: test to see if this rawMsg approach works to dynamically
// set the tileset type
func NewTileset(tp string, gid constants.ID) (*Tileset, error) {
	content, err := os.ReadFile(tp)
	if err != nil {
		return nil, fmt.Errorf("Error reading file at path: (%s) -- Error: %w", tp, err)
	}

	var rawMsg map[string]interface{}
	if err := json.Unmarshal(content, rawMsg); err != nil {
		return nil, fmt.Errorf("Error unmarshalling tileset: %w", err)
	}

	fmt.Printf("JSON: %+v", rawMsg)

	return nil, nil
}
