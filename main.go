package main

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"./config"
	"./hex"
	"./rsfont"
	"./tilemap"
)

const (
	minZoom = 1.0
)

var (
	camPos     = pixel.ZV
	camZoom    = minZoom
	msPerFrame = time.Millisecond * time.Duration(1000/config.Camera.FPS)
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

	toolTip := rsfont.NewText(pixel.ZV, 22)

	tileMap := tilemap.MakeMap(15, 15)
	tileMap.Draw(tilemap.Batch, pixel.IM)

	var (
		frames = 0
		second = time.Tick(time.Second)
		last   = time.Now()
	)

	for !win.Closed() {
		dt := time.Since(last)
		time.Sleep(msPerFrame - dt)
		last = time.Now()

		cam := zoomAndScroll(win, dt.Seconds())
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
	camZoom *= math.Pow(config.Camera.ZoomSpeed, win.MouseScroll().Y)
	camZoom = pixel.Clamp(camZoom, minZoom, config.Camera.MaxZoom)

	switch {
	case win.Pressed(pixelgl.KeyLeft):
		camPos.X -= config.Camera.Speed * dt / camZoom
	case win.Pressed(pixelgl.KeyRight):
		camPos.X += config.Camera.Speed * dt / camZoom
	case win.Pressed(pixelgl.KeyDown):
		camPos.Y -= config.Camera.Speed * dt / camZoom
	case win.Pressed(pixelgl.KeyUp):
		camPos.Y += config.Camera.Speed * dt / camZoom
	}

	return pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
}

func main() {
	pixelgl.Run(run)
}
