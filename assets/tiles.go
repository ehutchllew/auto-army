package assets

import (
	"image"

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
	Path string `json:"image"`
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
