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
		m.structureManager.Add(&initialStructures[i])
	}
}
