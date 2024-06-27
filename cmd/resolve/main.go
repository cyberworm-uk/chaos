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
	return fmt.Sprintf("chaos-resolve-%v-%v-%v-%v.gif", roots, prop, fuzz, width)
}

func main() {
	var w = flag.Int("width", 1080, "image width")
	var n = flag.Int("roots", 3, "n roots of unity")
	var ps = flag.Float64("prop-start", 1.0, "start proportion value")
	var pe = flag.Float64("prop-end", 0.01, "ending proportion value")
	var pi = flag.Float64("prop-inc", -0.01, "proportion increment value")
	var s = flag.Int("steps", 10000, "number of steps per frame (dots drawn)")
	var fuzz = flag.Float64("fuzz", 0, "how fuzzy it should be")
	flag.Parse()
	g, e := chaos.ResolveChaos(*w, *n, *s, *ps, *pe, *pi, *fuzz)
	if e != nil {
		log.Fatal(e)
	}
	// open file
	f, e := os.Create(fname(*n, *w, *ps, *fuzz))
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
