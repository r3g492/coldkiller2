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

var cachedShadowTiles []rl.Vector3
var lastPlayerPos rl.Vector3
var forceShadowUpdate = true

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

func DrawShadowFloor(playerPos rl.Vector3, sm *structure.Manager) {
	tileSize := float32(0.3)
	gridRadius := 95

	if rl.Vector3Distance(lastPlayerPos, playerPos) > (tileSize / 2) {
		forceShadowUpdate = true
	}

	if forceShadowUpdate {
		lastPlayerPos = playerPos
		cachedShadowTiles = make([]rl.Vector3, 0)

		startX := int(playerPos.X/tileSize) - gridRadius
		endX := int(playerPos.X/tileSize) + gridRadius
		startZ := int(playerPos.Z/tileSize) - gridRadius
		endZ := int(playerPos.Z/tileSize) + gridRadius

		eyePos := playerPos
		eyePos.Y = 0.0

		for x := startX; x <= endX; x++ {
			for z := startZ; z <= endZ; z++ {
				tilePos := rl.Vector3{
					X: float32(x) * tileSize,
					Y: 0.0,
					Z: float32(z) * tileSize,
				}

				targetPos := tilePos
				targetPos.Y = 0.0
				isVisible := false

				if rl.Vector3Distance(eyePos, targetPos) < 1.0 {
					isVisible = true
				} else if isWithinDistance3D(eyePos, targetPos, MaxSightDistance) {
					if hasLineOfSight3D(eyePos, targetPos, sm) {
						isVisible = true
					}
				}

				if !isVisible {
					cachedShadowTiles = append(cachedShadowTiles, tilePos)
				}
			}
		}
		forceShadowUpdate = false
	}

	shadowColor := rl.NewColor(0, 0, 0, 255)
	for i := 0; i < len(cachedShadowTiles); i++ {
		drawPos := cachedShadowTiles[i]
		drawPos.Y = 0.05
		rl.DrawPlane(drawPos, rl.Vector2{X: tileSize + 0.1, Y: tileSize + 0.1}, shadowColor)
	}
}
