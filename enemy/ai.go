package enemy

import (
	"coldkiller2/killer"
	"coldkiller2/structure"
	"math"

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
	if !e.IsAlive() {
		return false, rl.Vector3{}
	}

	distToPlayer := rl.Vector3Distance(e.Position, p.Position)
	dirToPlayer := rl.Vector3Normalize(rl.Vector3Subtract(p.Position, e.Position))

	if e.AiType == SimpleZombie {
		shouldAim := e.AimTimeLeft > 0 && distToPlayer <= e.AttackRange
		moveDir := rl.Vector3Normalize(rl.Vector3Subtract(p.Position, e.Position))
		return shouldAim, moveDir
	}
	if e.AiType == Elite {
		shouldAim := e.AimTimeLeft > 0 && distToPlayer <= e.AttackRange
		var moveDir rl.Vector3
		optimalRange := e.AttackRange * 0.8
		tooCloseRange := e.AttackRange * 0.4

		if distToPlayer > optimalRange {
			moveDir = dirToPlayer
		} else if distToPlayer < tooCloseRange {
			moveDir = rl.Vector3Scale(dirToPlayer, -1)
		} else {
			up := rl.Vector3{X: 0, Y: 1, Z: 0}
			moveDir = rl.Vector3CrossProduct(dirToPlayer, up)

			if myIdx%2 == 0 {
				moveDir = rl.Vector3Scale(moveDir, -1)
			}
		}

		lookAheadDist := e.Size * 3.0
		collisionSize := rl.Vector3{X: e.Size, Y: e.Size, Z: e.Size}

		rotateY := func(v rl.Vector3, angleRad float64) rl.Vector3 {
			cosA := float32(math.Cos(angleRad))
			sinA := float32(math.Sin(angleRad))
			return rl.Vector3{
				X: v.X*cosA - v.Z*sinA,
				Y: v.Y,
				Z: v.X*sinA + v.Z*cosA,
			}
		}

		probePos := rl.Vector3Add(e.Position, rl.Vector3Scale(moveDir, lookAheadDist))

		if structureManager.CheckCollision(probePos, e.Position, collisionSize) {
			dirRight := rotateY(moveDir, math.Pi/4)
			probeRight := rl.Vector3Add(e.Position, rl.Vector3Scale(dirRight, lookAheadDist))

			dirLeft := rotateY(moveDir, -math.Pi/4)
			probeLeft := rl.Vector3Add(e.Position, rl.Vector3Scale(dirLeft, lookAheadDist))

			if !structureManager.CheckCollision(probeRight, e.Position, collisionSize) {
				moveDir = dirRight
			} else if !structureManager.CheckCollision(probeLeft, e.Position, collisionSize) {
				moveDir = dirLeft
			} else {
				moveDir = rotateY(moveDir, math.Pi/2)
			}
		}

		if rl.Vector3Length(moveDir) > 0 {
			moveDir = rl.Vector3Normalize(moveDir)
		}

		return shouldAim, moveDir
	}
	// TODO: 다른 ai type 추가
	return false, rl.Vector3{}
}
