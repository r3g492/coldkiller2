package bullet

import (
	"coldkiller2/enemy"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	Position       rl.Vector3
	PrevPosition   rl.Vector3
	Direction      rl.Vector3
	Speed          float32
	Radius         float32
	Active         bool
	LifeTime       float32
	Age            float32
	SlowTimeActive bool
	Shooter        Shooter
	EnemyShooter   *enemy.Enemy
	Color          rl.Color
	Damage         int32

	IsHiddenFromKiller bool
}

type Shooter int

const (
	Player Shooter = iota
	Enemy
)

const initialDelay = 0.015

func (b *Bullet) Draw3D() {
	if b.IsHiddenFromKiller {
		return
	}
	var delay float32 = initialDelay
	if b.SlowTimeActive {
		delay = initialDelay * 2
	}
	if b.Age >= delay {
		tailStart := rl.Vector3Subtract(b.Position, rl.Vector3Scale(b.Direction, 5))
		tailColor := rl.NewColor(b.Color.R, b.Color.G, b.Color.B, 100)
		rl.DrawCapsule(tailStart, b.Position, b.Radius*0.3, 4, 1, tailColor)
		rl.DrawSphere(b.Position, b.Radius, b.Color)
	}
}

func (b *Bullet) Mutate(dt float32) {
	movement := rl.Vector3Scale(b.Direction, b.Speed*dt)
	b.PrevPosition = b.Position
	b.Position = rl.Vector3Add(b.Position, movement)
	b.LifeTime -= dt
	b.Age += dt
}
