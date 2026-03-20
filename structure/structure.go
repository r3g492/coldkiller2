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

func (s *Structure) getStructureCorners() []rl.Vector3 {
	halfX := float64(s.Size.X / 2.0)
	halfZ := float64(s.Size.Z / 2.0)

	angle := math.Atan2(float64(s.Direction.X), float64(s.Direction.Z))
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	localCorners := [4][2]float64{
		{-halfX, -halfZ},
		{halfX, -halfZ},
		{halfX, halfZ},
		{-halfX, halfZ},
	}

	worldCorners := make([]rl.Vector3, 4)
	for i := 0; i < 4; i++ {
		lx := localCorners[i][0]
		lz := localCorners[i][1]

		rotX := lx*cosA + lz*sinA
		rotZ := -lx*sinA + lz*cosA

		worldCorners[i] = rl.Vector3{
			X: s.Position.X + float32(rotX),
			Y: 0.0,
			Z: s.Position.Z + float32(rotZ),
		}
	}
	return worldCorners
}

func GetBoundaryRays(playerPos rl.Vector3, structures []*Structure) []rl.Ray {
	var rays []rl.Ray
	uniqueAngles := make(map[float64]bool)

	for i := 0; i < len(structures); i++ {
		corners := structures[i].getStructureCorners()

		for _, corner := range corners {
			dx := float64(corner.X - playerPos.X)
			dz := float64(corner.Z - playerPos.Z)
			angle := math.Atan2(dz, dx)

			uniqueAngles[angle-0.0001] = true
			uniqueAngles[angle] = true
			uniqueAngles[angle+0.0001] = true
		}
	}

	eyePos := playerPos
	eyePos.Y = 0.0

	for angle := range uniqueAngles {
		dirX := float32(math.Cos(angle))
		dirZ := float32(math.Sin(angle))

		dir := rl.Vector3{X: dirX, Y: 0.0, Z: dirZ}

		rays = append(rays, rl.Ray{
			Position:  eyePos,
			Direction: dir,
		})
	}

	return rays
}

func (s *Structure) GetStructureCorners() []rl.Vector3 {
	halfX := float64(s.Size.X / 2.0)
	halfZ := float64(s.Size.Z / 2.0)

	// FIX 1: Match the Atan2(X, Z) order used in Draw3D and RayCollisionOBB
	angle := math.Atan2(float64(s.Direction.X), float64(s.Direction.Z))
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	localCorners := [4][2]float64{
		{-halfX, -halfZ},
		{halfX, -halfZ},
		{halfX, halfZ},
		{-halfX, halfZ},
	}

	worldCorners := make([]rl.Vector3, 4)
	for i := 0; i < 4; i++ {
		lx := localCorners[i][0]
		lz := localCorners[i][1]

		// FIX 2: The exact mathematical inverse of your CheckCollision rotation logic
		rotX := lx*cosA + lz*sinA
		rotZ := -lx*sinA + lz*cosA

		worldCorners[i] = rl.Vector3{
			X: s.Position.X + float32(rotX),
			Y: 0.0,
			Z: s.Position.Z + float32(rotZ),
		}
	}
	return worldCorners
}
