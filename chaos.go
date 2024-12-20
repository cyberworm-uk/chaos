package chaos

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"math/rand"
)

var transparent, blue, green, cyan, red, magenta, yellow, white, black = color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 255, 255, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 255, 255}, color.RGBA{255, 255, 0, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255}

var myPalette color.Palette = color.Palette{
	transparent,
	red,
	green,
	blue,
	cyan,
	magenta,
	yellow,
	white,
	black,
}

func myPaletted(r image.Rectangle) *image.Paletted {
	img := image.NewPaletted(r, myPalette)
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			img.Set(x, y, myPalette[0])
		}
	}
	return img
}

func myRGBA(r image.Rectangle) *image.RGBA {
	img := image.NewRGBA(r)
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			img.Set(x, y, myPalette[0])
		}
	}
	return img
}

type Chaos struct {
	fuzz       float64
	roots      []complex128
	proportion complex128
	point      complex128
}

func NewChaos(n int, prop float64, fuzz float64) (*Chaos, error) {
	var roots, e = unity(n, fuzz)
	if e != nil {
		return &Chaos{}, e
	}

	var c = &Chaos{
		fuzz:       fuzz,
		roots:      roots,
		proportion: proportion(n, prop),
		point:      complex(0, 0),
	}
	c.fuzzer()
	return c, nil
}

func (c *Chaos) fuzzer() {
	shift := func() complex128 {
		x := rand.Float64() * c.fuzz
		if rand.Float64() >= 0.5 {
			x *= -1
		}
		y := rand.Float64() * c.fuzz
		if rand.Float64() >= 0.5 {
			y *= -1
		}
		return complex(x, y)
	}
	r := rand.Float64() * c.fuzz
	if rand.Float64() >= 0.5 {
		r *= -1
	}
	var scale = complex(1+r, 0)
	for i := range c.roots {
		c.roots[i] *= scale
		c.roots[i] += shift()
	}
	c.point *= scale
	c.point += shift()
}

func (c *Chaos) Step() complex128 {
	var direction = c.roots[rand.Intn(len(c.roots))]
	var new = (c.point + direction) * c.proportion
	c.point = new
	return c.point
}

func (c *Chaos) String() string {
	return fmt.Sprintf("roots: %v, proportion: %v, fuzz: %v", c.roots, c.proportion, c.fuzz)
}

func unity(n int, fuzz float64) ([]complex128, error) {
	var out = []complex128{}
	if n < 3 {
		return out, errors.New("bad number of roots of unity, must be 3 or more")
	}
	for i := 0; i < n; i++ {
		var θ = math.Pi * 2 * float64(i) / float64(n)
		θ += rand.Float64() * math.Pi * 2 * fuzz
		var root = complex(math.Cos(θ), math.Sin(θ))
		out = append(out, root)
	}
	return out, nil
}

func proportion(n int, prop float64) complex128 {
	var pf float64
	if prop != 0 {
		return complex(prop, 0)
	}
	switch n % 4 {
	case 0:
		pf = 1 / (1 + math.Tan(math.Pi/float64(n)))
	case 1:
		pf = 1 / (1 + 2*math.Sin(math.Pi/(2*float64(n))))
	case 2:
		pf = 1 / (1 + 2*math.Sin(math.Pi/float64(n)))
	case 3:
		pf = 1 / (1 + 2*math.Sin(math.Pi/(2*float64(n))))
	}
	return complex(1-pf, 0)
}

func frame(width, steps int, last *image.Paletted, c *Chaos, shade color.Color, extend bool) *image.Paletted {
	var raw = myPaletted(last.Bounds())
	if extend {
		copy(raw.Pix, last.Pix)
	}
	for i := 0; i < steps; i++ {
		var pos = c.Step()
		var x, y int
		x = int((real(pos) + 1) * float64(width) / 2)
		y = int((imag(pos) + 1) * float64(width) / 2)
		raw.Set(x, y, shade)
	}
	return raw
}

func RevealChaos(width, n, frames, steps int, prop float64, fuzz float64) (*gif.GIF, error) {
	var g = &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
		Disposal:  []byte{},
	}
	var c, e = NewChaos(n, prop, fuzz)
	if e != nil {
		return nil, e
	}
	log.Printf("%s\n", c)
	var f = myPaletted(image.Rectangle{image.Point{0, 0}, image.Point{width, width}})
	for i := 1; i <= frames; i++ {
		f = frame(width, steps, f, c, color.Black, true)
		log.Printf("Frame: %v of %v\n", i, frames)
		g.Image = append(g.Image, f)
		g.Delay = append(g.Delay, 10)
		g.Disposal = append(g.Disposal, gif.DisposalBackground)
	}
	return g, nil
}

func ResolveChaos(width, n, steps int, props, prope, propi float64, fuzz float64) (*gif.GIF, error) {
	var frames int
	if prope > props {
		frames = int(math.Abs(prope-props) / math.Abs(propi))
	} else {
		frames = int(math.Abs(props-prope) / math.Abs(propi))
	}
	var g = &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
		Disposal:  []byte{},
	}
	var c, e = NewChaos(n, props, fuzz)
	if e != nil {
		return nil, e
	}
	log.Printf("%s\n", c)
	var f = myPaletted(image.Rectangle{image.Point{0, 0}, image.Point{width, width}})
	for i := 1; i <= frames; i++ {
		c.proportion += complex(propi, 0)
		f = frame(width, steps, f, c, color.Black, false)
		log.Printf("Frame: %v of %v\n", i, frames)
		g.Image = append(g.Image, f)
		g.Delay = append(g.Delay, 10)
		g.Disposal = append(g.Disposal, gif.DisposalBackground)
		c.fuzzer()
	}
	return g, nil
}

func glitchFrame(width, steps int, last *image.Paletted, c *Chaos, shade color.Color, extend bool) *image.Paletted {
	var raw = myPaletted(last.Bounds())
	if extend {
		copy(raw.Pix, last.Pix)
	}
	for i := 0; i < steps; i++ {
		var pos = c.Step()
		var x, y int
		x = int((real(pos) + 1) * float64(width) / 2)
		y = int((imag(pos) + 1) * float64(width) / 2)
		r0, g0, b0, _ := raw.At(x, y).RGBA()
		r1, g1, b1, _ := shade.RGBA()
		raw.Set(x, y, color.RGBA{uint8(r0 | r1), uint8(g0 | g1), uint8(b0 | b1), 255})
	}
	return raw
}

func GlitchChaos(width, n, frames, steps int, prop float64, fuzz float64) (*gif.GIF, error) {
	var g = &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
		Disposal:  []byte{},
	}
	var c, e = NewChaos(n, prop, fuzz)
	if e != nil {
		return nil, e
	}
	log.Printf("%s\n", c)
	var f = myPaletted(image.Rectangle{image.Point{0, 0}, image.Point{width, width}})
	for i := 1; i <= frames; i++ {
		f = glitchFrame(width, rand.Intn(steps), f, c, color.RGBA{255, 0, 0, 255}, false)
		c.fuzzer()
		f = glitchFrame(width, rand.Intn(steps), f, c, color.RGBA{0, 255, 0, 255}, true)
		c.fuzzer()
		f = glitchFrame(width, rand.Intn(steps), f, c, color.RGBA{0, 0, 255, 255}, true)
		c.fuzzer()
		log.Printf("Frame: %v of %v\n", i, frames)
		g.Image = append(g.Image, f)
		g.Delay = append(g.Delay, 10)
		g.Disposal = append(g.Disposal, gif.DisposalBackground)
	}
	return g, nil
}

func CrispChaos(width, n, steps int, prop float64, fuzz float64) (*image.RGBA, error) {
	var raw = myRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, width}})
	var c, e = NewChaos(n, prop, fuzz)
	if e != nil {
		return nil, e
	}
	for i := 0; i < steps; i++ {
		var pos = c.Step()
		var x, y int
		x = int((real(pos) + 1) * float64(width) / 2)
		y = int((imag(pos) + 1) * float64(width) / 2)
		r, g, b, a := raw.At(x, y).RGBA()
		if a != 0 {
			r = (r + 10) % 256
			g = (g + 10) % 256
			b = (b + 10) % 256
		}
		raw.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
	}
	return raw, nil
}
