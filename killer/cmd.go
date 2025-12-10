package killer

import rl "github.com/gen2brain/raylib-go/raylib"

type BulletCmd struct {
	Pos rl.Vector3
	Dir rl.Vector3
}

type PushCmd struct {
	Position  rl.Vector3
	Direction rl.Vector3
	Radius    float32
	LifeTime  float32
	Force     float32
}
