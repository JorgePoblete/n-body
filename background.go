package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

type Background struct {
	Color color.RGBA
}

func (b *Background) Layout(gtx layout.Context) layout.Dimensions {
	defer op.Push(gtx.Ops).Pop()
	cs := gtx.Constraints
	paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
	paint.PaintOp{
		Rect: f32.Rectangle{Max: f32.Point{X: float32(cs.Max.X), Y: float32(cs.Max.Y)}},
	}.Add(gtx.Ops)
	op.InvalidateOp{}.Add(gtx.Ops)

	return layout.Dimensions{}
}
