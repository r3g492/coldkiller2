package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/structure"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	Difficulty       int
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
	m.structureManager = structureManager
	m.enemyManager = enemyManager
}

func (m *Manager) CreateNewStage(pPos rl.Vector3) {
	// Difficulty에 따라 enemy와 structure 변주 추가
	m.enemyManager.AddEnemy(pPos)
}
