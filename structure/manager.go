package structure

import rl "github.com/gen2brain/raylib-go/raylib"

type Manager struct {
	Structures []Structure
}

func CreateManager() *Manager {
	return &Manager{}
}

func (sm *Manager) Add(adding []Structure) {
	sm.Structures = append(sm.Structures, adding...)
}

func (sm *Manager) Unload() {
	sm.Structures = []Structure{}
}

func (sm *Manager) Draw3D() {
	for _, s := range sm.Structures {
		s.Draw3D()
	}
}

func (sm *Manager) CheckCollision(otherPos rl.Vector3, otherSize rl.Vector3) bool {
	for _, s := range sm.Structures {
		if s.CheckCollision(otherPos, otherSize) {
			return true
		}
	}
	return false
}

func (sm *Manager) Init() {
	sm.Add(
		[]Structure{
			{
				Position:  rl.Vector3{X: 5, Y: 0, Z: 5},
				Size:      rl.Vector3{X: 1, Y: 0, Z: 10},
				Direction: rl.Vector3{X: 1, Y: 1, Z: 1},
				Color:     rl.Green,
			},
		},
	)
}
