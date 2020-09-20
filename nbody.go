package main

import (
	"crypto/rand"
	"math"
	"math/big"
)

type Nbodys struct {
	Particles Particles
}

func (n *Nbodys) randFloat(max float64) float64 {
	random, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return float64(random.Int64())
	//return min + (float64(random.Int64())/max)*(max-min)
}

func (n *Nbodys) AddRandomParticles(maxX, maxY float64, numberOfParticles int) {
	for i := 0; i < numberOfParticles; i++ {
		n.Particles = append(
			n.Particles,
			Particle{
				X:    n.randFloat(maxX),
				Y:    n.randFloat(maxY),
				VX:   0.0,
				VY:   0.0,
				Mass: 1.0,
			},
		)
	}
}

func (n *Nbodys) Clear() {
	n.Particles = Particles{}
}

func (n *Nbodys) Compute(G, DT float64) {
	for i := range n.Particles {
		var Fx, Fy float64
		for j := range n.Particles {
			if j != i {
				dx := n.Particles[j].X - n.Particles[i].X
				dy := n.Particles[j].Y - n.Particles[i].Y
				drSquared := (dx * dx) + (dy * dy)
				m1m2 := n.Particles[i].Mass * n.Particles[j].Mass
				dr2 := math.Pow(math.Sqrt(drSquared), 2.0)
				F := (G * m1m2) / dr2
				Fx += dx * F
				Fy += dy * F
			}
		}
		n.Particles[i].VX += DT * Fx
		n.Particles[i].VY += DT * Fy
	}
	for i := range n.Particles {
		n.Particles[i].X += n.Particles[i].VX * DT
		n.Particles[i].Y += n.Particles[i].VY * DT
	}
}
