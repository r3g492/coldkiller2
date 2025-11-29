package enemy

import (
	"coldkiller2/bullet"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Enemy struct {
	Model                 rl.Model
	ModelAngleDeg         float32
	Animation             []rl.ModelAnimation
	AnimationIdx          int
	AnimationCurrentFrame int32
	AnimationFrameCounter float32
	AnimationFrameSpeed   float32
	MoveDirection         rl.Vector3
	TargetDirection       rl.Vector3
	Position              rl.Vector3
	Size                  float32
	MoveSpeed             float32
	Camera                rl.Camera3D
	ShotGunSound          rl.Sound
	ActionTimeLeft        float32
	State                 State
	Bullets               []bullet.Bullet
}

type State int

const (
	StateIdle   State = iota // 0
	StateMove                // 1
	StateAttack              // 2: Stationary shooting
	StateDash                // 3: Fast uncontrolled movement
	StateHit                 // 4: Stunned/Hurt
)
