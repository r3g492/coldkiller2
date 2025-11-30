package bullet

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
	Shooter   Shooter
}

type Shooter int

const (
	Player Shooter = iota
	Enemy
)

func (b *Bullet) DrawBullet() {
	rl.DrawSphere(b.Position, b.Radius, rl.Yellow)
}

func (b *Bullet) Mutate(dt float32) {
	movement := rl.Vector3Scale(b.Direction, b.Speed*dt)
	b.Position = rl.Vector3Add(b.Position, movement)
	b.LifeTime -= dt
}

type Manager struct {
	Bullets []Bullet
}

func (bm *Manager) NewPlayerBullet(
	pos,
	dir rl.Vector3,
) {
	b := Bullet{
		Position:  pos,
		Direction: dir,
		Speed:     40.0,
		Radius:    0.2,
		Active:    true,
		LifeTime:  2.0,
		Shooter:   Player,
	}
	bm.Bullets = append(bm.Bullets, b)
}

func (bm *Manager) Mutate(dt float32) {
	for i := 0; i < len(bm.Bullets); i++ {
		bm.Bullets[i].Mutate(dt)
		if bm.Bullets[i].LifeTime <= 0 || !bm.Bullets[i].Active {
			bm.Bullets[i] = bm.Bullets[len(bm.Bullets)-1]
			bm.Bullets = bm.Bullets[:len(bm.Bullets)-1]
			i--
		}
	}
}

func (bm *Manager) DrawBullets() {
	for _, b := range bm.Bullets {
		b.DrawBullet()
	}
}
