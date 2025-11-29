package killer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	Position  rl.Vector3
	Direction rl.Vector3
	Speed     float32
	Radius    float32
	Active    bool
	LifeTime  float32
}

func NewBullet(pos, dir rl.Vector3) Bullet {
	return Bullet{
		Position:  pos,
		Direction: dir,
		Speed:     40.0,
		Radius:    0.3,
		Active:    true,
		LifeTime:  2.0,
	}
}
