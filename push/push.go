package push

import rl "github.com/gen2brain/raylib-go/raylib"

type Push struct {
	Position  rl.Vector3
	Direction rl.Vector3
	Radius    float32
	Color     rl.Color
	// F = ma
	Force    float32
	Shooter  Shooter
	LifeTime float32
	Active   bool
}

type Shooter int

const (
	Player Shooter = iota
	Enemy
)

func (p *Push) DrawPush() {
	/*rl.DrawSphere(
		p.Position,
		p.Radius,
		p.Color,
	)*/
}

func (p *Push) Mutate(dt float32) {
	p.LifeTime -= dt
}
