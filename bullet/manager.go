package bullet

import (
	"coldkiller2/blast"
	"coldkiller2/enemy"
	"coldkiller2/killer"
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

func (bm *Manager) Init() {
	bm.Bullets = []Bullet{}
	bm.PlayerXp = 0
}

func (bm *Manager) KillerBulletCreate(
	bulletCmds []killer.BulletCmd,
) {
	for _, bc := range bulletCmds {
		b := Bullet{
			Position:  bc.Pos,
			Direction: bc.Dir,
			Speed:     75,
			Radius:    0.1,
			Active:    true,
			LifeTime:  1.0,
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
		const speed = 75
		lifeTime := bc.Range / speed
		b := Bullet{
			Position:     bc.Pos,
			Direction:    bc.Dir,
			Speed:        speed,
			Radius:       0.1,
			Active:       true,
			LifeTime:     lifeTime,
			Shooter:      Enemy,
			EnemyShooter: bc.Shooter,
			Color:        rl.Yellow,
			Damage:       bc.Damage,
		}
		bm.Bullets = append(bm.Bullets, b)
	}
}

func (bm *Manager) Mutate(
	dt float32,
	p *killer.Killer,
	el []*enemy.Enemy,
	structureManager *structure.Manager,
) []blast.Blast {

	var blasts []blast.Blast
	for i := 0; i < len(bm.Bullets); i++ {
		bm.Bullets[i].Mutate(dt)
		curBullet := bm.Bullets[i]

		if bm.Bullets[i].LifeTime >= 0.99 && bm.Bullets[i].Active {
			blasts = append(blasts, blast.Create(bm.Bullets[i].Position, bm.Bullets[i].Shooter == Player))
		}

		if structureManager.CheckCollision(curBullet.Position, curBullet.PrevPosition, rl.Vector3{X: curBullet.Radius, Y: curBullet.Radius, Z: curBullet.Radius}) {
			if bm.Bullets[i].Active {
				blasts = append(blasts, blast.CreateBig(bm.Bullets[i].Position, bm.Bullets[i].Shooter == Player))
				bm.Bullets[i].Active = false
			}
		}

		for j := 0; j < len(el); j++ {
			enemyPos := el[j].Position

			if (curBullet.Shooter == Player || (curBullet.Shooter == Enemy && el[j] != curBullet.EnemyShooter)) && el[j].IsAlive() {
				hitRadius := el[j].Size + curBullet.Radius
				if checkSegmentSphereCollision(curBullet.PrevPosition, curBullet.Position, enemyPos, hitRadius) {
					if bm.Bullets[i].Active {
						dir := bm.Bullets[i].Direction
						hitPos := bm.Bullets[i].Position
						el[j].Damage(bm.Bullets[i].Damage, dir)
						blasts = append(blasts, blast.CreateSplash(hitPos))
						for k := 0; k < 3; k++ {
							offset := rl.Vector3Scale(dir, float32(k+1)*0.3)
							blasts = append(blasts, blast.CreateDebris(rl.Vector3Add(hitPos, offset)))
						}
						bm.Bullets[i].Active = false
						bm.PlayerXp++
					}
				}
			}

			if curBullet.Shooter == Enemy && p.IsAlive() {
				hitRadius := p.Size + curBullet.Radius
				if checkSegmentSphereCollision(curBullet.PrevPosition, curBullet.Position, p.Position, hitRadius) {
					if bm.Bullets[i].Active {
						p.Damage(bm.Bullets[i].Damage)
						blasts = append(blasts, blast.Create(bm.Bullets[i].Position, bm.Bullets[i].Shooter == Player))
						bm.Bullets[i].Active = false
					}
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

func checkSegmentSphereCollision(start, end, sphereCenter rl.Vector3, radius float32) bool {
	ab := rl.Vector3Subtract(end, start)
	ap := rl.Vector3Subtract(sphereCenter, start)

	abLenSq := rl.Vector3DotProduct(ab, ab)
	if abLenSq == 0 {
		return rl.Vector3Distance(start, sphereCenter) < radius
	}

	t := rl.Vector3DotProduct(ap, ab) / abLenSq

	if t < 0.0 {
		t = 0.0
	}
	if t > 1.0 {
		t = 1.0
	}

	closestPoint := rl.Vector3Add(start, rl.Vector3Scale(ab, t))

	return rl.Vector3Distance(closestPoint, sphereCenter) < radius
}
