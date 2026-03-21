package enemy

import (
	"coldkiller2/killer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AiType int

const (
	SimpleZombie AiType = iota
	Elite
)

func deriveAimCondition(
	e *Enemy,
	distToPlayer float32,
) bool {
	if e.AiType == SimpleZombie {
		return e.AimTimeLeft > 0 && distToPlayer <= e.AttackRange && e.IsAlive()
	}
	// TODO: 다른 ai type 추가
	return false
}

func deriveMovementDirection(
	e *Enemy,
	p *killer.Killer,
) rl.Vector3 {
	if e.AiType == SimpleZombie {
		return rl.Vector3Normalize(
			rl.Vector3Subtract(
				p.Position,
				e.Position,
			),
		)
	}
	return rl.Vector3{}
}
