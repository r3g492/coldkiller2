package enemy

import (
	"coldkiller2/killer"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Enemies                    []Enemy
	SharedModel                rl.Model
	SharedAnimations           []rl.ModelAnimation
	EnemyGenerationLevel       int
	EnemyGenerationLevelUpUnit time.Duration
	LastLevelUp                time.Time
	EnemyGenerateUnit          time.Duration
	LastGenerated              time.Time
	EnemyLimit                 int
}

func CreateManager() *Manager {
	return &Manager{}
}

func (em *Manager) Init(p *killer.Killer) {
	em.Enemies = make([]Enemy, 0)
	em.SharedModel = rl.LoadModel("resources/unit_v3.glb")
	em.SharedAnimations = rl.LoadModelAnimations("resources/unit_v3.glb")
	em.EnemyGenerationLevel = 0
	em.EnemyGenerationLevelUpUnit = 8 * time.Second
	em.LastLevelUp = time.Now()
	em.EnemyGenerateUnit = 4 * time.Second
	em.LastGenerated = time.Now()
	em.EnemyLimit = 100
	em.Generate(p)
}

func (em *Manager) Mutate(dt float32, p *killer.Killer) []BulletCmd {
	var bulletCmds []BulletCmd

	for i := 0; i < len(em.Enemies); i++ {
		addBullets := em.Enemies[i].Mutate(dt, *p, em.Enemies, i)
		bulletCmds = append(bulletCmds, addBullets...)
	}

	for i := len(em.Enemies) - 1; i >= 0; i-- {
		if em.Enemies[i].IsDead {
			em.Enemies = append(em.Enemies[:i], em.Enemies[i+1:]...)
		}
	}

	if time.Since(em.LastGenerated) > em.EnemyGenerateUnit {
		em.Generate(p)
	}

	if time.Since(em.LastLevelUp) > em.EnemyGenerationLevelUpUnit {
		em.UpTheTempo()
	}

	return bulletCmds
}

func (em *Manager) DrawEnemies3D(p *killer.Killer) {
	for i := range em.Enemies {
		em.Enemies[i].Draw3D(p)
	}
}

func (em *Manager) DrawEnemiesUi(p *killer.Killer) {
	for i := range em.Enemies {
		em.Enemies[i].DrawUI(p)
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
	em.Enemies = []Enemy{}
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

func (e *Enemy) isColliding(myIdx int, others []Enemy, killerObstacle rl.BoundingBox) bool {
	myBox := e.GetBoundingBox()

	for i, other := range others {
		if i == myIdx || !other.IsAlive() {
			continue
		}

		if rl.CheckCollisionBoxes(myBox, other.GetBoundingBox()) {
			return true
		}
	}

	if rl.CheckCollisionBoxes(myBox, killerObstacle) {
		return true
	}
	return false
}

func (em *Manager) Generate(p *killer.Killer) {
	for i := 0; i <= em.EnemyGenerationLevel; i++ {
		em.addEnemy(p)
		em.addEnemy(p)
	}
	em.LastGenerated = time.Now()
}

func (em *Manager) addEnemy(p *killer.Killer) {
	if em.EnemyLimit < len(em.Enemies) {
		return
	}

	candidatePosition := getRandomPosition(p)
	candidate := Enemy{
		Model:           em.SharedModel,
		ModelAngleDeg:   0,
		Animation:       em.SharedAnimations,
		MoveDirection:   rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection: rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:        candidatePosition,
		Size:            1.0,
		MoveSpeed:       5.0,
		ActionTimeLeft:  0,
		Health:          100,
		IsDead:          false,
		AttackRange:     15,
		AimTimeLeft:     2,
		AimTimeUnit:     2,
	}
	trialLimit := 3
	trial := 0
	for ; candidate.isColliding(-1, em.Enemies, p.GetBoundingBox()) && trial < trialLimit; trial++ {
		candidatePosition = getRandomPosition(p)
		candidate = Enemy{
			Model:           em.SharedModel,
			ModelAngleDeg:   0,
			Animation:       em.SharedAnimations,
			MoveDirection:   rl.Vector3{X: 0, Y: 0, Z: 0},
			TargetDirection: rl.Vector3{X: 0, Y: 0, Z: 0},
			Position:        candidatePosition,
			Size:            1.0,
			MoveSpeed:       5.0,
			ActionTimeLeft:  0,
			Health:          100,
			IsDead:          false,
			AttackRange:     15,
			AimTimeLeft:     2,
			AimTimeUnit:     2,
		}
	}
	if trial >= trialLimit {
		return
	}
	em.Enemies = append(em.Enemies, candidate)
}

func getRandomPosition(p *killer.Killer) rl.Vector3 {
	angle := rand.Float64() * 2 * math.Pi
	distance := float32(28.0)
	offsetX := float32(math.Cos(angle)) * distance
	offsetZ := float32(math.Sin(angle)) * distance
	return rl.Vector3{
		X: p.Position.X + offsetX,
		Y: 0,
		Z: p.Position.Z + offsetZ,
	}
}

func (em *Manager) UpTheTempo() {
	em.EnemyGenerationLevel++
	em.LastLevelUp = time.Now()
}
