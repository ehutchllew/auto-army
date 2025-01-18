package assets

import (
	"encoding/json"
	"os"
	"path"

	"github.com/ehutchllew/autoarmy/constants"
)

type TileMapLayerJson struct {
	Data   []int  `json:"data"`
	Height int    `json:"height"`
	Name   string `json:"name"`
	Width  int    `json:"width"`
}

type TileMapTilesetJson struct {
	Firstgid uint8  `json:"firstgid"`
	Source   string `json:"source"`
}

type TileMapJson struct {
	Layers   []TileMapLayerJson   `json:"layers"`
	Tilesets []TileMapTilesetJson `json:"tilesets"`
}

func (t *TileMapJson) GenTilesets() ([]*Tileset, error) {
	ts := make([]*Tileset, 0)
	for _, tilesetData := range t.Tilesets {
		tilesetPath := path.Join("assets/maps/", tilesetData.Source)
		tileset, err := NewTileset(tilesetPath, constants.ID(tilesetData.Firstgid))
		if err != nil {
			return nil, err
		}

		ts = append(ts, tileset)
	}

	return ts, nil
}

func NewTileMapJson(fp string) (*TileMapJson, error) {
	contents, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var tileMapJson TileMapJson
	err = json.Unmarshal(contents, &tileMapJson)
	if err != nil {
		return nil, err
	}

	return &tileMapJson, nil
}
