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

func DrawBoundaryRayToStructures(playerPos rl.Vector3, sm *structure.SpatialManager) {
	eyePos := playerPos
	eyePos.Y = 0.0

	structures := sm.GetStructuresNearPosition(eyePos, structure.RADIUS)
	rays := structure.GetBoundaryRays(eyePos, structures)

	// Changed to red to represent the "shadow" or "blocked" path
	shadowRayColor := rl.NewColor(255, 0, 0, 150)

	for i := 0; i < len(rays); i++ {
		ray := rays[i]
		closestHitDist := MaxSightDistance

		for j := 0; j < len(structures); j++ {
			hitInfo := structures[j].RayCollisionOBB(ray)
			if hitInfo.Hit && hitInfo.Distance < closestHitDist {
				closestHitDist = hitInfo.Distance
			}
		}

		// If the ray actually hit a structure before reaching the max distance...
		if closestHitDist < MaxSightDistance {
			// 1. Calculate where the ray hit the wall (The New Start Point)
			hitPos := rl.Vector3{
				X: ray.Position.X + ray.Direction.X*closestHitDist,
				Y: 0.0,
				Z: ray.Position.Z + ray.Direction.Z*closestHitDist,
			}

			// 2. Calculate the absolute edge of the vision range (The New End Point)
			maxPos := rl.Vector3{
				X: ray.Position.X + ray.Direction.X*MaxSightDistance,
				Y: 0.0,
				Z: ray.Position.Z + ray.Direction.Z*MaxSightDistance,
			}

			// 3. Draw the line extending outward from the back of the structure
			rl.DrawLine3D(hitPos, maxPos, shadowRayColor)
		}
	}
}
