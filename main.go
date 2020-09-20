package main

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var WHITE = color.RGBA{R: 255, G: 255, B: 255, A: 255}
var GREEN = color.RGBA{R: 0, G: 255, B: 0, A: 255}
var LIGHT_GREEN = color.RGBA{R: 173, G: 255, B: 47, A: 255}

var G = 6.67408e-11
var DT = 10.0
var NumberOfParticles = 1000

var theme *material.Theme
var particles Particles
var quadtree QuadTree
var programControls controls
var programBackground = Background{Color: color.RGBA{A: 255}}

func main() {
	go func() {
		w := app.NewWindow()
		theme = material.NewTheme(gofont.Collection())
		if err := loop(w); err != nil {
			return
		}
	}()

	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			programBackground.Layout(gtx)
			layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						dims := programControls.Layout(gtx)
						return dims
					}),
					layout.Flexed(0.98, func(gtx layout.Context) layout.Dimensions {
						if programControls.drawTree.Value {
							quadtree.Layout(gtx)
						}
						return particles.Layout(gtx)
					}),
				)
			})
			e.Frame(gtx.Ops)
		}
	}
}
