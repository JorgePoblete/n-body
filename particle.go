package main

import (
	"crypto/rand"
	"image/color"
	"math/big"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Particle struct {
	X     float64
	Y     float64
	VX    float64
	VY    float64
	Mass  float64
	Ratio float64
	Color color.RGBA
}

func (p *Particle) Layout(gtx layout.Context) layout.Dimensions {
	defer op.Push(gtx.Ops).Pop()
	op.Offset(f32.Point{X: float32(p.X), Y: float32(p.Y)}).Add(gtx.Ops)

	paint.ColorOp{Color: p.Color}.Add(gtx.Ops)
	clip.RRect{
		Rect: f32.Rectangle{Max: f32.Point{X: float32(p.Ratio), Y: float32(p.Ratio)}},
		NE:   float32(p.Ratio) * .5, NW: float32(p.Ratio) * .5, SE: float32(p.Ratio) * .5, SW: float32(p.Ratio) * .5,
	}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(p.Ratio), Y: float32(p.Ratio)}}}.Add(gtx.Ops)
	op.InvalidateOp{}.Add(gtx.Ops)

	return layout.Dimensions{}
}

type Particles []Particle

const (
	PARTICLE_SIZE float64 = 3
)

func (p *Particles) randFloat(max float64) float64 {
	random, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return float64(random.Int64())
}

func (p *Particles) NewRandomParticle(maxX, maxY float64) Particle {
	return Particle{
		X:     p.randFloat(maxX),
		Y:     p.randFloat(maxY),
		VX:    0.0,
		VY:    0.0,
		Mass:  1.0,
		Ratio: PARTICLE_SIZE,
		Color: color.RGBA{
			A: 255,
			R: uint8(p.randFloat(255)),
			G: uint8(p.randFloat(255)),
			B: uint8(p.randFloat(255)),
		},
	}
}

func (p *Particles) NewParticle(x, y float64) Particle {
	return Particle{
		X:     x,
		Y:     y,
		VX:    0.0,
		VY:    0.0,
		Mass:  1.0,
		Ratio: PARTICLE_SIZE,
		Color: color.RGBA{
			A: 255,
			R: uint8(p.randFloat(255)),
			G: uint8(p.randFloat(255)),
			B: uint8(p.randFloat(255)),
		},
	}
}

func (p *Particles) AddParticle(x, y float64) {
	*p = append(*p, p.NewParticle(x, y))
}

func (p *Particles) AddRandomParticles(maxX, maxY float64, numberOfParticles int) {
	for i := 0; i < numberOfParticles; i++ {
		*p = append(*p, p.NewRandomParticle(maxX, maxY))
	}
}

func (p *Particles) Clear() {
	*p = Particles{}
}

func (p *Particles) Layout(gtx layout.Context) layout.Dimensions {
	for _, particle := range *p {
		particle.Layout(gtx)
	}
	return layout.Dimensions{}
}
