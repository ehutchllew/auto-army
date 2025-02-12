package assets

import (
	"encoding/json"
	"os"
	"path"

	"github.com/ehutchllew/autoarmy/constants"
)

type TileMapObjectPropsJson struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type TileMapObjectsJson struct {
	Height     int                      `json:"height"`
	Gid        constants.ID             `json:"gid"`
	Name       string                   `json:"name"`
	Properties []TileMapObjectPropsJson `json:"properties"`
	Type       string                   `json:"type"`
	Width      int                      `json:"width"`
	X          float64                  `json:"x"`
	Y          float64                  `json:"y"`
}

type TileMapLayerJson struct {
	Data    []int                `json:"data,omitempty"`
	Height  int                  `json:"height"`
	Name    string               `json:"name"`
	Objects []TileMapObjectsJson `json:"objects,omitempty"`
	Width   int                  `json:"width"`
}

type TileMapTilesetJson struct {
	Firstgid uint8  `json:"firstgid"`
	Source   string `json:"source"`
}

type TileMapJson struct {
	Layers   []TileMapLayerJson   `json:"layers"`
	Tilesets []TileMapTilesetJson `json:"tilesets"`
}

func (t *TileMapJson) GenTilesets() ([]Tileset, error) {
	ts := make([]Tileset, 0)
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
