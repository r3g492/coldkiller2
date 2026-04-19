package enemy

import (
	"coldkiller2/killer"
	"coldkiller2/model"
	"coldkiller2/structure"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Enemies                    []*Enemy
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
	AliveEnemyCount            int
}

const CellSize = 5.0
const PlayerSeeArea = 35.0
const Player3dSeeSquare = 1000

func CreateManager() *Manager {
	return &Manager{}
}

func (em *Manager) Init(p *killer.Killer) {
	em.Enemies = make([]*Enemy, 0)
	em.SharedModel = model.SoldierModel
	em.SharedAnimations = model.SoldierAnimation
	em.EnemyGenerationLevel = 0
	em.EnemyGenerationLevelUpUnit = 8 * time.Second
	em.LastLevelUp = time.Now()
	em.EnemyGenerateUnit = 4 * time.Second
	em.LastGenerated = time.Now()
	em.EnemyLimit = 100
	em.CellSize = CellSize
	em.Grid = make(map[int][]int)
}

func (em *Manager) Mutate(
	dt float32,
	p *killer.Killer,
	structureManager *structure.Manager,
) []BulletCmd {
	em.AliveEnemyCount = 0
	em.updateGrid()
	em.BulletBuffer = em.BulletBuffer[:0]

	// dash push: player collides with enemies while dashing
	if p.DashPushTimeLeft > 0 {
		for _, e := range em.Enemies {
			if !e.IsAlive() {
				continue
			}
			if rl.CheckCollisionSpheres(p.Position, p.Size*1.5, e.Position, e.Size) {
				dir := rl.Vector3Subtract(e.Position, p.Position)
				if rl.Vector3LengthSqr(dir) < 0.0001 {
					dir = rl.Vector3{X: 1}
				}
				e.ApplyKnockback(rl.Vector3Scale(rl.Vector3Normalize(dir), 25.0), 0.25)
			}
		}
	}

	for i := 0; i < len(em.Enemies); i++ {
		addBullets := em.Enemies[i].Mutate(dt, *p, em, i, structureManager)
		em.BulletBuffer = append(em.BulletBuffer, addBullets...)
	}

	for i := len(em.Enemies) - 1; i >= 0; i-- {
		if em.Enemies[i].ShouldBeDeleted {
			em.Enemies = append(em.Enemies[:i], em.Enemies[i+1:]...)
			continue
		}
		if em.Enemies[i].IsAlive() {
			em.AliveEnemyCount++
		}
	}

	if time.Since(em.LastLevelUp) > em.EnemyGenerationLevelUpUnit {
		em.UpTheTempo()
	}

	return em.BulletBuffer
}

func (em *Manager) Draw3D(p *killer.Killer) {
	for i := range em.Enemies {
		if rl.CheckCollisionSpheres(em.Enemies[i].Position, 1.0, p.Position, PlayerSeeArea) {
			em.Enemies[i].Draw3D(p)
		}
	}
}

func (em *Manager) DrawUi(p *killer.Killer) {
	for i := range em.Enemies {
		if rl.Vector3DistanceSqr(em.Enemies[i].Position, p.Position) < Player3dSeeSquare {
			em.Enemies[i].DrawUI(p)
		}
	}
}

func (em *Manager) ProcessAnimation(dt float32, p *killer.Killer) {
	for i := range em.Enemies {
		if rl.Vector3DistanceSqr(em.Enemies[i].Position, p.Position) > Player3dSeeSquare {
			continue
		}
		em.Enemies[i].ResolveAnimation()
		em.Enemies[i].PlanAnimate(dt)
	}
}

func (em *Manager) Unload() {
	rl.UnloadModel(em.SharedModel)
	if len(em.SharedAnimations) > 0 {
		rl.UnloadModelAnimations(em.SharedAnimations)
	}
	em.Enemies = []*Enemy{}
}

func (em *Manager) AliveCount() int {
	count := 0
	for _, e := range em.Enemies {
		if e.IsAlive() {
			count++
		}
	}
	return count
}

func (em *Manager) GetBoundingBoxes() []rl.BoundingBox {
	boxes := make([]rl.BoundingBox, 0, len(em.Enemies))
	for _, e := range em.Enemies {
		if e.IsAlive() {
			boxes = append(boxes, e.GetBoundingBox())
		}
	}
	return boxes
}

func (em *Manager) findCollidingEnemy(excludeIdx int, e *Enemy) *Enemy {
	myBox := e.GetBoundingBox()
	for i, other := range em.Enemies {
		if i == excludeIdx || !other.IsAlive() {
			continue
		}
		if rl.CheckCollisionBoxes(myBox, other.GetBoundingBox()) {
			return other
		}
	}
	return nil
}

func (em *Manager) Add(e *Enemy) {
	if em.EnemyLimit < len(em.Enemies) {
		return
	}
	em.Enemies = append(em.Enemies, e)
	newEnemyIndex := len(em.Enemies) - 1
	gridID := em.getGridID(e.Position)
	em.Grid[gridID] = append(em.Grid[gridID], newEnemyIndex)
}

func getRandomPosition(pPos rl.Vector3) rl.Vector3 {
	angle := rand.Float64() * 2 * math.Pi
	distance := float32(28.0)
	offsetX := float32(math.Cos(angle)) * distance
	offsetZ := float32(math.Sin(angle)) * distance
	return rl.Vector3{
		X: pPos.X + offsetX,
		Y: 0,
		Z: pPos.Z + offsetZ,
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
