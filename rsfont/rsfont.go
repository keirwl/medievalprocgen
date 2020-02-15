package rsfont

import (
	"io/ioutil"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

const (
	fontFilename = "assets/runescape_uf.ttf"
	fontSize     = 24
)

var atlas *text.Atlas

func init() {
	face, err := loadTTF(fontFilename, fontSize)
	if err != nil {
		panic(err)
	}

	atlas = text.NewAtlas(face, text.ASCII)
}

func NewText(v pixel.Vec) *text.Text {
	t := text.New(v, atlas)
	t.Color = colornames.Black
	return t
}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
