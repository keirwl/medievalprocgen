package main

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"

	"./rsfont"
	"./tilemap"
)

const (
	camSpeed     = 500.0
	minZoom      = 1.0
	maxZoom      = 10.0
	camZoomSpeed = 1.2
)

var (
	camPos  = pixel.ZV
	camZoom = minZoom
)

func run() {
	screenWidth, screenHeight := pixelgl.PrimaryMonitor().Size()
	cfg := pixelgl.WindowConfig{
		Title:     "Keir's Medieval Simulator",
		Bounds:    pixel.R(0, 0, screenWidth, screenHeight),
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	testText := rsfont.NewText(pixel.ZV)
	for _, r := range text.ASCII {
		testText.WriteRune(r)
	}

	var (
		frames = 0
		second = time.Tick(time.Second)
		last   = time.Now()
	)

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := zoomAndScroll(win, dt)
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			tile := tilemap.Random()
			mouse := cam.Unproject(win.MousePosition())
			tile.Draw(tilemap.Batch, pixel.IM.Moved(mouse))
		}

		win.Clear(colornames.Aliceblue)
		tilemap.Batch.Draw(win)
		testText.Draw(win, pixel.IM)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func zoomAndScroll(win *pixelgl.Window, dt float64) pixel.Matrix {
	camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
	if camZoom < minZoom {
		camZoom = minZoom
	}
	if camZoom > maxZoom {
		camZoom = maxZoom
	}

	switch {
	case win.Pressed(pixelgl.KeyLeft):
		camPos.X -= camSpeed * dt
	case win.Pressed(pixelgl.KeyRight):
		camPos.X += camSpeed * dt
	case win.Pressed(pixelgl.KeyDown):
		camPos.Y -= camSpeed * dt
	case win.Pressed(pixelgl.KeyUp):
		camPos.Y += camSpeed * dt
	}

	return pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
}

func main() {
	pixelgl.Run(run)
}
