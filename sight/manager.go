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
	structureManager *structure.Manager,
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
	for i := 0; i < len(structureManager.Structures); i++ {
		structureManager.Structures[i].IsHiddenFromKiller = true
	}

	pPos := player.Position
	for i := 0; i < len(structureManager.Structures); i++ {
		s := &structureManager.Structures[i]
		if isWithinDistance3D(pPos, s.Position, MaxSightDistance) {
			s.IsHiddenFromKiller = false
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

func hasLineOfSight3D(start, end rl.Vector3, sm *structure.Manager) bool {
	direction := rl.Vector3Subtract(end, start)

	distanceToTarget := rl.Vector3Length(direction)

	direction = rl.Vector3Normalize(direction)

	ray := rl.Ray{
		Position:  start,
		Direction: direction,
	}

	for i := 0; i < len(sm.Structures); i++ {
		s := &sm.Structures[i]

		hitInfo := s.RayCollisionOBB(ray)

		if hitInfo.Hit {
			if hitInfo.Distance < distanceToTarget {
				return false
			}
		}
	}

	return true
}
