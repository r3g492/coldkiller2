package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Persistent scorch/scuff marks left on the floor by deaths and explosions, so
// the battlefield visibly accumulates the fight. Purely cosmetic: flat discs
// laid just above the floor that fade out over several seconds.

const (
	decalFloorY = -1.98 // just above the floor plane at Y = -2
	maxDecals   = 80
)

type decal struct {
	pos     rl.Vector3
	radius  float32
	life    float32
	maxLife float32
	color   rl.Color
}

var decals []decal

func resetDecals() {
	decals = decals[:0]
}

func addDecal(pos rl.Vector3, radius, life float32, color rl.Color) {
	pos.Y = decalFloorY
	decals = append(decals, decal{pos: pos, radius: radius, life: life, maxLife: life, color: color})
	if len(decals) > maxDecals {
		decals = decals[len(decals)-maxDecals:]
	}
}

// addScuff marks where an enemy fell.
func addScuff(pos rl.Vector3) {
	addDecal(pos, 0.5, 6.0, rl.NewColor(22, 20, 26, 150))
}

// addScorch marks where an explosion went off, sized to the blast.
func addScorch(pos rl.Vector3, blastRadius float32) {
	addDecal(pos, blastRadius*0.7, 9.0, rl.NewColor(16, 11, 7, 195))
}

func updateDecals(dt float32) {
	for i := len(decals) - 1; i >= 0; i-- {
		decals[i].life -= dt
		if decals[i].life <= 0 {
			decals = append(decals[:i], decals[i+1:]...)
		}
	}
}

func drawDecals3D() {
	for i := range decals {
		d := &decals[i]
		c := d.color
		c.A = uint8(float32(c.A) * (d.life / d.maxLife))
		rl.DrawCylinder(d.pos, d.radius, d.radius, 0.01, 20, c)
	}
}
