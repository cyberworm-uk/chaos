package main

import (
	"flag"
	"fmt"
	"image/gif"
	"log"
	"os"

	"github.com/cyberworm-uk/chaos"
)

func fname(roots, width int, prop, fuzz float64) string {
	return fmt.Sprintf("chaos-glitch-%v-%v-%v-%v.gif", roots, prop, fuzz, width)
}

func main() {
	var w = flag.Int("width", 1080, "image width")
	var n = flag.Int("roots", 3, "n roots of unity")
	var p = flag.Float64("prop", 0, "proportion of distance to move per step (0 for automatic)")
	var s = flag.Int("steps", 100000, "number of steps per frame (dots drawn)")
	var frames = flag.Int("frames", 100, "number of frames")
	var fuzz = flag.Float64("fuzz", 0.05, "how fuzzy it should be")
	flag.Parse()
	g, e := chaos.GlitchChaos(*w, *n, *frames, *s, *p, *fuzz)
	if e != nil {
		log.Fatal(e)
	}
	// open file
	f, e := os.Create(fname(*n, *w, *p, *fuzz))
	if e != nil {
		log.Fatal(e)
	}
	// write file
	e = gif.EncodeAll(f, g)
	if e != nil {
		log.Fatal(e)
	}
	// close file
	f.Close()
}
