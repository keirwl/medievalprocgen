package config

import (
	"encoding/json"
	"os"
)

const configFilename = "data/config.json"

type Cam struct {
	Speed     float64
	MaxZoom   float64
	ZoomSpeed float64
	FPS       int
}

type TS struct {
	Filename string
	Width    float64
	Height   float64
	FullH    float64 `json:"FullHeight"`
	Size     float64
}

var (
	Camera  Cam
	Tileset TS
	All     struct {
		Camera  Cam
		Tileset TS
	}
)

func init() {
	f, err := os.Open(configFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&All)

	Camera = All.Camera
	Tileset = All.Tileset
}
