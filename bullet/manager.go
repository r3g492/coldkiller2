package bullet

import (
	"coldkiller2/blast"
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/sound"
	"coldkiller2/structure"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Bullets  []Bullet
	PlayerXp int
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
			Speed:     100.0,
			Radius:    0.03,
			Active:    true,
			LifeTime:  2.0,
			Shooter:   Player,
			Color:     rl.Yellow,
			Damage:    bc.Damage,
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
			Speed:     50.0,
			Radius:    0.03,
			Active:    true,
			LifeTime:  2.0,
			Shooter:   Enemy,
			Color:     rl.Yellow,
			Damage:    bc.Damage,
		}
		bm.Bullets = append(bm.Bullets, b)
	}
}

func (bm *Manager) Mutate(
	dt float32,
	p *killer.Killer,
	el []enemy.Enemy,
	structureManager *structure.Manager,
) []blast.Blast {

	var blasts []blast.Blast
	for i := 0; i < len(bm.Bullets); i++ {
		bm.Bullets[i].Mutate(dt)
		for j := 0; j < len(el); j++ {
			enemyPos := el[j].Position
			enemySize := el[j].Size
			curBullet := bm.Bullets[i]

			if structureManager.CheckCollision(curBullet.Position, rl.Vector3{X: curBullet.Radius, Y: curBullet.Radius, Z: curBullet.Radius}) {
				if bm.Bullets[i].Active {
					blasts = append(blasts, blast.Create(bm.Bullets[i].Position))
					bm.Bullets[i].Active = false
				}
			}

			if curBullet.Shooter == Player && rl.Vector3Distance(enemyPos, curBullet.Position) < enemySize && el[j].Health > 0 {
				if bm.Bullets[i].Active {
					el[j].Damage(bm.Bullets[i].Damage)
					blasts = append(blasts, blast.Create(bm.Bullets[i].Position))
					sound.PlaySound3D(sound.ShotNew, bm.Bullets[i].Position, p.Position, 1)
					bm.Bullets[i].Active = false
					bm.PlayerXp++
				}
			}

			if curBullet.Shooter == Enemy && rl.Vector3Distance(p.Position, curBullet.Position) < p.Size && p.Health > 0 {
				if bm.Bullets[i].Active {
					p.Damage(bm.Bullets[i].Damage)
					blasts = append(blasts, blast.Create(bm.Bullets[i].Position))
					bm.Bullets[i].Active = false
				}
			}
		}

		if bm.Bullets[i].LifeTime <= 0 || !bm.Bullets[i].Active {
			bm.Bullets[i] = bm.Bullets[len(bm.Bullets)-1]
			bm.Bullets = bm.Bullets[:len(bm.Bullets)-1]
			i--
		}
	}
	return blasts
}

func (bm *Manager) Draw3D() {
	for _, b := range bm.Bullets {
		b.Draw3D()
	}
}

func (bm *Manager) Unload() {
	bm.Bullets = []Bullet{}
	bm.PlayerXp = 0
}
