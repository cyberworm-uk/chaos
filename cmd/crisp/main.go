package main

import (
	"chaos"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
)

func fname(roots, width int, prop, fuzz float64) string {
	return fmt.Sprintf("chaos-crisp-%v-%v-%v-%v.png", roots, prop, fuzz, width)
}

func main() {
	var w = flag.Int("width", 1080, "image width")
	var n = flag.Int("roots", 3, "n roots of unity (which n-gon)")
	var p = flag.Float64("prop", 0, "proportion of distance to move per step (0 for automatic)")
	var s = flag.Int("steps", 5000000, "number of steps per frame (dots drawn)")
	var fuzz = flag.Float64("fuzz", 0, "how fuzzy it should be")
	flag.Parse()
	g, e := chaos.CrispChaos(*w, *n, *s, *p, *fuzz)
	if e != nil {
		log.Fatal(e)
	}
	// open file
	f, e := os.Create(fname(*n, *w, *p, *fuzz))
	if e != nil {
		log.Fatal(e)
	}
	// write file
	e = png.Encode(f, g)
	if e != nil {
		log.Fatal(e)
	}
	// close file
	f.Close()
}
