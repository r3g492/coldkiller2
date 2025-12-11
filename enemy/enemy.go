package enemy

import (
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
	AttackSound           rl.Sound
	ActionTimeLeft        float32
	PushedTimeLeft        float32
	PushDirection         rl.Vector3
	PushForce             float32
	Health                int32
	IsDead                bool
}

func (e *Enemy) Draw3D() {
	anim := e.Animation[e.AnimationIdx]
	rl.UpdateModelAnimation(e.Model, anim, e.AnimationCurrentFrame)
	rl.PushMatrix()
	rl.Translatef(e.Position.X, e.Position.Y, e.Position.Z)
	rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, e.Size*2, e.Size*2, e.Size*2, rl.Red)
	rl.Rotatef(-60, 1, 0, 0)
	rl.Rotatef(e.ModelAngleDeg, 0, 1, 0)
	rl.DrawModel(e.Model, rl.NewVector3(0, -e.Size, 0), 0.7, rl.White)
	rl.PopMatrix()
	rl.DrawRay(rl.NewRay(e.Position, e.TargetDirection), rl.Green)
}

func (e *Enemy) Mutate(dt float32) ([]BulletCmd, []PushCmd) {
	if e.PushedTimeLeft > 0 {
		e.PushedTimeLeft -= dt
		moveAmount := rl.Vector3Scale(e.PushDirection, e.PushForce*dt)
		move := rl.Vector3Length(moveAmount) > 0
		if move {
			e.Position = rl.Vector3Add(e.Position, moveAmount)
		}
		return []BulletCmd{}, []PushCmd{}
	}

	var bulletCmds []BulletCmd
	var pushCmds []PushCmd
	if e.ActionTimeLeft > 0 {
		e.ActionTimeLeft -= dt
		return bulletCmds, pushCmds
	}
	return []BulletCmd{}, []PushCmd{}
}

func (e *Enemy) Damage(d int32) {
	e.Health -= d
	if e.Health <= 0 {
		e.IsDead = true
	}
}

func (e *Enemy) Push(
	pushDirection rl.Vector3,
	force float32,
) {
	e.PushedTimeLeft = 1.0
	e.PushDirection = pushDirection
	e.PushForce = force
	e.AnimationIdx = 1
	e.AnimationFrameSpeed = 58
	e.AnimationCurrentFrame = 0
}

func (e *Enemy) PlanAnimate(dt float32) {
	e.AnimationFrameCounter += e.AnimationFrameSpeed * dt
	anim := e.Animation[e.AnimationIdx]
	for e.AnimationFrameCounter >= 1.0 {
		e.AnimationCurrentFrame++
		e.AnimationFrameCounter -= 1.0

		if e.AnimationIdx == 1 && e.AnimationCurrentFrame >= anim.FrameCount-1 {
			e.AnimationCurrentFrame = anim.FrameCount - 1
			return
		}

		if e.AnimationCurrentFrame >= anim.FrameCount {
			e.AnimationCurrentFrame = 0
		}
	}
}

func (e *Enemy) Unload() {
	rl.UnloadModel(e.Model)
	rl.UnloadModelAnimations(e.Animation)
}
