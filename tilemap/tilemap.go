package tilemap

import (
	"image"
	"math/rand"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
)

const (
	tilesetFilename = "assets/fantasyhextiles_v3.png"
	tileWidth       = 32
	tileHeight      = 24
	tileFullH       = 48
	tileSide        = 16
)

var (
	spritesheet pixel.Picture
	tiles       = make([]pixel.Rect, Max)
	initialised bool
)

func New() (batch *pixel.Batch, err error) {
	spritesheet, err = loadPicture(tilesetFilename)
	if err != nil {
		return nil, err
	}

	var t TileType
outer:
	for y := spritesheet.Bounds().Max.Y; y > spritesheet.Bounds().Min.Y; y -= tileFullH {
		for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += tileWidth {

			tiles[t] = pixel.R(x, y-tileFullH, x+tileWidth, y)
			t++
			if t >= Max {
				break outer
			}
		}
	}

	batch = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	initialised = true
	return batch, nil
}

func Sprite(t TileType) *pixel.Sprite {
	if !initialised {
		panic("tilemap spritesheet has not been initialised!")
	}
	return pixel.NewSprite(spritesheet, tiles[t])
}

func Random() *pixel.Sprite {
	return Sprite(TileType(rand.Intn(int(Max))))
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}
