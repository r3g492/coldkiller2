package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/structure"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Difficulty       int
	HighestBeaten    int
	StructureManager *structure.Manager
	EnemyManager     *enemy.Manager
	Player           *killer.Killer
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
