package blast

import rl "github.com/gen2brain/raylib-go/raylib"

var Radius float32 = 0.25
var BigRadius float32 = 0.25
var COLOR = rl.Yellow
var LIFETIME float32 = 1.00

type Blast struct {
	Position           rl.Vector3
	Radius             float32
	MaxLifeTime        float32
	LifeTime           float32
	Color              rl.Color
	AlwaysShow         bool
	IsHiddenFromKiller bool
}

func (b *Blast) Draw3D() {
	if !b.AlwaysShow && b.IsHiddenFromKiller {
		return
	}

	rl.DrawSphere(b.Position, b.Radius, b.Color)
}

func (b *Blast) Mutate(dt float32) {
	b.LifeTime -= dt

	lifeRatio := b.LifeTime / b.MaxLifeTime
	if lifeRatio < 0 {
		lifeRatio = 0
	}

	b.Radius = b.Radius * lifeRatio
}

func Create(position rl.Vector3, isByPlayer bool) Blast {
	return Blast{
		Position:    position,
		Radius:      Radius,
		MaxLifeTime: LIFETIME,
		LifeTime:    LIFETIME,
		Color:       COLOR,
		AlwaysShow:  isByPlayer,
	}
}

func CreateBig(position rl.Vector3, isByPlayer bool) Blast {
	return Blast{
		Position:    position,
		Radius:      BigRadius,
		MaxLifeTime: LIFETIME,
		LifeTime:    LIFETIME,
		Color:       COLOR,
		AlwaysShow:  isByPlayer,
	}
}

func CreateSplash(position rl.Vector3) Blast {
	return Blast{
		Position:    position,
		Radius:      0.55,
		MaxLifeTime: 0.25,
		LifeTime:    0.25,
		Color:       rl.NewColor(220, 80, 20, 230),
		AlwaysShow:  true,
	}
}

func CreateDebris(position rl.Vector3) Blast {
	return Blast{
		Position:    position,
		Radius:      0.2,
		MaxLifeTime: 0.35,
		LifeTime:    0.35,
		Color:       rl.NewColor(180, 40, 10, 200),
		AlwaysShow:  true,
	}
}

func createMuzzleFlash(position rl.Vector3) Blast {
	return Blast{
		Position:    position,
		Radius:      0.45,
		MaxLifeTime: 0.10,
		LifeTime:    0.10,
		Color:       rl.NewColor(255, 210, 60, 255),
		AlwaysShow:  true,
	}
}

func CreateMuzzleBlast(position rl.Vector3, direction rl.Vector3) []Blast {
	position.Z -= 0.5
	position = rl.Vector3Add(position, rl.Vector3Scale(direction, 1.5))
	blasts := []Blast{createMuzzleFlash(position)}
	right := rl.Vector3Normalize(rl.Vector3{X: -direction.Z, Y: 0, Z: direction.X})
	offsets := []rl.Vector3{
		rl.Vector3Scale(right, 0.25),
		rl.Vector3Scale(right, -0.25),
		rl.Vector3Scale(direction, 0.35),
	}
	for _, off := range offsets {
		blasts = append(blasts, Blast{
			Position:    rl.Vector3Add(position, off),
			Radius:      0.18,
			MaxLifeTime: 0.14,
			LifeTime:    0.14,
			Color:       rl.NewColor(255, 140, 30, 210),
			AlwaysShow:  true,
		})
	}
	return blasts
}
