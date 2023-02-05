package common

import (
	"bytes"
	"fmt"
	"image/color"
	"sync"

	sc "github.com/taigrr/simplecolorpalettes/simplecolor"
)

var (
	colorsLoaded sync.Once
	colorScheme  []sc.SimpleColor
)

func CreateGraph() bytes.Buffer {
	var x bytes.Buffer
	return x
}

func init() {
	colors := []string{"#767960", "#a7297f", "#e8ca89", "#f5efd6", "#158266"}
	colors = []string{"#0e4429", "#006d32", "#26a641", "#39d353"}
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
	fmt.Printf("f: %d, min: %d, max: %d\n", freq, min, max)
	if freq == 0 {
		return sc.SimpleColor(0)
	}
	spread := max - min
	if spread < len(colorScheme)-1 {
		return colorScheme[freq-min+1]
	}
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

func MinMax(f []int) (int, int) {
	max := 0
	for _, x := range f {
		if x > max {
			max = x
		}
	}
	min := max
	for _, x := range f {
		if x < min && x != 0 {
			min = x
		}
	}

	return min, max
}
