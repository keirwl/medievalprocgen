package hex

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
)

const sqrt3 = 1.7320508075688772

var (
	Size   pixel.Vec
	Origin pixel.Vec

	dirHexes = []*Hex{
		New(1, 0),
		New(1, -1),
		New(0, -1),
		New(-1, 0),
		New(-1, 1),
		New(0, 1),
	}

	toPixel = [4]float64{
		3.0 / 2.0,
		0.0,
		sqrt3 / 2.0,
		sqrt3,
	}

	toHex = [4]float64{
		2.0 / 3.0,
		0.0,
		-1.0 / 3.0,
		sqrt3 / 3.0,
	}
)

type Hex struct {
	Q, R    int
	corners [6]pixel.Vec
	centre  pixel.Vec
}

func (h *Hex) String() string {
	return fmt.Sprintf("Hex(%+d, %+d)", h.Q, h.R)
}

func New(q, r int) *Hex {
	return &Hex{Q: q, R: r}
}

func FromPixel(p pixel.Vec) *Hex {
	pt := pixel.Vec{
		X: (p.X - Origin.X) / Size.X,
		Y: (p.Y - Origin.Y) / Size.Y,
	}

	q := toHex[0]*pt.X + toHex[1]*pt.Y
	r := toHex[2]*pt.X + toHex[3]*pt.Y
	qi, ri := hex_round(q, r)
	return &Hex{Q: qi, R: ri}
}

func hex_round(q, r float64) (int, int) {
	s := -q - r
	qi := int(math.Round(q))
	ri := int(math.Round(r))
	si := int(math.Round(s))

	qDiff := math.Abs(float64(qi) - q)
	rDiff := math.Abs(float64(ri) - r)
	sDiff := math.Abs(float64(si) - s)

	if qDiff > rDiff && qDiff > sDiff {
		qi = -ri - si
	} else if rDiff > sDiff {
		ri = -qi - si
	}

	return qi, ri
}

func (h *Hex) ToPixel() pixel.Vec {
	if h.centre != pixel.ZV {
		return h.centre
	}

	centre := pixel.Vec{
		X: Origin.X + (Size.X * (toPixel[0]*float64(h.Q) + toPixel[1]*float64(h.R))),
		Y: Origin.Y + (Size.Y * (toPixel[2]*float64(h.Q) + toPixel[3]*float64(h.R))),
	}
	h.centre = centre
	return centre
}

func (h *Hex) Corners() [6]pixel.Vec {
	corners := [6]pixel.Vec{}
	if h.corners != corners {
		return h.corners
	}

	centre := h.ToPixel()

	for i := 0; i < 6; i++ {
		offset := hex_corner_offset(i)
		corners[i] = pixel.V(centre.X+offset.X, centre.Y+offset.Y)
	}
	h.corners = corners
	return corners
}

func hex_corner_offset(corner int) pixel.Vec {
	angle := 2.0 * math.Pi * float64(corner) / 6.0
	return pixel.V(Size.X*math.Cos(angle), Size.Y*math.Sin(angle))
}

func Grid(height, width int) map[Hex]struct{} {
	grid := make(map[Hex]struct{}, height*width)
	for q := 0; q < width; q++ {
		qOffset := (q + 1) / 2
		for r := -qOffset; r < height-qOffset; r++ {
			grid[Hex{Q: q, R: r}] = struct{}{}
		}
	}

	return grid
}

func (this *Hex) Equals(other *Hex) bool {
	return this.Q == other.Q &&
		this.R == other.R
}

func Equals(a, b *Hex) bool {
	return a.Q == b.Q &&
		a.R == b.R
}

func (this *Hex) Add(other *Hex) {
	this.Q += other.Q
	this.Q += other.R
}

func Add(a, b *Hex) *Hex {
	return &Hex{
		Q: a.Q + b.Q,
		R: a.R + b.R,
	}
}

func (this *Hex) Sub(other *Hex) {
	this.Q -= other.Q
	this.R -= other.R
}

func Sub(a, b *Hex) *Hex {
	return &Hex{
		Q: a.Q - b.Q,
		R: a.R - b.R,
	}
}

func (this *Hex) Mult(other *Hex) {
	this.Q *= other.Q
	this.R *= other.R
}

func Mult(a, b *Hex) *Hex {
	return &Hex{
		Q: a.Q * b.Q,
		R: a.R * b.R,
	}
}

func (h *Hex) Length() int {
	return int(
		(float64(abs(h.Q)) + float64(abs(h.R)) +
			float64(abs(-h.Q-h.R))) / 2.)
}

func Distance(a, b *Hex) int {
	return Sub(a, b).Length()
}

func Direction(dir int) *Hex {
	return dirHexes[(6+(dir%6))%6]
}

func (h *Hex) Neighbour(dir int) *Hex {
	return Add(h, Direction(dir))
}

func abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}
