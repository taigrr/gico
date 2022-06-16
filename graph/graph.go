package graph

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"

	svg "github.com/ajstarks/svgo"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	sc "github.com/taigrr/go-colorpallettes/simplecolor"
)

var colorsLoaded sync.Once
var colorScheme []sc.SimpleColor

func CreateGraph() bytes.Buffer {
	var x bytes.Buffer
	return x
}

func init() {
	colors := []string{"#767960", "#a7297f", "#e8ca89", "#f5efd6", "#858966"}
	for _, c := range colors {
		color := sc.FromHexString(c)
		_ = c
		colorScheme = append(colorScheme, color)
	}

}

func SetColorScheme(c []color.Color) {
	for _, c := range c {
		colorScheme = append(colorScheme, sc.FromRGBA(c.RGBA()))
	}
}

func ColorForFrequency(freq, min, max int) sc.SimpleColor {
	spread := max - min
	interval := float64(spread) / float64(len(colorScheme))
	colorIndex := 0
	for i := float64(min); i < float64(freq); i += float64(interval) {
		colorIndex++
	}
	if colorIndex > len(colorScheme)-1 {
		colorIndex = len(colorScheme) - 1
	}
	return colorScheme[colorIndex]
}
func GetImage(frequencies []int) bytes.Buffer {
	squareColors := []sc.SimpleColor{}
	min, max := minmax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, ColorForFrequency(f, min, max))
	}
	return drawImage(squareColors)
}

func svgToPng() {
	w, h := 512, 512

	in, err := os.Open("in.svg")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	icon, _ := oksvg.ReadIconStream(in)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	out, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, rgba)
	if err != nil {
		panic(err)
	}
}

func drawImage(c []sc.SimpleColor) bytes.Buffer {
	//TODO here, draw suqares in appropriate colors, hopefully as an svg
	var sb bytes.Buffer
	sbw := bufio.NewWriter(&sb)
	width := 717
	height := 112
	squareLength := 10
	canvas := svg.New(sbw)
	canvas.Start(width, height)
	for i, c := range c {
		canvas.Square(10+squareLength+width/52*(i/7), squareLength/2+height/7*(i%7), squareLength, fmt.Sprintf("fill:%s", c.HexString()))
	}
	canvas.Text(width/100, squareLength*2+10, "Mon", "text-anchor:middle;font-size:10px;fill:black")
	canvas.End()
	sbw.Flush()
	return sb
}

func minmax(f []int) (int, int) {
	min, max := 0, 0
	for _, x := range f {
		if x < min {
			min = x
		} else if x > max {
			max = x
		}
	}
	return min, max
}
