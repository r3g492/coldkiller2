package blast

type Manager struct {
	Blasts []Blast
}

func CreateManager() *Manager {
	return &Manager{}
}

func (bm *Manager) AddBlasts(addingBlasts []Blast) {
	bm.Blasts = append(bm.Blasts, addingBlasts...)
}

func (bm *Manager) Unload() {
	bm.Blasts = []Blast{}
}

func (bm *Manager) Mutate(dt float32) {
	for i := 0; i < len(bm.Blasts); i++ {
		bm.Blasts[i].Mutate(dt)
		if bm.Blasts[i].LifeTime <= 0 {
			bm.Blasts[i] = bm.Blasts[len(bm.Blasts)-1]
			bm.Blasts = bm.Blasts[:len(bm.Blasts)-1]
			i--
		}
	}
}

func (bm *Manager) DrawBlasts3D() {
	for _, b := range bm.Blasts {
		b.Draw()
	}
}
