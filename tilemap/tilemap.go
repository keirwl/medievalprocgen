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

func RandTile(q, r int) Tile {
	return New(q, r, TileType(rand.Intn(int(Max-1))+1))
}

func (t Tile) Draw(target pixel.Target, matrix pixel.Matrix) {
	t.Sprite.Draw(target, matrix.Moved(t.Rect.Min))
}

type Map struct {
	tiles  [][]Tile
	height int
	width  int
}

func MakeMap(height, width int) *Map {
	m := &Map{height: height, width: width}
	m.tiles = make([][]Tile, height)

	for row := 0; row < height; row++ {
		m.tiles[row] = make([]Tile, width)

		for col := 0; col < width; col++ {
			q, r := hex.ToAxial(col, row)
			m.tiles[row][col] = RandTile(q, r)
		}
	}

	return m
}

func (m *Map) At(q, r int) Tile {
	col, row := hex.ToOffset(q, r)
	return m.tiles[row][col]
}

func (m *Map) SetTileType(q, r int, t TileType) {
	col, row := hex.ToOffset(q, r)
	m.tiles[row][col].Type = t
	m.tiles[row][col].Sprite = sprites[t]
}

func (m *Map) Draw(t pixel.Target, matrix pixel.Matrix) {
	for row := m.height - 1; row >= 0; row-- {
		for col := 0; col < m.width; col += 2 {
			m.tiles[row][col].Draw(t, matrix)
		}
		for col := 1; col < m.width; col += 2 {
			m.tiles[row][col].Draw(t, matrix)
		}
	}
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
