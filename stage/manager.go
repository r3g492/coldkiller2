package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/structure"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Difficulty       int
	StageWon         int
	structureManager *structure.Manager
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
	structureManager *structure.Manager,
	enemyManager *enemy.Manager,
) {
	m.Difficulty = 0
	m.StageWon = 0
	m.structureManager = structureManager
	m.enemyManager = enemyManager
}

func (m *Manager) CreateNewStage(pPos rl.Vector3) {
	m.enemyManager.AddEnemy(pPos)
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
