// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"fmt"
	"image"

	rw "github.com/mattn/go-runewidth"

	. "github.com/gizak/termui"
)

type StackedBarChart struct {
	Block
	BarAttrs   []Attribute
	LabelAttrs []Attribute
	NumAttrs   []Attribute
	NumFmt     func(int) string
	Data       [][]int
	Labels     []string
	BarWidth   int
	BarGap     int
	MaxVal     int
}

func NewStackedBarChart() *StackedBarChart {
	return &StackedBarChart{
		Block:      *NewBlock(),
		BarAttrs:   Theme.StackedBarChart.Bars,
		LabelAttrs: Theme.StackedBarChart.Labels,
		NumAttrs:   Theme.StackedBarChart.Nums,
		NumFmt:     func(n int) string { return fmt.Sprint(n) },
		BarGap:     1,
		BarWidth:   3,
	}
}

func (bc *StackedBarChart) Draw(buf *Buffer) {
	bc.Block.Draw(buf)

	maxVal := bc.MaxVal
	if maxVal == 0 {
		for _, data := range bc.Data {
			maxVal = MaxInt(maxVal, SumIntSlice(data))
		}
	}

	barXCoordinate := bc.Inner.Min.X

	for i, bar := range bc.Data {
		// draw stacked bars
		stackedBarYCoordinate := 0
		for j, data := range bar {
			// draw each stacked bar
			height := int((float64(data) / float64(maxVal)) * float64(bc.Inner.Dy()-1))
			for x := barXCoordinate; x < MinInt(barXCoordinate+bc.BarWidth, bc.Inner.Max.X); x++ {
				for y := (bc.Inner.Max.Y - 2) - stackedBarYCoordinate; y > (bc.Inner.Max.Y-2)-stackedBarYCoordinate-height; y-- {
					c := Cell{' ', AttrPair{ColorDefault, SelectAttr(bc.BarAttrs, j)}}
					buf.SetCell(c, image.Pt(x, y))
				}
			}

			// draw number
			numberXCoordinate := barXCoordinate + int((float64(bc.BarWidth) / 2)) - 1
			buf.SetString(
				fmt.Sprintf("%d", data),
				AttrPair{
					SelectAttr(bc.NumAttrs, j+1),
					SelectAttr(bc.BarAttrs, j),
				},
				image.Pt(numberXCoordinate, (bc.Inner.Max.Y-2)-stackedBarYCoordinate),
			)

			stackedBarYCoordinate += height
		}

		// draw label
		labelXCoordinate := barXCoordinate + MaxInt(
			int((float64(bc.BarWidth)/2))-int((float64(rw.StringWidth(bc.Labels[i]))/2)),
			0,
		)
		buf.SetString(
			TrimString(bc.Labels[i], bc.BarWidth),
			AttrPair{SelectAttr(bc.LabelAttrs, i), ColorDefault},
			image.Pt(labelXCoordinate, bc.Inner.Max.Y-1),
		)

		barXCoordinate += (bc.BarWidth + bc.BarGap)
	}
}