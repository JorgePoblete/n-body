package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type controls struct {
	start    widget.Clickable
	stop     widget.Clickable
	generate widget.Clickable
	clear    widget.Clickable
	drawTree widget.Bool
	DT       widget.Float
	run      bool
	running  bool
}

func (c *controls) Layout(gtx layout.Context) layout.Dimensions {
	for c.generate.Clicked() {
		log.Printf("generating")
		cs := gtx.Constraints
		X := float64(cs.Max.X)
		Y := float64(cs.Max.Y)
		particles.AddRandomParticles(X, Y, NumberOfParticles)
	}

	if len(particles) > 0 {
		for c.clear.Clicked() {
			log.Printf("clearing")
			c.run = false
			for c.running {
				time.Sleep(1 * time.Millisecond)
			}
			particles.Clear()
			quadtree.Clear()
		}
	}
	if !c.running {
		for c.start.Clicked() {
			c.run = true
			go func() {
				log.Printf("started")
				c.running = true
				for c.run {
					quadtree.Compute(particles, G, DT*math.Pow(10.0, float64(c.DT.Value)))
				}
				c.running = false
				log.Printf("stoped")
			}()
		}
	}
	if c.running {
		for c.stop.Clicked() {
			c.run = false
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(theme, &c.generate, "Generate").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(theme, &c.clear, "Clear").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(theme, &c.start, "Start").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(theme, &c.stop, "Stop").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Switch(theme, &c.drawTree).Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			widget := material.H6(theme, fmt.Sprintf("Speed: %d", int(c.DT.Value)))
			widget.Color = WHITE
			return widget.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(theme, &c.DT, 0.0, 10.0).Layout(gtx)
		}),
	)
}
