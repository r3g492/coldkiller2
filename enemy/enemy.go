package enemy

import (
	"coldkiller2/animation"
	"coldkiller2/killer"
	"coldkiller2/sound"
	"math"
	"time"

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
	ActionTimeLeft  float32
	Health          int32
	IsDead          bool
	AttackCooldown  time.Duration
	LastAttack      time.Time
	AttackRange     float32
	AimTimeLeft     float32
	AimTimeUnit     float32
}

func (e *Enemy) Draw3D(p *killer.Killer) {
	anim := e.Animation[e.AnimationIdx]
	rl.UpdateModelAnimation(e.Model, anim, e.AnimationCurrentFrame)
	rl.PushMatrix()
	rl.Translatef(e.Position.X, e.Position.Y, e.Position.Z)
	rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, e.Size*2, e.Size*2, e.Size*2, rl.Red)
	rl.Rotatef(-60, 1, 0, 0)
	rl.Rotatef(e.ModelAngleDeg, 0, 1, 0)
	rl.DrawModel(e.Model, rl.NewVector3(0, -e.Size, 0), 0.45, rl.White)
	rl.PopMatrix()

	if e.AnimationState == animation.StateAiming {
		rl.DrawLine3D(e.Position, p.Position, rl.Red)
	}
}

func (e *Enemy) Mutate(
	dt float32,
	p killer.Killer,
	enemyObstacles []rl.BoundingBox,
	myIdx int,
) []BulletCmd {
	distToPlayer := rl.Vector3Distance(e.Position, p.Position)
	vecToPlayer := rl.Vector3Subtract(p.Position, e.Position)
	var _ = rl.Vector3Normalize(vecToPlayer)

	var bulletCmds []BulletCmd
	if e.ActionTimeLeft > 0 {
		e.ActionTimeLeft -= dt
		return []BulletCmd{}
	}
	if e.ActionTimeLeft <= 0 {
		e.AnimationState = animation.StateIdle
	}

	if e.ActionTimeLeft <= 0 && e.Health <= 0 {
		e.IsDead = true
	}

	if e.AimTimeLeft <= 0 && distToPlayer <= e.AttackRange {
		e.TargetDirection = vecToPlayer
		angleRad := math.Atan2(float64(e.TargetDirection.X), float64(e.TargetDirection.Z))
		e.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))

		e.ActionTimeLeft = 1
		e.AnimationState = animation.StateAttacking
		e.AnimationCurrentFrame = 0

		e.AimTimeLeft = e.AimTimeUnit
		rl.PlaySound(sound.ShotgunSound)
		dir := rl.Vector3Normalize(e.TargetDirection)
		spawnPos := rl.Vector3Add(e.Position, rl.Vector3{X: 0, Y: 0, Z: 0})
		bulletCmds = append(bulletCmds, BulletCmd{spawnPos, dir, 200})
		return bulletCmds
	}

	if e.AimTimeLeft > 0 && distToPlayer <= e.AttackRange && e.IsAlive() {
		e.TargetDirection = vecToPlayer
		angleRad := math.Atan2(float64(e.TargetDirection.X), float64(e.TargetDirection.Z))
		e.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))

		e.AimTimeLeft -= dt
		e.AnimationState = animation.StateAiming
		e.AnimationCurrentFrame = 0
		return []BulletCmd{}
	}
	e.AimTimeLeft = e.AimTimeUnit

	e.MoveDirection = rl.Vector3Normalize(
		rl.Vector3Subtract(
			p.Position,
			e.Position,
		),
	)

	moveAmount := rl.Vector3Scale(e.MoveDirection, e.MoveSpeed*dt)

	oldPos := e.Position
	e.Position.X += moveAmount.X
	if e.isColliding(myIdx, enemyObstacles, p.GetBoundingBox()) {
		e.Position.X = oldPos.X
	}
	e.Position.Z += moveAmount.Z
	if e.isColliding(myIdx, enemyObstacles, p.GetBoundingBox()) {
		e.Position.Z = oldPos.Z
	}

	moving := rl.Vector3Distance(oldPos, e.Position) > 0.01
	if moving {
		e.AnimationState = animation.StateRunning
	}

	e.TargetDirection = e.MoveDirection
	angleRad := math.Atan2(float64(e.TargetDirection.X), float64(e.TargetDirection.Z))
	e.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))
	return bulletCmds
}

func (e *Enemy) Damage(d int32) {
	e.Health -= d
	e.AnimationState = animation.StateDying
	e.ActionTimeLeft = 0.1
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
	switch e.AnimationState {
	case animation.StateIdle:
		e.setAnim(0, 24, true)
	case animation.StateRunning:
		e.setAnim(1, 180, true)
	case animation.StateAttacking:
		e.setAnim(2, 150, false)
	case animation.StateDying:
		e.setAnim(3, 200, false)
	case animation.StateAiming:
		e.setAnim(2, 0, false)
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

func (e *Enemy) GetBoundingBox() rl.BoundingBox {
	return rl.BoundingBox{
		Min: rl.Vector3{X: e.Position.X - e.Size, Y: e.Position.Y - e.Size, Z: e.Position.Z - e.Size},
		Max: rl.Vector3{X: e.Position.X + e.Size, Y: e.Position.Y + e.Size, Z: e.Position.Z + e.Size},
	}
}

func (e *Enemy) IsAlive() bool {
	return e.Health > 0
}
