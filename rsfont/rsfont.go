package rsfont

import (
	"io/ioutil"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
)

const (
	fontFilename = "assets/runescape_uf.ttf"
)

var face *truetype.Font
var atlases map[float64]*text.Atlas

func init() {
	var err error
	face, err = loadTTF(fontFilename)
	if err != nil {
		panic(err)
	}

	atlases = make(map[float64]*text.Atlas)
}

func NewText(v pixel.Vec, size float64) *text.Text {
	var atlas *text.Atlas
	atlas, ok := atlases[size]
	if !ok {
		atlas = text.NewAtlas(truetype.NewFace(face, &truetype.Options{
			Size:              size,
			GlyphCacheEntries: 1,
		}), text.ASCII)
		atlases[size] = atlas
	}

	t := text.New(v, atlas)
	t.Color = colornames.Black
	return t
}

func loadTTF(path string) (*truetype.Font, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return truetype.Parse(bytes)
}
