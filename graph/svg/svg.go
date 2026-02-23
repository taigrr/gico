package svg

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	svg "github.com/ajstarks/svgo"
	sc "github.com/taigrr/simplecolorpalettes/simplecolor"

	"github.com/taigrr/gico/graph/common"
)

func GetWeekSVG(frequencies []int, shouldHighlight bool) bytes.Buffer {
	squareColors := []sc.SimpleColor{}
	min, max := common.MinMax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, common.ColorForFrequency(f, min, max))
	}
	return drawWeekImage(squareColors, frequencies, shouldHighlight)
}

func drawWeekImage(c []sc.SimpleColor, freq []int, shouldHighlight bool) bytes.Buffer {
	var sb bytes.Buffer
	sbw := bufio.NewWriter(&sb)
	squareLength := 10
	width := len(c)*squareLength*2 + squareLength
	height := squareLength * 2
	canvas := svg.New(sbw)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:black")
	for i, s := range c {
		if shouldHighlight && i == len(c)-1 {
			if freq[i] == 0 {
				s = sc.FromHexString("#FF0000")
			}
		}
		canvas.Square(squareLength*2*(i)+squareLength, squareLength/2, squareLength, fmt.Sprintf("fill:%s; value:%d", s.ToHex(), freq[i]))
	}
	canvas.End()
	sbw.Flush()
	return sb
}

func GetYearSVG(frequencies []int, shouldHighlight bool) bytes.Buffer {
	squareColors := []sc.SimpleColor{}
	min, max := common.MinMax(frequencies)
	for _, f := range frequencies {
		squareColors = append(squareColors, common.ColorForFrequency(f, min, max))
	}
	return drawYearImage(squareColors, frequencies, shouldHighlight)
}

func drawYearImage(c []sc.SimpleColor, freq []int, shouldHighlight bool) bytes.Buffer {
	var sb bytes.Buffer
	now := time.Now()
	sbw := bufio.NewWriter(&sb)
	squareLength := 10
	width := (len(c)/7+1)*squareLength*2 + squareLength*5
	height := squareLength*9 + squareLength*3
	canvas := svg.New(sbw)
	canvas.Start(width, height)
	for i, s := range c {
		if shouldHighlight && i == now.YearDay()-1 {
			if freq[i] == 0 {
				s = sc.FromHexString("#FF0000")
			}
		}
		canvas.Square(2*squareLength+width/(len(c)/7+1)*(i/7)+squareLength*2, squareLength/2+height/7*(i%7), squareLength, fmt.Sprintf("fill:%s; value:%d", s.ToHex(), freq[i]))
	}
	canvas.End()
	sbw.Flush()
	return sb
}
