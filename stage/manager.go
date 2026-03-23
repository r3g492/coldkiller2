package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/structure"
)

type Manager struct {
	Difficulty       int
	StageWon         int
	structureManager *structure.SpatialManager
	enemyManager     *enemy.Manager
}

func CreateManager() *Manager {
	return &Manager{}
}

func (m *Manager) GenerateNewStage() {

}

func (m *Manager) Unload() {
}

func (m *Manager) Init(
	structureManager *structure.SpatialManager,
	enemyManager *enemy.Manager,
) {
	m.Difficulty = 0
	m.StageWon = 0
	m.structureManager = structureManager
	m.enemyManager = enemyManager
}

func (m *Manager) SetDifficulty(
	difficulty int,
) {
	m.Difficulty = difficulty
}

func (m *Manager) ResetScore() {
	m.StageWon = 0
}

func (m *Manager) ScoreUp() {
	m.StageWon++
}
