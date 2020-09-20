package main

import (
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

type QuadTree struct {
	Quadrant *Quadrant
}

func (qt *QuadTree) Compute(particles Particles, G, DT float64) {
	Xmin := -1.0
	Xmax := -1.0
	Ymin := -1.0
	Ymax := -1.0

	for _, particle := range particles {
		if Xmin == -1.0 || Xmin > particle.X {
			Xmin = particle.X
		}
		if Xmax == -1.0 || Xmax < particle.X {
			Xmax = particle.X
		}
		if Ymin == -1.0 || Ymin > particle.Y {
			Ymin = particle.Y
		}
		if Ymax == -1.0 || Ymax < particle.Y {
			Ymax = particle.Y
		}
	}
	qt.Quadrant = &Quadrant{
		RemainingQuadrants: 1000,
		Xmin:               Xmin,
		Xmax:               Xmax,
		Ymin:               Ymin,
		Ymax:               Ymax,
	}
	qt.Quadrant = qt.Quadrant.SubQuadrant(particles)

	for i, particle := range particles {
		Fx, Fy := qt.Quadrant.ForceOverParticle(particle, G)
		particles[i].VX += DT * Fx
		particles[i].VY += DT * Fy
		particles[i].X += particles[i].VX
		particles[i].Y += particles[i].VY
	}
}

func (qt *QuadTree) Clear() {
	qt.Quadrant = nil
}

func (qt *QuadTree) Layout(gtx layout.Context) layout.Dimensions {
	if qt.Quadrant != nil {
		qt.Quadrant.Layout(gtx)
	}
	return layout.Dimensions{}
}

type Quadrant struct {
	Quadrants          []*Quadrant
	Particles          Particles
	RemainingQuadrants int
	Mass               float64
	Xmin               float64
	Xmax               float64
	Ymin               float64
	Ymax               float64
	CenterOfMassX      float64
	CenterOfMassY      float64
}

func (q *Quadrant) Layout(gtx layout.Context) layout.Dimensions {
	defer op.Push(gtx.Ops).Pop()
	var lineSize float32 = 0.001
	paint.ColorOp{Color: LIGHT_GREEN}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(float32(q.Xmin), float32(q.Ymin), float32(q.Xmax), float32(q.Ymin)+lineSize)}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(float32(q.Xmin), float32(q.Ymax), float32(q.Xmax), float32(q.Ymax)-lineSize)}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(float32(q.Xmin), float32(q.Ymin), float32(q.Xmin)+lineSize, float32(q.Ymax))}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(float32(q.Xmax)-lineSize, float32(q.Ymin), float32(q.Xmax), float32(q.Ymax))}.Add(gtx.Ops)
	op.InvalidateOp{}.Add(gtx.Ops)

	for _, quadrant := range q.Quadrants {
		quadrant.Layout(gtx)
	}
	return layout.Dimensions{}
}

func (q *Quadrant) NewQuadrant(x1, x2, y1, y2 float64, remaining int) *Quadrant {
	return &Quadrant{
		RemainingQuadrants: remaining,
		Xmin:               x1,
		Xmax:               x2,
		Ymin:               y1,
		Ymax:               y2,
	}
}

func (q *Quadrant) SubQuadrant(particles Particles) *Quadrant {
	for _, particle := range particles {
		if particle.X >= q.Xmin && particle.X <= q.Xmax && particle.Y >= q.Ymin && particle.Y <= q.Ymax {
			q.Particles = append(q.Particles, particle)
			q.Mass += particle.Mass
			q.CenterOfMassX += particle.Mass * particle.X
			q.CenterOfMassY += particle.Mass * particle.Y
		}
	}
	q.CenterOfMassX /= q.Mass
	q.CenterOfMassY /= q.Mass

	if len(q.Particles) >= 4 && q.RemainingQuadrants > 0 {
		mx := (q.Xmin + q.Xmax) / 2
		my := (q.Ymin + q.Ymax) / 2
		q.Quadrants = append(q.Quadrants, q.NewQuadrant(q.Xmin, mx, q.Ymin, my, q.RemainingQuadrants-1).SubQuadrant(q.Particles))
		q.Quadrants = append(q.Quadrants, q.NewQuadrant(mx, q.Xmax, q.Ymin, my, q.RemainingQuadrants-1).SubQuadrant(q.Particles))
		q.Quadrants = append(q.Quadrants, q.NewQuadrant(q.Xmin, mx, my, q.Ymax, q.RemainingQuadrants-1).SubQuadrant(q.Particles))
		q.Quadrants = append(q.Quadrants, q.NewQuadrant(mx, q.Xmax, my, q.Ymax, q.RemainingQuadrants-1).SubQuadrant(q.Particles))
	}
	return q
}

func (q *Quadrant) ForceOverParticle(particle Particle, G float64) (float64, float64) {
	if len(q.Particles) == 0 {
		return 0.0, 0.0
	}
	ratio := math.Sqrt(math.Pow(q.CenterOfMassX-particle.X, 2) + math.Pow(q.CenterOfMassY-particle.Y, 2))
	if ratio == 0.0 {
		return 0.0, 0.0
	}

	width := ((q.Xmax - q.Xmin) + (q.Ymax - q.Ymin)) / 2
	if width/ratio < 0.5 {
		dx := q.CenterOfMassX - particle.X
		dy := q.CenterOfMassY - particle.Y
		drSquared := (dx * dx) + (dy * dy)
		m1m2 := particle.Mass * q.Mass
		dr2 := math.Pow(math.Sqrt(drSquared), 2)

		F := (G * m1m2) / dr2
		return F * dx, F * dy
	} else if len(q.Quadrants) == 0 {
		var Fx, Fy float64
		for _, p := range q.Particles {
			if particle.X != p.X && particle.Y != p.Y {
				dx := p.X - particle.X
				dy := p.Y - particle.Y
				drSquared := (dx * dx) + (dy * dy)
				m1m2 := particle.Mass * p.Mass
				dr2 := math.Pow(math.Sqrt(drSquared), 2)

				F := (G * m1m2) / dr2
				Fx += dx * F
				Fy += dy * F
			}
		}
		return Fx, Fy
	} else {
		var Fx, Fy float64
		for _, quadrant := range q.Quadrants {
			qFx, qFy := quadrant.ForceOverParticle(particle, G)
			Fx += qFx
			Fy += qFy
		}
		return Fx, Fy
	}
}
