package blast

import rl "github.com/gen2brain/raylib-go/raylib"

var RADIUS float32 = 0.25
var COLOR = rl.Yellow
var LIFETIME float32 = 0.06

type Blast struct {
	Position rl.Vector3
	Radius   float32
	LifeTime float32
	Color    rl.Color
}

func (b *Blast) Draw() {
	rl.DrawSphere(b.Position, b.Radius, b.Color)
}

func (b *Blast) Mutate(dt float32) {
	b.LifeTime -= dt
}

func Create(
	position rl.Vector3,
) Blast {
	return Blast{
		Position: position,
		Radius:   RADIUS,
		LifeTime: LIFETIME,
		Color:    COLOR,
	}
}
