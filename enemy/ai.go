package enemy

import (
	"coldkiller2/killer"
	"coldkiller2/structure"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AiType int

const (
	SimpleZombie AiType = iota
	Elite
	Charger
)

type MoveMode int

const (
	ModeDirect MoveMode = iota
	ModeZigzag
	ModeStrafe
)

const (
	moveModeMinTime = 1.2
	moveModeMaxTime = 2.8
)

const (
	zombieZigzagDeg  = 35.0
	zombieZigzagFreq = 4.0

	chargerZigzagDeg  = 18.0
	chargerZigzagFreq = 7.0
)

type aiContext struct {
	distToPlayer  float32
	dirToPlayer   rl.Vector3
	aimObstructed bool
	shouldAim     bool
	dt            float32
	myIdx         int
}

func newAiContext(e *Enemy, p *killer.Killer, sm *structure.Manager, myIdx int, dt float32) aiContext {
	distToPlayer := rl.Vector3Distance(e.Position, p.Position)
	dirToPlayer := rl.Vector3Normalize(rl.Vector3Subtract(p.Position, e.Position))
	aimObstructed := sm.RayObstructed(e.Position, p.Position)
	shouldAim := e.AimTimeLeft > 0 && distToPlayer <= e.AttackRange && !aimObstructed
	return aiContext{
		distToPlayer:  distToPlayer,
		dirToPlayer:   dirToPlayer,
		aimObstructed: aimObstructed,
		shouldAim:     shouldAim,
		dt:            dt,
		myIdx:         myIdx,
	}
}

func deriveAi(
	e *Enemy,
	em *Manager,
	myIdx int,
	p *killer.Killer,
	structureManager *structure.Manager,
	dt float32,
) (bool, rl.Vector3) {
	if !e.IsAlive() {
		return false, rl.Vector3{}
	}

	ctx := newAiContext(e, p, structureManager, myIdx, dt)

	switch e.AiType {
	case SimpleZombie:
		return deriveSimpleZombie(e, ctx, structureManager)
	case Charger:
		return deriveCharger(e, ctx, structureManager)
	case Elite:
		return deriveElite(e, ctx, structureManager)
	}
	return false, rl.Vector3{}
}

func deriveSimpleZombie(e *Enemy, ctx aiContext, sm *structure.Manager) (bool, rl.Vector3) {
	shouldAim := e.AimTimeLeft > 0 && ctx.distToPlayer <= e.AttackRange
	moveDir := approachDir(e, ctx, zombieZigzagDeg, zombieZigzagFreq, true)
	moveDir = steerAroundObstacles(e, sm, moveDir)
	return shouldAim, moveDir
}

func deriveCharger(e *Enemy, ctx aiContext, sm *structure.Manager) (bool, rl.Vector3) {
	moveDir := approachDir(e, ctx, chargerZigzagDeg, chargerZigzagFreq, false)
	moveDir = steerAroundObstacles(e, sm, moveDir)
	return ctx.shouldAim, moveDir
}

func deriveElite(e *Enemy, ctx aiContext, sm *structure.Manager) (bool, rl.Vector3) {
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
		if ctx.myIdx%2 == 0 {
			moveDir = rl.Vector3Scale(moveDir, -1)
		}
	}

	moveDir = steerAroundObstacles(e, sm, moveDir)
	return ctx.shouldAim, moveDir
}

func approachDir(e *Enemy, ctx aiContext, zigDeg, zigFreq float32, allowStrafe bool) rl.Vector3 {
	updateMoveMode(e, ctx.dt, allowStrafe)

	switch e.MoveMode {
	case ModeZigzag:
		return zigzag(e, ctx, ctx.dirToPlayer, zigDeg, zigFreq)
	case ModeStrafe:
		return strafeDir(e, ctx)
	default:
		return ctx.dirToPlayer
	}
}

func updateMoveMode(e *Enemy, dt float32, allowStrafe bool) {
	e.MoveModeTimer -= dt
	if e.MoveModeTimer > 0 {
		return
	}
	e.MoveMode = pickMoveMode(allowStrafe)
	e.MoveModeTimer = moveModeMinTime + rand.Float32()*(moveModeMaxTime-moveModeMinTime)
	if e.MoveMode == ModeStrafe {
		e.StrafeSign = 1
		if rand.Intn(2) == 0 {
			e.StrafeSign = -1
		}
	}
}

func pickMoveMode(allowStrafe bool) MoveMode {
	r := rand.Float32()
	if !allowStrafe {
		if r < 0.5 {
			return ModeDirect
		}
		return ModeZigzag
	}
	switch {
	case r < 0.45:
		return ModeDirect
	case r < 0.8:
		return ModeZigzag
	default:
		return ModeStrafe
	}
}

func zigzag(e *Enemy, ctx aiContext, dir rl.Vector3, amplitudeDeg, freq float32) rl.Vector3 {
	e.WanderPhase += freq * ctx.dt
	angle := float64(amplitudeDeg) * math.Sin(float64(e.WanderPhase)+float64(ctx.myIdx))
	return rotateY(dir, angle)
}

func strafeDir(e *Enemy, ctx aiContext) rl.Vector3 {
	up := rl.Vector3{X: 0, Y: 1, Z: 0}
	perp := rl.Vector3Scale(rl.Vector3CrossProduct(ctx.dirToPlayer, up), e.StrafeSign)
	dir := rl.Vector3Add(rl.Vector3Scale(perp, 0.85), rl.Vector3Scale(ctx.dirToPlayer, 0.35))
	if rl.Vector3Length(dir) > 0 {
		dir = rl.Vector3Normalize(dir)
	}
	return dir
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
