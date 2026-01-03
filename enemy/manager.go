package enemy

import (
	"coldkiller2/animation"
	"coldkiller2/killer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Enemies          []Enemy
	SharedModel      rl.Model
	SharedAnimations []rl.ModelAnimation
}

func CreateManager() *Manager {
	return &Manager{}
}

func (em *Manager) Init() {
	em.Enemies = make([]Enemy, 0)

	// Load assets ONCE
	em.SharedModel = rl.LoadModel("resources/unit_v3.glb")
	em.SharedAnimations = rl.LoadModelAnimations("resources/unit_v3.glb")

	// TODO: change unit init
	enemyPosition := rl.Vector3{X: 10, Y: 0, Z: 0}
	addEnemy1 := Enemy{
		Model:           em.SharedModel,
		ModelAngleDeg:   0,
		Animation:       em.SharedAnimations,
		MoveDirection:   rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection: rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:        enemyPosition,
		Size:            1.5,
		MoveSpeed:       2.0,
		ActionTimeLeft:  0,
		Health:          100,
		IsDead:          false,
		AttackRange:     15,
		AimTimeLeft:     2,
		AimTimeUnit:     2,
	}
	em.Enemies = append(em.Enemies, addEnemy1)

	addEnemy2 := Enemy{
		Model:           em.SharedModel,
		ModelAngleDeg:   0,
		Animation:       em.SharedAnimations,
		AnimationState:  animation.StateIdle,
		MoveDirection:   rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection: rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:        rl.Vector3{X: 5, Y: 0, Z: 0},
		Size:            1.5,
		MoveSpeed:       2.0,
		ActionTimeLeft:  0,
		Health:          100,
		IsDead:          false,
		AttackRange:     15,
		AimTimeLeft:     2,
		AimTimeUnit:     2,
	}
	em.Enemies = append(em.Enemies, addEnemy2)
}

func (em *Manager) Mutate(dt float32, p *killer.Killer) []BulletCmd {
	var bulletCmds []BulletCmd
	allEnemyBoxes := em.GetEnemyBoundingBoxes()
	for i := 0; i < len(em.Enemies); i++ {
		var addBullets = em.Enemies[i].Mutate(dt, *p, allEnemyBoxes, i)
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
	for i := range em.Enemies {
		em.Enemies[i].Draw3D()
	}
}

func (em *Manager) ProcessAnimation(dt float32) {
	for i, _ := range em.Enemies {
		em.Enemies[i].ResolveAnimation()
		em.Enemies[i].PlanAnimate(dt)
		em.Enemies[i].Animate()
	}
}

func (em *Manager) Unload() {
	rl.UnloadModel(em.SharedModel)
	rl.UnloadModelAnimations(em.SharedAnimations)
	em.Enemies = nil
}

func (em *Manager) GetEnemyBoundingBoxes() []rl.BoundingBox {
	boxes := make([]rl.BoundingBox, 0, len(em.Enemies))
	for _, e := range em.Enemies {
		if e.Health > 0 {
			boxes = append(boxes, e.GetBoundingBox())
		}
	}
	return boxes
}

func (e *Enemy) isColliding(myIdx int, obstacles []rl.BoundingBox, killerObstacle rl.BoundingBox) bool {
	myBox := e.GetBoundingBox()
	for i, box := range obstacles {
		if i == myIdx {
			continue
		}
		if rl.CheckCollisionBoxes(myBox, box) {
			return true
		}
	}
	if rl.CheckCollisionBoxes(myBox, killerObstacle) {
		return true
	}
	return false
}
