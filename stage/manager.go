package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/structure"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Difficulty       int
	HighestBeaten    int
	StructureManager *structure.Manager
	EnemyManager     *enemy.Manager
	Player           *killer.Killer

	StageElapsed float32
	PenaltyTimer float32
}

const (
	penaltyStartDelay   = float32(60)
	penaltyInterval     = float32(30)
	penaltyStepSeconds  = float32(60)
	penaltyMaxCount     = 12
	penaltyCountStep    = 4
	penaltySpawnRadius  = float32(50)
	penaltySpawnRetries = 10
)

func CreateManager() *Manager {
	return &Manager{}
}

func (m *Manager) GenerateNewStage() {

}

func (m *Manager) Unload() {
}

func (m *Manager) Init(
	structureManager *structure.Manager,
	enemyManager *enemy.Manager,
	player *killer.Killer,
) {
	m.StructureManager = structureManager
	m.EnemyManager = enemyManager
	m.Player = player
}

func (m *Manager) CreateNewStage(pPos rl.Vector3) {
	stageData := Stages[m.Difficulty-1]

	for _, e := range stageData.Enemies {
		m.EnemyManager.Add(e)
	}

	for _, s := range stageData.Structures {
		m.StructureManager.Add(s)
	}

	m.StageElapsed = 0
	m.PenaltyTimer = 0
}

func (m *Manager) Mutate(dt float32) {
	m.StageElapsed += dt
	if m.StageElapsed < penaltyStartDelay {
		return
	}
	if m.StageWon() {
		return
	}
	m.PenaltyTimer -= dt
	if m.PenaltyTimer > 0 {
		return
	}
	m.PenaltyTimer = penaltyInterval

	count := m.penaltyRobotCount()
	for i := 0; i < count; i++ {
		m.spawnPenaltyRobot()
	}
}

func (m *Manager) penaltyRobotCount() int {
	elapsed := m.StageElapsed - penaltyStartDelay
	bucket := int(elapsed/penaltyStepSeconds) + 1
	count := bucket * penaltyCountStep
	if count > penaltyMaxCount {
		count = penaltyMaxCount
	}
	return count
}

func (m *Manager) spawnPenaltyRobot() {
	pPos := m.Player.Position
	size := rl.Vector3{X: killer.CharSize, Y: killer.CharSize, Z: killer.CharSize}
	for range penaltySpawnRetries {
		angle := rand.Float64() * 2 * math.Pi
		pos := rl.Vector3{
			X: pPos.X + float32(math.Cos(angle))*penaltySpawnRadius,
			Y: 0,
			Z: pPos.Z + float32(math.Sin(angle))*penaltySpawnRadius,
		}
		if m.StructureManager.CheckCollision(pos, pos, size) {
			continue
		}
		m.EnemyManager.Add(enemy.SuperRobot(pos.X, pos.Z))
		return
	}
}

func (m *Manager) StageWon() bool {
	return m.EnemyManager.AliveEnemyCount == 0
}

func (m *Manager) GameWon() bool {
	return m.Difficulty > len(Stages)
}

func (m *Manager) StageLost() bool {
	return !m.Player.IsAlive() && m.Player.ActionTimeLeft <= 0
}
