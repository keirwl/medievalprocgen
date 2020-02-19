package main

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"./hex"
	"./rsfont"
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
	hex.Size = pixel.V(64, 64)
}

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

	testDraw := imdraw.New(nil)

	testDraw.Color = colornames.Black

	t := rsfont.NewText(pixel.ZV, 18)
	hexGrid := hex.Grid(5, 5)
	for h, _ := range hexGrid {
		corners := h.Corners()
		for _, c := range corners {
			testDraw.Push(c)
		}
		testDraw.Push(corners[0])
		testDraw.Line(1)

		t.Dot = h.ToPixel()
		fmt.Fprintf(t, "%+d, %+d", h.Q, h.R)
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

		mousePos := win.MousePosition()
		unProj := cam.Unproject(mousePos)
		hCoords := hex.FromPixel(unProj)

		toolTip.Clear()
		toolTip.Orig = mousePos
		fmt.Fprintf(toolTip, "(%.f, %.f)\n(%.f, %.f)\n(%+d, %+d)",
			mousePos.X, mousePos.Y, unProj.X, unProj.Y, hCoords.Q, hCoords.R)

		win.Clear(colornames.Aliceblue)

		testDraw.Draw(win)
		t.Draw(win, pixel.IM)

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
