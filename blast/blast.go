package blast

import rl "github.com/gen2brain/raylib-go/raylib"

var RADIUS float32 = 0.25
var COLOR = rl.Yellow
var LIFETIME float32 = 0.06

type Blast struct {
	Position    rl.Vector3
	Radius      float32
	MaxLifeTime float32
	LifeTime    float32
	Color       rl.Color

	IsHiddenFromKiller bool
}

func (b *Blast) Draw3D() {
	if b.IsHiddenFromKiller {
		return
	}

	rl.DrawSphere(b.Position, b.Radius, b.Color)
}

func (b *Blast) Mutate(dt float32) {
	b.LifeTime -= dt

	lifeRatio := b.LifeTime / b.MaxLifeTime
	if lifeRatio < 0 {
		lifeRatio = 0
	}

	b.Radius = RADIUS * lifeRatio
}

func Create(position rl.Vector3) Blast {
	return Blast{
		Position:    position,
		Radius:      RADIUS,
		MaxLifeTime: LIFETIME,
		LifeTime:    LIFETIME,
		Color:       COLOR,
	}
}
