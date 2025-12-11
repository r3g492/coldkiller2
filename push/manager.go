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
			LifeTime:  pc.LifeTime,
			Shooter:   Player,
			Color:     rl.Green,
			Active:    true,
			Force:     pc.Force,
		}
		pm.PushList = append(pm.PushList, p)
	}
}

func (pm *Manager) EnemyPushCreate(
	cmds []enemy.PushCmd,
) {
	for _, ec := range cmds {
		p := Push{
			Position:  ec.Position,
			Direction: ec.Direction,
			Radius:    ec.Radius,
			LifeTime:  ec.LifeTime,
			Shooter:   Enemy,
			Color:     rl.Red,
			Active:    true,
			Force:     ec.Force,
		}
		pm.PushList = append(pm.PushList, p)
	}
}

func (pm *Manager) Mutate(dt float32, p *killer.Killer, enemyList []enemy.Enemy) {
	for i := 0; i < len(pm.PushList); i++ {
		pm.PushList[i].Mutate(dt)
		for j := 0; j < len(enemyList); j++ {
			enemyPos := enemyList[j].Position
			curPush := pm.PushList[i]
			if rl.Vector3Distance(enemyPos, curPush.Position) < curPush.Radius {
				enemyList[j].Push(
					curPush.Direction,
					curPush.Force,
				)
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
