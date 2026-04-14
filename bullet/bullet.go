package bullet

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	Position     rl.Vector3
	PrevPosition rl.Vector3
	Direction    rl.Vector3
	Speed        float32
	Radius       float32
	Active       bool
	LifeTime     float32
	Shooter      Shooter
	Color        rl.Color
	Damage       int32

	IsHiddenFromKiller bool
}

type Shooter int

const (
	Player Shooter = iota
	Enemy
)

func (b *Bullet) Draw3D() {
	if b.IsHiddenFromKiller {
		return
	}
	tailStart := rl.Vector3Subtract(b.Position, rl.Vector3Scale(b.Direction, 0.8))
	tailColor := rl.NewColor(b.Color.R, b.Color.G, b.Color.B, 100)
	rl.DrawCapsule(tailStart, b.Position, b.Radius*0.3, 4, 1, tailColor)
	rl.DrawSphere(b.Position, b.Radius, b.Color)
}

func (b *Bullet) Mutate(dt float32) {
	movement := rl.Vector3Scale(b.Direction, b.Speed*dt)
	b.PrevPosition = b.Position
	b.Position = rl.Vector3Add(b.Position, movement)
	b.LifeTime -= dt
}
