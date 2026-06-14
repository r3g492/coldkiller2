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
	Charger
)

// aiContext holds values derived once per tick and shared across AI behaviors.
type aiContext struct {
	distToPlayer  float32
	dirToPlayer   rl.Vector3
	aimObstructed bool
	shouldAim     bool
}

func newAiContext(e *Enemy, p *killer.Killer, sm *structure.Manager) aiContext {
	distToPlayer := rl.Vector3Distance(e.Position, p.Position)
	dirToPlayer := rl.Vector3Normalize(rl.Vector3Subtract(p.Position, e.Position))
	aimObstructed := sm.RayObstructed(e.Position, p.Position)
	shouldAim := e.AimTimeLeft > 0 && distToPlayer <= e.AttackRange && !aimObstructed
	return aiContext{
		distToPlayer:  distToPlayer,
		dirToPlayer:   dirToPlayer,
		aimObstructed: aimObstructed,
		shouldAim:     shouldAim,
	}
}

func deriveAi(
	e *Enemy,
	em *Manager,
	myIdx int,
	p *killer.Killer,
	structureManager *structure.Manager,
) (bool, rl.Vector3) {
	if !e.IsAlive() {
		return false, rl.Vector3{}
	}

	ctx := newAiContext(e, p, structureManager)

	switch e.AiType {
	case SimpleZombie:
		return deriveSimpleZombie(e, ctx)
	case Charger:
		return deriveCharger(e, ctx, structureManager)
	case Elite:
		return deriveElite(e, myIdx, ctx, structureManager)
	}
	return false, rl.Vector3{}
}

func deriveSimpleZombie(e *Enemy, ctx aiContext) (bool, rl.Vector3) {
	shouldAim := e.AimTimeLeft > 0 && ctx.distToPlayer <= e.AttackRange
	return shouldAim, ctx.dirToPlayer
}

func deriveCharger(e *Enemy, ctx aiContext, sm *structure.Manager) (bool, rl.Vector3) {
	moveDir := steerAroundObstacles(e, sm, ctx.dirToPlayer)
	return ctx.shouldAim, moveDir
}

func deriveElite(e *Enemy, myIdx int, ctx aiContext, sm *structure.Manager) (bool, rl.Vector3) {
	optimalRange := e.AttackRange * 0.8
	tooCloseRange := e.AttackRange * 0.4

	var moveDir rl.Vector3
	switch {
	case ctx.distToPlayer > optimalRange || ctx.aimObstructed:
		moveDir = ctx.dirToPlayer
	case ctx.distToPlayer < tooCloseRange:
		moveDir = rl.Vector3Scale(ctx.dirToPlayer, -1)
	default:
		up := rl.Vector3{X: 0, Y: 1, Z: 0}
		moveDir = rl.Vector3CrossProduct(ctx.dirToPlayer, up)
		if myIdx%2 == 0 {
			moveDir = rl.Vector3Scale(moveDir, -1)
		}
	}

	moveDir = steerAroundObstacles(e, sm, moveDir)
	return ctx.shouldAim, moveDir
}

func steerAroundObstacles(e *Enemy, sm *structure.Manager, moveDir rl.Vector3) rl.Vector3 {
	lookAheadDist := e.Size * 3.0
	collisionSize := rl.Vector3{X: e.Size, Y: e.Size, Z: e.Size}

	blocked := func(dir rl.Vector3) bool {
		probe := rl.Vector3Add(e.Position, rl.Vector3Scale(dir, lookAheadDist))
		return sm.CheckCollision(probe, e.Position, collisionSize)
	}

	if blocked(moveDir) {
		dirRight := rotateY(moveDir, 45)
		dirLeft := rotateY(moveDir, -45)
		if !blocked(dirRight) {
			moveDir = dirRight
		} else if !blocked(dirLeft) {
			moveDir = dirLeft
		} else {
			moveDir = rotateY(moveDir, 90)
		}
	}

	if rl.Vector3Length(moveDir) > 0 {
		moveDir = rl.Vector3Normalize(moveDir)
	}
	return moveDir
}
