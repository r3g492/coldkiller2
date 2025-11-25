package killer

import rl "github.com/gen2brain/raylib-go/raylib"

type Killer struct {
	model rl.Model

	animation             []rl.ModelAnimation
	animationIdx          int
	animationCurrentFrame int32
	animationFrameCounter float32
	animationFrameSpeed   float32

	direction rl.Vector3
	Position  rl.Vector3
	Size      float32
}

func Init(
	model rl.Model,
	animation []rl.ModelAnimation,
	animationIdx int,
	animationCurrentFrame int32,
	animationFrameCounter float32,
	animationFrameSpeed float32,
	direction rl.Vector3,
	position rl.Vector3,
	size float32,
) *Killer {
	return &Killer{
		model:                 model,
		animation:             animation,
		animationIdx:          animationIdx,
		animationCurrentFrame: animationCurrentFrame,
		animationFrameCounter: animationFrameCounter,
		animationFrameSpeed:   animationFrameSpeed,
		direction:             direction,
		Position:              position,
		Size:                  size,
	}
}

func (k *Killer) Draw3D() {
	anim := k.animation[k.animationIdx]
	rl.UpdateModelAnimation(k.model, anim, k.animationCurrentFrame)

	rl.PushMatrix()
	rl.Translatef(k.Position.X, k.Position.Y, k.Position.Z)
	rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, k.Size*2, k.Size*2, k.Size*2, rl.Purple)

	// movement.RotateByDirection(s.Direction)
	rl.Rotatef(270, 1, 0, 0)
	rl.DrawModel(k.model, rl.NewVector3(0, -k.Size, 0), 0.7, rl.White)
	rl.PopMatrix()
}
