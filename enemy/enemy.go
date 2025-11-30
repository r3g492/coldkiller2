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
	Health                int32
	IsDead                bool
}

func (e *Enemy) Draw3D() {
	if e.IsDead {
		return
	}
	anim := e.Animation[e.AnimationIdx]
	rl.UpdateModelAnimation(e.Model, anim, e.AnimationCurrentFrame)
	rl.PushMatrix()
	rl.Translatef(e.Position.X, e.Position.Y, e.Position.Z)
	rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, e.Size*2, e.Size*2, e.Size*2, rl.Red)
	rl.Rotatef(270, 1, 0, 0)
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
	// TODO: implment movement and bullet cmds
	return bulletCmds
}

func (e *Enemy) Damage(d int32) {
	e.Health -= d
	if e.Health <= 0 {
		e.IsDead = true
	}
}
