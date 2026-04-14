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
