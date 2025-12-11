package enemy

import (
	"coldkiller2/killer"
	"coldkiller2/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Enemies []Enemy
}

func CreateManager() *Manager {
	return &Manager{}
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
		AnimationIdx:          2,
		AnimationCurrentFrame: 0,
		AnimationFrameCounter: 0,
		AnimationFrameSpeed:   24,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:              enemyPosition,
		Size:                  2,
		MoveSpeed:             10.0,
		AttackSound:           shotGunSound,
		ActionTimeLeft:        0,
		PushedTimeLeft:        0,
		Health:                100,
		IsDead:                false,
	}
	em.Enemies = append(em.Enemies, addEnemy1)

	addEnemy2 := Enemy{
		Model:                 enemyModel,
		ModelAngleDeg:         0,
		Animation:             enemyAnimation,
		AnimationIdx:          2,
		AnimationCurrentFrame: 0,
		AnimationFrameCounter: 0,
		AnimationFrameSpeed:   24,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:              rl.Vector3{X: 5, Y: 0, Z: 0},
		Size:                  2,
		MoveSpeed:             10.0,
		AttackSound:           shotGunSound,
		ActionTimeLeft:        0,
		PushedTimeLeft:        0,
		Health:                100,
		IsDead:                false,
	}
	em.Enemies = append(em.Enemies, addEnemy2)
}

func (em *Manager) Mutate(dt float32, p *killer.Killer) ([]BulletCmd, []PushCmd) {
	var bulletCmds []BulletCmd
	var pushCmds []PushCmd
	for i := 0; i < len(em.Enemies); i++ {
		var addBullets, addPush = em.Enemies[i].Mutate(dt)
		bulletCmds = append(bulletCmds, addBullets...)
		pushCmds = append(pushCmds, addPush...)
		if em.Enemies[i].IsDead {
			em.Enemies[i] = em.Enemies[len(em.Enemies)-1]
			em.Enemies = em.Enemies[:len(em.Enemies)-1]
			i--
		}
	}
	return bulletCmds, pushCmds
}

func (em *Manager) DrawEnemies3D() {
	for i := range em.Enemies {
		em.Enemies[i].Draw3D()
	}
}

func (em *Manager) PlanAnimate(dt float32) {
	for i, _ := range em.Enemies {
		em.Enemies[i].PlanAnimate(dt)
	}
}
