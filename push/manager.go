package push

import (
	"coldkiller2/enemy"
	"coldkiller2/killer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Manager struct {
	PushList []Push
}

func CreateManager() *Manager {
	return &Manager{}
}

func (pm *Manager) KillerPushCreate(
	cmds []killer.PushCmd,
) {
	for _, pc := range cmds {
		p := Push{
			Position:  pc.Position,
			Direction: pc.Direction,
			Radius:    pc.Radius,
			LifeTime:  1.0,
			Shooter:   Player,
			Color:     rl.Green,
			Active:    true,
			Force:     pc.Force,
		}
		pm.PushList = append(pm.PushList, p)
	}
}

func (pm *Manager) Mutate(dt float32, p *killer.Killer, el []enemy.Enemy) {
	for i := 0; i < len(pm.PushList); i++ {
		pm.PushList[i].Mutate(dt)
		for j := 0; j < len(el); j++ {
			enemyPos := el[j].Position
			enemySize := el[j].Size
			curPush := pm.PushList[i]
			if rl.Vector3Distance(enemyPos, curPush.Position) < enemySize {
				el[j].Damage(50)
				pm.PushList[i].Active = false
			}
		}

		if pm.PushList[i].LifeTime <= 0 || !pm.PushList[i].Active {
			pm.PushList[i] = pm.PushList[len(pm.PushList)-1]
			pm.PushList = pm.PushList[:len(pm.PushList)-1]
			i--
		}
	}
}

func (pm *Manager) DrawPush3D() {
	for _, p := range pm.PushList {
		p.DrawPush()
	}
}
