package tilemap

import (
	"image"
	"math/rand"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"

	"../hex"
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
	sprites     = make([]*pixel.Sprite, Max)
	Batch       *pixel.Batch
)

type Tile struct {
	hex.Hex
	Type   TileType
	Sprite *pixel.Sprite
	Rect   pixel.Rect
}

func New(q int, r int, t TileType) Tile {
	tile := Tile{
		hex.Hex{Q: q, R: r},
		t,
		sprites[t],
		pixel.Rect{},
	}

	c := tile.ToPixel()
	tile.Rect = pixel.R(
		c.X-tileWidth/2,
		c.Y-tileHeight/2,
		c.X+tileWidth/2,
		c.Y+tileFullH-tileHeight/2,
	)

	return tile
}

func init() {
	var err error
	spritesheet, err = loadPicture(tilesetFilename)
	if err != nil {
		panic(err)
	}

	sprites[None] = pixel.NewSprite(spritesheet, pixel.ZR)
	var t TileType = Grass
outer:
	for y := spritesheet.Bounds().Max.Y; y > spritesheet.Bounds().Min.Y; y -= tileFullH {
		for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += tileWidth {

			sprites[t] = pixel.NewSprite(spritesheet, pixel.R(x, y-tileFullH, x+tileWidth, y))
			t++
			if t >= Max {
				break outer
			}
		}
	}

	Batch = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
}

func Sprite(t TileType) *pixel.Sprite {
	return sprites[t]
}

func Random() *pixel.Sprite {
	return Sprite(TileType(rand.Intn(int(Max-1)) + 1))
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
