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
	// colors = []string{"#000000", "#0e4429", "#006d32", "#26a641", "#39d353"}
	for _, c := range colors {
		color := sc.FromHexString(c)
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
func GetWeekImage(frequencies []int) bytes.Buffer {
	squareColors := []sc.SimpleColor{}
	min, max := minmax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, ColorForFrequency(f, min, max))
	}
	return drawWeekImage(squareColors)
}

func drawWeekImage(c []sc.SimpleColor) bytes.Buffer {
	var sb bytes.Buffer
	sbw := bufio.NewWriter(&sb)
	squareLength := 10
	width := (len(c) + 1) * squareLength * 2
	height := squareLength * 2
	canvas := svg.New(sbw)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:black")
	for i, s := range c {
		canvas.Square(squareLength*2*(i+1), squareLength/2, squareLength, fmt.Sprintf("fill:%s", s.HexString()))
	}
	canvas.End()
	sbw.Flush()
	return sb
}
func GetYearImage(frequencies []int) bytes.Buffer {
	squareColors := []sc.SimpleColor{}
	min, max := minmax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, ColorForFrequency(f, min, max))
	}
	return drawYearImage(squareColors)
}

func drawYearImage(c []sc.SimpleColor) bytes.Buffer {
	//TODO here, draw suqares in appropriate colors, hopefully as an svg
	var sb bytes.Buffer
	sbw := bufio.NewWriter(&sb)
	squareLength := 10
	width := (len(c)/7+1)*squareLength*2 + squareLength*5
	height := squareLength*9 + squareLength*3
	canvas := svg.New(sbw)
	canvas.Start(width, height)
	for i, s := range c {
		canvas.Square(2*squareLength+width/(len(c)/7+1)*(i/7)+squareLength*2, squareLength/2+height/7*(i%7), squareLength, fmt.Sprintf("fill:%s", s.HexString()))
	}
	canvas.Text(2*squareLength, squareLength*3, "Mon", fmt.Sprintf("text-anchor:middle;font-size:%dpx;fill:black", squareLength))
	canvas.Text(2*squareLength, int(float64(squareLength)*6.5), "Wed", fmt.Sprintf("text-anchor:middle;font-size:%dpx;fill:black", squareLength))
	canvas.Text(2*squareLength, int(float64(squareLength))*10, "Fri", fmt.Sprintf("text-anchor:middle;font-size:%dpx;fill:black", squareLength))
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
