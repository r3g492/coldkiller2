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

func (s *Structure) CheckCollision(curPos rl.Vector3, prevPos rl.Vector3, otherSize rl.Vector3) bool {
	// 1. Calculate the local orientation of the structure
	angleRad := math.Atan2(float64(s.Direction.X), float64(s.Direction.Z))
	cosA := math.Cos(-angleRad)
	sinA := math.Sin(-angleRad)

	// Helper function to transform a global position into the structure's local space
	toLocal := func(pos rl.Vector3) (float64, float64, float64) {
		relX := float64(pos.X - s.Position.X)
		relY := float64(pos.Y - s.Position.Y)
		relZ := float64(pos.Z - s.Position.Z)

		localX := relX*cosA + relZ*sinA
		localZ := -relX*sinA + relZ*cosA
		localY := relY
		return localX, localY, localZ
	}

	// 2. Transform both previous and current positions
	prevX, prevY, prevZ := toLocal(prevPos)
	curX, curY, curZ := toLocal(curPos)

	// 3. Define the bounding limits (Minkowski sum of both half-sizes)
	limitX := float64((s.Size.X + otherSize.X) / 2)
	limitY := float64((s.Size.Y + otherSize.Y) / 2)
	limitZ := float64((s.Size.Z + otherSize.Z) / 2)

	// 4. Perform Line-Segment vs AABB intersection (Slab Method)
	// tMin and tMax represent the segment from prevPos (t=0.0) to curPos (t=1.0)
	tMin := 0.0
	tMax := 1.0

	// X-Axis check
	dx := curX - prevX
	if math.Abs(dx) < 1e-8 { // Parallel to the plane
		if prevX < -limitX || prevX > limitX {
			return false
		}
	} else {
		t1 := (-limitX - prevX) / dx
		t2 := (limitX - prevX) / dx
		if t1 > t2 {
			t1, t2 = t2, t1 // Swap so t1 is always the smaller intersection
		}
		if t1 > tMin {
			tMin = t1
		}
		if t2 < tMax {
			tMax = t2
		}
		if tMin > tMax {
			return false
		} // Ray misses the box
	}

	// Y-Axis check
	dy := curY - prevY
	if math.Abs(dy) < 1e-8 {
		if prevY < -limitY || prevY > limitY {
			return false
		}
	} else {
		t1 := (-limitY - prevY) / dy
		t2 := (limitY - prevY) / dy
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		if t1 > tMin {
			tMin = t1
		}
		if t2 < tMax {
			tMax = t2
		}
		if tMin > tMax {
			return false
		}
	}

	// Z-Axis check
	dz := curZ - prevZ
	if math.Abs(dz) < 1e-8 {
		if prevZ < -limitZ || prevZ > limitZ {
			return false
		}
	} else {
		t1 := (-limitZ - prevZ) / dz
		t2 := (limitZ - prevZ) / dz
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		if t1 > tMin {
			tMin = t1
		}
		if t2 < tMax {
			tMax = t2
		}
		if tMin > tMax {
			return false
		}
	}

	// If we make it here, the segment [0, 1] overlaps the bounding box limits
	return true
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
