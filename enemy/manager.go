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
	BulletBuffer               []BulletCmd
	Grid                       map[int][]int
	CellSize                   float32
}

const CELL_SIZE = 5.0
const PLAYER_SEE_AREA = 35.0
const PLAYER_3D_SEE_SQUARE = 1000

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
	em.CellSize = CELL_SIZE
	em.Grid = make(map[int][]int)
	em.Generate(p)
}

func (em *Manager) Mutate(dt float32, p *killer.Killer) []BulletCmd {
	em.updateGrid()
	em.BulletBuffer = em.BulletBuffer[:0]

	for i := 0; i < len(em.Enemies); i++ {
		addBullets := em.Enemies[i].Mutate(dt, *p, em, i)
		em.BulletBuffer = append(em.BulletBuffer, addBullets...)
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

	return em.BulletBuffer
}

func (em *Manager) DrawEnemies3D(p *killer.Killer) {
	for i := range em.Enemies {
		if rl.CheckCollisionSpheres(em.Enemies[i].Position, 1.0, p.Position, PLAYER_SEE_AREA) {
			em.Enemies[i].Draw3D(p)
		}
	}
}

func (em *Manager) DrawEnemiesUi(p *killer.Killer) {
	for i := range em.Enemies {
		if rl.Vector3DistanceSqr(em.Enemies[i].Position, p.Position) < PLAYER_3D_SEE_SQUARE {
			em.Enemies[i].DrawUI(p)
		}
	}
}

func (em *Manager) ProcessAnimation(dt float32, p *killer.Killer) {
	for i := range em.Enemies {
		if rl.Vector3DistanceSqr(em.Enemies[i].Position, p.Position) > PLAYER_3D_SEE_SQUARE {
			continue
		}
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
		MoveSpeed:       10.0,
		ActionTimeLeft:  0,
		Health:          100,
		IsDead:          false,
		AttackRange:     10,
		AimTimeLeft:     2,
		AimTimeUnit:     2,
	}
	trialLimit := 3
	trial := 0
	for trial < trialLimit && candidate.isCollidingWithGrid(-1, em, p.GetBoundingBox()) {
		candidatePosition = getRandomPosition(p)
		candidate.Position = candidatePosition
		trial++
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

func (em *Manager) getGridID(pos rl.Vector3) int {
	gx := int(math.Floor(float64(pos.X / em.CellSize)))
	gz := int(math.Floor(float64(pos.Z / em.CellSize)))
	return gx*73856093 ^ gz*19349663
}

func (em *Manager) updateGrid() {
	for k := range em.Grid {
		em.Grid[k] = em.Grid[k][:0]
	}
	for i, e := range em.Enemies {
		if !e.IsAlive() {
			continue
		}
		id := em.getGridID(e.Position)
		em.Grid[id] = append(em.Grid[id], i)
	}
}

func (e *Enemy) isCollidingWithGrid(myIdx int, em *Manager, playerBox rl.BoundingBox) bool {
	myBox := e.GetBoundingBox()

	for x := float32(-1); x <= 1; x++ {
		for z := float32(-1); z <= 1; z++ {
			checkPos := rl.Vector3{
				X: e.Position.X + (x * em.CellSize),
				Z: e.Position.Z + (z * em.CellSize),
			}
			gridID := em.getGridID(checkPos)

			for _, otherIdx := range em.Grid[gridID] {
				if otherIdx == myIdx {
					continue
				}

				other := em.Enemies[otherIdx]
				if rl.CheckCollisionBoxes(myBox, other.GetBoundingBox()) {
					return true
				}
			}
		}
	}

	return rl.CheckCollisionBoxes(myBox, playerBox)
}
