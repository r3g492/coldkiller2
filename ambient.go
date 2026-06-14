package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ambient dust motes that drift slowly through the arena to give the static
// stage a sense of life. They wrap around the player so the visible area is
// always populated, and they never move fast enough to be distracting.

const (
	moteCount   = 70
	moteRange   = 28.0 // half-extent of the box kept around the player
	moteCeiling = 6.0  // motes rise to this height then loop back down
)

type mote struct {
	pos   rl.Vector3
	drift rl.Vector3
	size  float32
	alpha float32
	phase float32
}

var motes []mote

func initAmbient() {
	motes = make([]mote, moteCount)
	for i := range motes {
		motes[i] = newMote()
		motes[i].pos = rl.Vector3{
			X: (rand.Float32()*2 - 1) * moteRange,
			Y: rand.Float32() * moteCeiling,
			Z: (rand.Float32()*2 - 1) * moteRange,
		}
	}
}

func newMote() mote {
	return mote{
		drift: rl.Vector3{
			X: (rand.Float32()*2 - 1) * 0.35,
			Y: rand.Float32()*0.25 + 0.08,
			Z: (rand.Float32()*2 - 1) * 0.35,
		},
		size:  rand.Float32()*0.04 + 0.03,
		alpha: rand.Float32()*0.35 + 0.15,
		phase: rand.Float32() * math.Pi * 2,
	}
}

func updateAmbient(dt float32, center rl.Vector3) {
	for i := range motes {
		m := &motes[i]
		m.phase += dt
		// gentle horizontal sway layered on the base drift
		m.pos.X += (m.drift.X + float32(math.Sin(float64(m.phase)))*0.08) * dt
		m.pos.Z += (m.drift.Z + float32(math.Cos(float64(m.phase*0.8)))*0.08) * dt
		m.pos.Y += m.drift.Y * dt

		if m.pos.Y > moteCeiling {
			m.pos.Y = 0
		}
		if m.pos.X < center.X-moteRange {
			m.pos.X = center.X + moteRange
		} else if m.pos.X > center.X+moteRange {
			m.pos.X = center.X - moteRange
		}
		if m.pos.Z < center.Z-moteRange {
			m.pos.Z = center.Z + moteRange
		} else if m.pos.Z > center.Z+moteRange {
			m.pos.Z = center.Z - moteRange
		}
	}
}

func drawAmbient3D() {
	for i := range motes {
		m := &motes[i]
		c := rl.NewColor(255, 248, 225, uint8(m.alpha*255))
		rl.DrawSphere(m.pos, m.size, c)
	}
}
