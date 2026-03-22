package enemy

import (
	"coldkiller2/killer"
	"coldkiller2/structure"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AiType int

const (
	SimpleZombie AiType = iota
	Elite
)

func deriveAi(
	e *Enemy,
	em *Manager,
	myIdx int,
	p *killer.Killer,
	structureManager *structure.SpatialManager,
) (bool, rl.Vector3) {
	if e.AiType == SimpleZombie {
		return e.AimTimeLeft > 0 && rl.Vector3Distance(e.Position, p.Position) <= e.AttackRange && e.IsAlive(),
			rl.Vector3Normalize(
				rl.Vector3Subtract(
					p.Position,
					e.Position,
				),
			)
	}
	// TODO: 다른 ai type 추가
	return false, rl.Vector3{}
}
