package structure

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
