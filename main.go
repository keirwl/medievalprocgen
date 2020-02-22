package main

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"./hex"
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

func init() {
	hex.Size = pixel.V(16, 16)
}

func run() {
	screenWidth, screenHeight := pixelgl.PrimaryMonitor().Size()
	cfg := pixelgl.WindowConfig{
		Title:     "Keir's Medieval Simulator",
		Bounds:    pixel.R(0, 0, screenWidth, screenHeight),
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	toolTip := rsfont.NewText(pixel.ZV, 22)

	tileMap := tilemap.MakeMap(15, 15)
	tileMap.Draw(tilemap.Batch, pixel.IM)

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

		mousePos := win.MousePosition()
		unProj := cam.Unproject(mousePos)
		hCoords := hex.FromPixel(unProj)
		col, row := hCoords.ToOffset()

		t := ""
		if col >= 0 && col < 15 && row >= 0 && row < 15 {
			t = tileMap.At(hCoords.Q, hCoords.R).Type.String()
		}

		toolTip.Clear()
		toolTip.Orig = mousePos
		fmt.Fprintf(toolTip, "(%.f, %.f)(%.f, %.f)\n(%+d, %+d)(%+d, %+d)\n%s",
			mousePos.X, mousePos.Y, unProj.X, unProj.Y, hCoords.Q, hCoords.R, col, row, t)

		win.Clear(colornames.Aliceblue)

		tilemap.Batch.Draw(win)

		win.SetMatrix(pixel.IM)
		toolTip.Draw(win, pixel.IM)

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
