package bullet

import (
	"coldkiller2/enemy"
	"coldkiller2/killer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Bullets []Bullet
}

func CreateManager() *Manager {
	return &Manager{}
}

func (bm *Manager) KillerBulletCreate(
	bulletCmds []killer.BulletCmd,
) {
	for _, bc := range bulletCmds {
		b := Bullet{
			Position:  bc.Pos,
			Direction: bc.Dir,
			Speed:     40.0,
			Radius:    0.2,
			Active:    true,
			LifeTime:  2.0,
			Shooter:   Player,
			Color:     rl.Yellow,
		}
		bm.Bullets = append(bm.Bullets, b)
	}
}

func (bm *Manager) EnemyBulletCreate(
	bulletCmds []enemy.BulletCmd,
) {
	for _, bc := range bulletCmds {
		b := Bullet{
			Position:  bc.Pos,
			Direction: bc.Dir,
			Speed:     40.0,
			Radius:    0.2,
			Active:    true,
			LifeTime:  2.0,
			Shooter:   Enemy,
			Color:     rl.Red,
		}
		bm.Bullets = append(bm.Bullets, b)
	}
}

func (bm *Manager) Mutate(dt float32, p *killer.Killer, el []enemy.Enemy) {
	for i := 0; i < len(bm.Bullets); i++ {
		bm.Bullets[i].Mutate(dt)
		for j := 0; j < len(el); j++ {
			enemyPos := el[j].Position
			enemySize := el[j].Size
			curBullet := bm.Bullets[i]
			if rl.Vector3Distance(enemyPos, curBullet.Position) < enemySize && el[j].Health > 0 {
				el[j].Damage(50)
				bm.Bullets[i].Active = false
			}
		}

		if bm.Bullets[i].LifeTime <= 0 || !bm.Bullets[i].Active {
			bm.Bullets[i] = bm.Bullets[len(bm.Bullets)-1]
			bm.Bullets = bm.Bullets[:len(bm.Bullets)-1]
			i--
		}
	}
}

func (bm *Manager) DrawBullets3D() {
	for _, b := range bm.Bullets {
		b.DrawBullet()
	}
}
