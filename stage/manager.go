package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/structure"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Difficulty           int
	GameEndingDifficulty int
	StructureManager     *structure.Manager
	EnemyManager         *enemy.Manager
	Player               *killer.Killer
}

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
	m.GameEndingDifficulty = 100
	m.StructureManager = structureManager
	m.EnemyManager = enemyManager
	m.Player = player
}

func (m *Manager) CreateNewStage(pPos rl.Vector3) {
	// TODO: Difficulty에 따라 enemy와 structure 변주 추가
	m.EnemyManager.AddEnemy(pPos)
	initialStructures := []structure.Structure{
		{
			Position:  rl.Vector3{X: 5, Y: 0, Z: 5},
			Size:      rl.Vector3{X: 1, Y: 1, Z: 10},
			Direction: rl.Vector3{X: 1, Y: 0.2, Z: 1},
			Color:     rl.DarkGray,
		},
		{
			Position:  rl.Vector3{X: -5, Y: 0, Z: -5},
			Size:      rl.Vector3{X: 1, Y: 1, Z: 10},
			Direction: rl.Vector3{X: 0, Y: 0, Z: 0},
			Color:     rl.DarkGray,
		},
		{
			Position:  rl.Vector3{X: 10, Y: 0, Z: 0},
			Size:      rl.Vector3{X: 2, Y: 2, Z: 2},
			Direction: rl.Vector3{X: 0, Y: 0, Z: 0},
			Color:     rl.DarkGray,
		},
		{
			Position:  rl.Vector3{X: 15, Y: 0, Z: 0},
			Size:      rl.Vector3{X: 2, Y: 2, Z: 2},
			Direction: rl.Vector3{X: 0, Y: 0, Z: 0},
			Color:     rl.DarkGray,
		},
		{
			Position:  rl.Vector3{X: 15, Y: 0, Z: 10},
			Size:      rl.Vector3{X: 2, Y: 2, Z: 2},
			Direction: rl.Vector3{X: 0, Y: 0, Z: 0},
			Color:     rl.DarkGray,
		},
	}

	for i := range initialStructures {
		m.StructureManager.Add(&initialStructures[i])
	}
}

func (m *Manager) StageWon() bool {
	return m.EnemyManager.AliveEnemyCount == 0
}

func (m *Manager) GameWon() bool {
	return m.EnemyManager.AliveEnemyCount == 0 && m.Difficulty >= m.GameEndingDifficulty
}

func (m *Manager) StageLost() bool {
	return !m.Player.IsAlive() && m.Player.ActionTimeLeft <= 0
}
