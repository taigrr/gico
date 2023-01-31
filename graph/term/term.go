package term

import (
	"os"
	"strings"
	"sync"

	"github.com/muesli/termenv"
	sc "github.com/taigrr/simplecolorpalettes/simplecolor"

	"github.com/taigrr/gico/graph/common"
)

var (
	colorsLoaded sync.Once
	colorScheme  []sc.SimpleColor
)

func GetWeekUnicode(frequencies []int) string {
	squareColors := []sc.SimpleColor{}
	min, max := common.MinMax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, common.ColorForFrequency(f, min, max))
	}
	return drawWeekUnicode(squareColors)
}

func drawWeekUnicode(c []sc.SimpleColor) string {
	// o := termenv.NewOutput(os.Stdout)
	s := strings.Builder{}
	o := termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	for w, color := range c {
		style := o.String(block).Foreground(termenv.TrueColor.Color(color.ToHex()))
		s.WriteString(style.String())
		//	termenv.SetForegroundColor(termenv.ForegroundColor())
		if w == len(c)-1 {
			s.WriteString("\n")
		} else {
			s.WriteString(" ")
		}
	}
	return s.String()
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
	o := termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	weekRows := [7][]sc.SimpleColor{{}}
	for i := 0; i < 7; i++ {
		weekRows[i] = []sc.SimpleColor{}
	}
	for i := 0; i < len(c); i++ {
		weekRows[i%7] = append(weekRows[i%7], c[i])
	}
	for _, row := range weekRows {
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
