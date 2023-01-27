package term

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/muesli/termenv"
	sc "github.com/taigrr/simplecolorpalettes/simplecolor"

	"github.com/taigrr/gitgraph/common"
)

var (
	colorsLoaded sync.Once
	colorScheme  []sc.SimpleColor
)

func GetWeekUnicode(frequencies []int) {
	squareColors := []sc.SimpleColor{}
	min, max := common.MinMax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, common.ColorForFrequency(f, min, max))
	}
	drawWeekUnicode(squareColors)
}

func drawWeekUnicode(c []sc.SimpleColor) {
	// o := termenv.NewOutput(os.Stdout)
	o := termenv.NewOutputWithProfile(os.Stdout, termenv.TrueColor)
	for w, color := range c {
		style := o.String(block).Foreground(termenv.TrueColor.Color(color.ToHex()))
		fmt.Print(style.String())
		//	termenv.SetForegroundColor(termenv.ForegroundColor())
		if w == len(c)-1 {
			fmt.Println()
		} else {
			fmt.Print(" ")
		}
	}
}

func GetYearUnicode(frequencies []int) string {
	squareColors := []sc.SimpleColor{}
	min, max := common.MinMax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, common.ColorForFrequency(f, min, max))
	}
	return drawYearUnicode(squareColors)
}

func drawYearUnicode(c []sc.SimpleColor) string {
	// o := termenv.NewOutput(os.Stdout)
	var s strings.Builder
	o := termenv.NewOutputWithProfile(os.Stdout, termenv.TrueColor)
	weeks := [7][]sc.SimpleColor{{}}
	for i := 0; i < 7; i++ {
		weeks[i] = []sc.SimpleColor{}
	}
	for i := range c {
		weeks[i%7] = append(weeks[i%7], c[i])
	}
	for _, row := range weeks {
		for w, d := range row {
			style := o.String(block).Foreground(termenv.TrueColor.Color(d.ToHex()))
			s.WriteString(style.String())
			if w == len(row)-1 {
				s.WriteString("\n")
			} else {
				s.WriteString(" ")
			}

		}
	}
	return s.String()
}
