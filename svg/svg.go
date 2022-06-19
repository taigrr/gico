package svg

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"

	svg "github.com/ajstarks/svgo"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/taigrr/gitgraph/common"
	sc "github.com/taigrr/go-colorpallettes/simplecolor"
)

var colorsLoaded sync.Once
var colorScheme []sc.SimpleColor

func GetWeekSVG(frequencies []int) bytes.Buffer {
	squareColors := []sc.SimpleColor{}
	min, max := common.MinMax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, common.ColorForFrequency(f, min, max))
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
func GetYearSVG(frequencies []int) bytes.Buffer {
	squareColors := []sc.SimpleColor{}
	min, max := common.MinMax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, common.ColorForFrequency(f, min, max))
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
