package enemy

import (
	"coldkiller2/killer"
	"coldkiller2/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Enemies []Enemy
}

func (em *Manager) Init() {
	enemyModel := rl.LoadModel("resources/robot.glb")
	enemyAnimation := rl.LoadModelAnimations("resources/robot.glb")
	enemyPosition := rl.Vector3{X: 0, Y: 0, Z: 0}
	shotGunSound := util.LoadSoundFromEmbedded("shotgun-03-38220.mp3")

	addEnemy1 := Enemy{
		Model:                 enemyModel,
		ModelAngleDeg:         0,
		Animation:             enemyAnimation,
		AnimationIdx:          0,
		AnimationCurrentFrame: 0,
		AnimationFrameCounter: 0,
		AnimationFrameSpeed:   0.1,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:              enemyPosition,
		Size:                  2,
		MoveSpeed:             10.0,
		AttackSound:           shotGunSound,
		ActionTimeLeft:        0,
		Health:                100,
		IsDead:                false,
	}
	em.Enemies = append(em.Enemies, addEnemy1)
}

func (em *Manager) Mutate(dt float32, p *killer.Killer) []BulletCmd {
	var bulletCmds []BulletCmd
	for i := 0; i < len(em.Enemies); i++ {
		var addBullets []BulletCmd = em.Enemies[i].Mutate(dt)
		bulletCmds = append(bulletCmds, addBullets...)
		if em.Enemies[i].IsDead {
			em.Enemies[i] = em.Enemies[len(em.Enemies)-1]
			em.Enemies = em.Enemies[:len(em.Enemies)-1]
			i--
		}
	}
	return bulletCmds
}

func (em *Manager) DrawEnemies3D() {
	for _, e := range em.Enemies {
		e.Draw3D()
	}
}
