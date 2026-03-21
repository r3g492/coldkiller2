package structure

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const RADIUS = 100

type GridCoord struct {
	X, Y, Z int
}

type SpatialManager struct {
	CellSize float32
	Grid     map[GridCoord][]*Structure
}

func CreateSpatialManager(cellSize float32) *SpatialManager {
	return &SpatialManager{
		CellSize: cellSize,
		Grid:     make(map[GridCoord][]*Structure),
	}
}

func (sm *SpatialManager) GetCoord(pos rl.Vector3) GridCoord {
	return GridCoord{
		X: int(math.Floor(float64(pos.X / sm.CellSize))),
		Y: int(math.Floor(float64(pos.Y / sm.CellSize))),
		Z: int(math.Floor(float64(pos.Z / sm.CellSize))),
	}
}

func (sm *SpatialManager) Add(s *Structure) {
	coord := sm.GetCoord(s.Position)
	sm.Grid[coord] = append(sm.Grid[coord], s)
}

func (sm *SpatialManager) GetStructuresNearPosition(pos rl.Vector3, searchRadius float32) []*Structure {
	var nearby []*Structure
	centerCoord := sm.GetCoord(pos)
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				checkCoord := GridCoord{
					X: centerCoord.X + x,
					Y: centerCoord.Y + y,
					Z: centerCoord.Z + z,
				}
				if structuresInCell, exists := sm.Grid[checkCoord]; exists {
					for _, s := range structuresInCell {
						dist := rl.Vector3Distance(pos, s.Position)
						if dist <= searchRadius {
							nearby = append(nearby, s)
						}
					}
				}
			}
		}
	}

	return nearby
}

func (sm *SpatialManager) CheckCollision(otherPos rl.Vector3, prevPos rl.Vector3, otherSize rl.Vector3) bool {
	for _, s := range sm.GetStructuresNearPosition(otherPos, RADIUS) {
		if s.CheckCollision(otherPos, prevPos, otherSize) {
			return true
		}
	}
	return false
}

func (sm *SpatialManager) Unload() {
	clear(sm.Grid)
}

func (sm *SpatialManager) Init() {
	initialStructures := []Structure{
		{
			Position:  rl.Vector3{X: 5, Y: 0, Z: 5},
			Size:      rl.Vector3{X: 1, Y: 1, Z: 10},
			Direction: rl.Vector3{X: 1, Y: 0.2, Z: 1},
			Color:     rl.DarkGray,
		},
		{
			Position:  rl.Vector3{X: -5, Y: 0, Z: -5},
			Size:      rl.Vector3{X: 1, Y: 1, Z: 10},
			Direction: rl.Vector3{X: 0, Y: 0, Z: 0},
			Color:     rl.DarkGray,
		},
		{
			Position:  rl.Vector3{X: 10, Y: 0, Z: 0},
			Size:      rl.Vector3{X: 1, Y: 1, Z: 1},
			Direction: rl.Vector3{X: 0, Y: 0, Z: 0},
			Color:     rl.DarkGray,
		},
	}

	for i := range initialStructures {
		sm.Add(&initialStructures[i])
	}
}

func (sm *SpatialManager) Draw3D(playerPosition rl.Vector3) {
	for _, s := range sm.GetStructuresNearPosition(playerPosition, RADIUS) {
		s.Draw3D()
	}
}
