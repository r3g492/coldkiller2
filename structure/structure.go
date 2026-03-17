package structure

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Structure struct {
	Position  rl.Vector3
	Size      rl.Vector3
	Direction rl.Vector3
	Color     rl.Color

	IsHiddenFromKiller bool
}

func (s *Structure) Draw3D() {
	if s.IsHiddenFromKiller {
		return
	}

	rl.PushMatrix()
	rl.Translatef(s.Position.X, s.Position.Y, s.Position.Z)

	angle := float32(math.Atan2(float64(s.Direction.X), float64(s.Direction.Z))) * rl.Rad2deg
	rl.Rotatef(angle, 0, 1, 0)

	rl.DrawCube(rl.Vector3{}, s.Size.X, s.Size.Y, s.Size.Z, s.Color)
	rl.DrawCubeWires(rl.Vector3{}, s.Size.X, s.Size.Y, s.Size.Z, rl.Black)
	rl.PopMatrix()
}

func (s *Structure) CheckCollision(otherPos rl.Vector3, otherSize rl.Vector3) bool {
	angleRad := math.Atan2(float64(s.Direction.X), float64(s.Direction.Z))

	relX := float64(otherPos.X - s.Position.X)
	relY := float64(otherPos.Y - s.Position.Y)
	relZ := float64(otherPos.Z - s.Position.Z)

	cosA := math.Cos(-angleRad)
	sinA := math.Sin(-angleRad)

	localX := relX*cosA + relZ*sinA
	localZ := -relX*sinA + relZ*cosA
	localY := relY

	limitX := float64((s.Size.X + otherSize.X) / 2)
	limitY := float64((s.Size.Y + otherSize.Y) / 2)
	limitZ := float64((s.Size.Z + otherSize.Z) / 2)

	return math.Abs(localX) <= limitX &&
		math.Abs(localY) <= limitY &&
		math.Abs(localZ) <= limitZ
}

func (s *Structure) RayCollisionOBB(ray rl.Ray) rl.RayCollision {
	angleRad := math.Atan2(float64(s.Direction.X), float64(s.Direction.Z))

	cosA := float32(math.Cos(-angleRad))
	sinA := float32(math.Sin(-angleRad))

	relPosX := ray.Position.X - s.Position.X
	relPosY := ray.Position.Y - s.Position.Y
	relPosZ := ray.Position.Z - s.Position.Z

	localRayPos := rl.Vector3{
		X: relPosX*cosA + relPosZ*sinA,
		Y: relPosY,
		Z: -relPosX*sinA + relPosZ*cosA,
	}

	localRayDir := rl.Vector3{
		X: ray.Direction.X*cosA + ray.Direction.Z*sinA,
		Y: ray.Direction.Y,
		Z: -ray.Direction.X*sinA + ray.Direction.Z*cosA,
	}

	localRay := rl.Ray{Position: localRayPos, Direction: localRayDir}

	localBox := rl.BoundingBox{
		Min: rl.Vector3{X: -s.Size.X / 2, Y: -s.Size.Y / 2, Z: -s.Size.Z / 2},
		Max: rl.Vector3{X: s.Size.X / 2, Y: s.Size.Y / 2, Z: s.Size.Z / 2},
	}

	return rl.GetRayCollisionBox(localRay, localBox)
}
