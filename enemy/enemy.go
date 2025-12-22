package enemy

import (
	"coldkiller2/animation"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Enemy struct {
	Model         rl.Model
	ModelAngleDeg float32

	Animation             []rl.ModelAnimation
	AnimationState        animation.ActionState
	AnimationIdx          int
	AnimationCurrentFrame int32
	AnimationFrameCounter float32
	AnimationFrameSpeed   float32
	AnimationReplay       bool

	MoveDirection   rl.Vector3
	TargetDirection rl.Vector3
	Position        rl.Vector3
	Size            float32
	MoveSpeed       float32
	AttackSound     rl.Sound
	ActionTimeLeft  float32
	Health          int32
	IsDead          bool
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

func (e *Enemy) Mutate(dt float32) []BulletCmd {
	var bulletCmds []BulletCmd
	if e.ActionTimeLeft > 0 {
		e.ActionTimeLeft -= dt
		return bulletCmds
	}

	if e.ActionTimeLeft <= 0 {
		e.AnimationState = animation.StateIdle
	}

	if e.ActionTimeLeft <= 0 && e.Health <= 0 {
		e.IsDead = true
	}

	return []BulletCmd{}
}

func (e *Enemy) Damage(d int32) {
	e.Health -= d
	e.AnimationState = animation.StateAttacking
	e.ActionTimeLeft = 0.2
	if e.Health <= 0 {
		e.AnimationState = animation.StateDying
		e.ActionTimeLeft = 10
	}
}

func (e *Enemy) Unload() {
	rl.UnloadModel(e.Model)
	rl.UnloadModelAnimations(e.Animation)
}

func (e *Enemy) ResolveAnimation() {
	// 0 dance
	// 1 death
	// 2 idle
	// 3 jump
	// 4 no
	// 5 punch
	// 6 running
	// 7 sitting
	// 8 standing
	// 9 thumbsup
	switch e.AnimationState {
	case animation.StateIdle:
		e.setAnim(2, 24, true)
	case animation.StateRunning:
		e.setAnim(6, 100, true)
	case animation.StateAttacking:
		e.setAnim(3, 150, false)
	case animation.StateDying:
		e.setAnim(1, 56, false)
	}
}

func (e *Enemy) setAnim(idx int, speed float32, loop bool) {
	if e.AnimationIdx != idx {
		e.AnimationIdx = idx
		e.AnimationCurrentFrame = 0
		e.AnimationFrameCounter = 0
	}
	e.AnimationFrameSpeed = speed
	e.AnimationReplay = loop
}

func (e *Enemy) PlanAnimate(dt float32) {
	e.AnimationFrameCounter += e.AnimationFrameSpeed * dt
	anim := e.Animation[e.AnimationIdx]
	for e.AnimationFrameCounter >= 1.0 {
		e.AnimationCurrentFrame++
		e.AnimationFrameCounter -= 1.0
		if e.AnimationReplay == false && e.AnimationCurrentFrame >= anim.FrameCount-5 {
			e.AnimationCurrentFrame = anim.FrameCount - 5
			return
		}
	}
}

func (e *Enemy) Animate() {
	anim := e.Animation[e.AnimationIdx]
	rl.UpdateModelAnimation(e.Model, anim, e.AnimationCurrentFrame)
}
