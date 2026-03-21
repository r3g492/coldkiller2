package sight

import (
	"coldkiller2/blast"
	"coldkiller2/bullet"
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/structure"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MaxSightDistance float32 = 800.0

func UpdateSight(
	blastManager *blast.Manager,
	bulletManager *bullet.Manager,
	enemyManager *enemy.Manager,
	structureManager *structure.SpatialManager,
	player *killer.Killer,
) {
	for i := 0; i < len(blastManager.Blasts); i++ {
		blastManager.Blasts[i].IsHiddenFromKiller = true
	}
	for i := 0; i < len(bulletManager.Bullets); i++ {
		bulletManager.Bullets[i].IsHiddenFromKiller = true
	}
	for i := 0; i < len(enemyManager.Enemies); i++ {
		enemyManager.Enemies[i].IsHiddenFromKiller = true
	}
	var structureList = structureManager.GetStructuresNearPosition(player.Position, structure.RADIUS)
	for i := 0; i < len(structureList); i++ {
		structureList[i].IsHiddenFromKiller = true
	}

	pPos := player.Position
	for i := 0; i < len(structureList); i++ {
		if isWithinDistance3D(pPos, structureList[i].Position, MaxSightDistance) {
			structureList[i].IsHiddenFromKiller = false
		}
	}

	for i := 0; i < len(enemyManager.Enemies); i++ {
		e := &enemyManager.Enemies[i]
		if isWithinDistance3D(pPos, e.Position, MaxSightDistance) {
			if hasLineOfSight3D(pPos, e.Position, structureManager) {
				e.IsHiddenFromKiller = false
			}
		}
	}

	for i := 0; i < len(bulletManager.Bullets); i++ {
		b := &bulletManager.Bullets[i]
		if isWithinDistance3D(pPos, b.Position, MaxSightDistance) {
			if hasLineOfSight3D(pPos, b.Position, structureManager) {
				b.IsHiddenFromKiller = false
			}
		}
	}

	for i := 0; i < len(blastManager.Blasts); i++ {
		b := &blastManager.Blasts[i]
		if isWithinDistance3D(pPos, b.Position, MaxSightDistance) {
			if hasLineOfSight3D(pPos, b.Position, structureManager) {
				b.IsHiddenFromKiller = false
			}
		}
	}
}

func isWithinDistance3D(p1, p2 rl.Vector3, maxDist float32) bool {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return (dx*dx + dy*dy + dz*dz) <= (maxDist * maxDist)
}

func hasLineOfSight3D(start, end rl.Vector3, sm *structure.SpatialManager) bool {
	direction := rl.Vector3Subtract(end, start)

	distanceToTarget := rl.Vector3Length(direction)

	direction = rl.Vector3Normalize(direction)

	ray := rl.Ray{
		Position:  start,
		Direction: direction,
	}

	var structureList = sm.GetStructuresNearPosition(start, structure.RADIUS)
	for i := 0; i < len(structureList); i++ {
		s := structureList[i]

		hitInfo := s.RayCollisionOBB(ray)

		if hitInfo.Hit {
			if hitInfo.Distance < distanceToTarget {
				return false
			}
		}
	}

	return true
}

func DrawSolidShadows(playerPos rl.Vector3, sm *structure.SpatialManager) {
	eyePos := playerPos
	eyePos.Y = 0.0

	shadowColor := rl.NewColor(0, 0, 0, 255) // Solid Black

	structures := sm.GetStructuresNearPosition(eyePos, structure.RADIUS)

	for _, s := range structures {
		corners := s.GetStructureCorners()

		for i := 0; i < 4; i++ {
			A := corners[i]
			B := corners[(i+1)%4]

			dx := B.X - A.X
			dz := B.Z - A.Z

			nx := dz
			nz := -dx

			cx := (A.X+B.X)/2.0 - eyePos.X
			cz := (A.Z+B.Z)/2.0 - eyePos.Z

			dot := cx*nx + cz*nz

			if dot > 0 {

				dirA := rl.Vector3Normalize(rl.Vector3Subtract(A, eyePos))
				dirB := rl.Vector3Normalize(rl.Vector3Subtract(B, eyePos))

				projA := rl.Vector3{
					X: A.X + dirA.X*MaxSightDistance,
					Y: 0.01, // Slightly raised so it doesn't z-fight/glitch with the floor
					Z: A.Z + dirA.Z*MaxSightDistance,
				}
				projB := rl.Vector3{
					X: B.X + dirB.X*MaxSightDistance,
					Y: 0.01,
					Z: B.Z + dirB.Z*MaxSightDistance,
				}

				A.Y = 0.01
				B.Y = 0.01

				rl.DrawTriangle3D(A, B, projA, shadowColor)
				rl.DrawTriangle3D(B, projB, projA, shadowColor)

				rl.DrawTriangle3D(A, projA, B, shadowColor)
				rl.DrawTriangle3D(B, projA, projB, shadowColor)
			}
		}
	}
}
