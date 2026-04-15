package enemy

import (
	"coldkiller2/animation"
	"coldkiller2/killer"
	"coldkiller2/model"
	"coldkiller2/sound"
	"coldkiller2/structure"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var enemyRed = rl.NewColor(150, 20, 25, 255)

type Enemy struct {
	Model         rl.Model
	ModelAngleDeg float32
	ModelRatio    float32

	Animation             []rl.ModelAnimation
	AnimationState        animation.ActionState
	AnimationIdx          int
	AnimationCurrentFrame int32
	AnimationFrameCounter float32
	AnimationFrameSpeed   float32
	AnimationReplay       bool

	MoveDirection         rl.Vector3
	TargetDirection       rl.Vector3
	Position              rl.Vector3
	PrevPosition          rl.Vector3
	Size                  float32
	MoveSpeed             float32
	ActionTimeLeft        float32
	Health                int32
	ShouldBeDeleted       bool
	AttackCooldown        time.Duration
	LastAttack            time.Time
	AttackRange           float32
	AimTimeLeft           float32
	AimTimeUnit           float32
	AimDirection          rl.Vector3
	FootstepSoundTimeLeft float32
	FootstepSoundTimeUnit float32
	FootstepSound         rl.Sound
	IsHiddenFromKiller    bool
	AiType                AiType
	KnockbackVelocity     rl.Vector3
	KnockbackTimeLeft     float32
}

func (e *Enemy) IsAlive() bool {
	return e.Health > 0
}

func (e *Enemy) ApplyKnockback(velocity rl.Vector3, duration float32) {
	e.KnockbackVelocity = velocity
	e.KnockbackTimeLeft = duration
}

func (e *Enemy) Draw3D(p *killer.Killer) {
	if e.IsHiddenFromKiller {
		return
	}

	if e.IsAlive() {
		rl.DrawCylinder(
			rl.Vector3{X: e.Position.X, Y: -1, Z: e.Position.Z + 0.3},
			e.Size*0.4, e.Size*0.4, 0.01, 16,
			rl.NewColor(0, 0, 0, 40),
		)
	}

	anim := e.Animation[e.AnimationIdx]
	rl.UpdateModelAnimation(e.Model, anim, e.AnimationCurrentFrame)
	rl.PushMatrix()
	rl.Translatef(e.Position.X, e.Position.Y, e.Position.Z)
	rl.Rotatef(-30, 1, 0, 0)
	rl.Rotatef(e.ModelAngleDeg, 0, 1, 0)
	if e.IsAlive() {
		rl.DrawModel(e.Model, rl.NewVector3(0, -e.Size, 0), e.ModelRatio, enemyRed)
		// rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, e.Size*2, e.Size*2, e.Size*2, enemyRed)
	} else {
		rl.DrawModel(e.Model, rl.NewVector3(0, -e.Size, 0), e.ModelRatio, rl.DarkGray)
	}
	rl.PopMatrix()
	if e.AnimationState == animation.StateAiming && e.AimDirection != (rl.Vector3{}) {
		aimEnd := rl.Vector3Add(e.Position, rl.Vector3Scale(rl.Vector3Normalize(e.AimDirection), e.AttackRange))
		rl.DrawLine3D(e.Position, aimEnd, enemyRed)
	}
}

func (e *Enemy) DrawUI(p *killer.Killer) {
	if e.IsHiddenFromKiller {
		return
	}

	uiWorldPos := rl.Vector3{X: e.Position.X, Y: e.Position.Y + 3.0, Z: e.Position.Z}
	screenPos := rl.GetWorldToScreen(uiWorldPos, p.Camera)

	if e.AimTimeLeft > 0 && e.AimTimeLeft != e.AimTimeUnit && e.IsAlive() {
		barWidth := float32(40)
		barHeight := float32(8)
		pct := e.AimTimeLeft / e.AimTimeUnit
		fillWidth := pct * barWidth

		barX := screenPos.X - barWidth/2
		barY := screenPos.Y + 25

		rl.DrawRectangleRec(rl.NewRectangle(barX, barY, barWidth, barHeight), rl.DarkGray)
		rl.DrawRectangleRec(rl.NewRectangle(barX, barY, fillWidth, barHeight), rl.Yellow)
		rl.DrawRectangleLinesEx(rl.NewRectangle(barX, barY, barWidth, barHeight), 1, rl.Black)
	}
}

func (e *Enemy) Mutate(
	dt float32,
	p killer.Killer,
	em *Manager,
	myIdx int,
	structureManager *structure.Manager,
) []BulletCmd {
	if e.KnockbackTimeLeft > 0 {
		kbMove := rl.Vector3Scale(e.KnockbackVelocity, dt)
		e.PrevPosition = e.Position
		oldPos := e.Position
		e.Position.X += kbMove.X
		if structureManager.CheckCollision(e.Position, e.PrevPosition, rl.Vector3{X: e.Size, Y: e.Size, Z: e.Size}) {
			e.Position.X = oldPos.X
		}
		e.Position.Z += kbMove.Z
		if structureManager.CheckCollision(e.Position, e.PrevPosition, rl.Vector3{X: e.Size, Y: e.Size, Z: e.Size}) {
			e.Position.Z = oldPos.Z
		}
		// cascade: transfer knockback to any enemy we moved into
		if hit := em.findCollidingEnemy(myIdx, e); hit != nil {
			speed := rl.Vector3Length(e.KnockbackVelocity)
			if speed > 2 {
				dir := rl.Vector3Subtract(hit.Position, e.Position)
				if rl.Vector3LengthSqr(dir) < 0.0001 {
					dir = e.KnockbackVelocity
				}
				hit.ApplyKnockback(rl.Vector3Scale(rl.Vector3Normalize(dir), speed*0.8), e.KnockbackTimeLeft*0.8)
			}
		}
		e.KnockbackVelocity = rl.Vector3Scale(e.KnockbackVelocity, 1-8*dt)
		e.KnockbackTimeLeft -= dt
	}

	vecToPlayer := rl.Vector3Subtract(p.Position, e.Position)
	var bulletCmds []BulletCmd
	if e.ActionTimeLeft > 0 {
		e.ActionTimeLeft -= dt
		return []BulletCmd{}
	}
	if e.ActionTimeLeft <= 0 {
		e.AnimationState = animation.StateIdle
	}
	if e.ActionTimeLeft <= 0 && !e.IsAlive() && e.KnockbackTimeLeft <= 0 {
		e.ShouldBeDeleted = true
	}

	var derivedAimStart, derivedMovement = deriveAi(e, em, myIdx, &p, structureManager)
	if e.AimDirection != (rl.Vector3{}) {
		derivedAimStart = true
	}

	if e.AimTimeLeft <= 0 {
		dir := rl.Vector3Normalize(e.AimDirection)
		e.ActionTimeLeft = 1
		e.AnimationState = animation.StateAttacking
		e.AnimationCurrentFrame = 0
		e.AimTimeLeft = e.AimTimeUnit
		e.AimDirection = rl.Vector3{}
		rl.PlaySound(sound.ShotgunSound)
		spawnPos := e.Position
		bulletCmds = append(bulletCmds, BulletCmd{Pos: spawnPos, Dir: dir, Damage: 200, Range: e.AttackRange, Shooter: e})
		return bulletCmds
	}

	if derivedAimStart {
		if e.AimDirection == (rl.Vector3{}) {
			e.AimDirection = vecToPlayer
		}
		e.TargetDirection = e.AimDirection
		angleRad := math.Atan2(float64(e.TargetDirection.X), float64(e.TargetDirection.Z))
		e.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))
		e.AimTimeLeft -= dt
		e.AnimationState = animation.StateAiming
		e.AnimationCurrentFrame = 0
		return []BulletCmd{}
	}

	e.AimTimeLeft = e.AimTimeUnit
	e.AimDirection = rl.Vector3{}
	e.MoveDirection = derivedMovement

	moveAmount := rl.Vector3Scale(e.MoveDirection, e.MoveSpeed*dt)

	e.PrevPosition = e.Position
	oldPos := e.Position
	e.Position.X += moveAmount.X
	if e.isCollidingWithGrid(myIdx, em, p.GetBoundingBox()) || structureManager.CheckCollision(e.Position, e.PrevPosition, rl.Vector3{X: e.Size, Y: e.Size, Z: e.Size}) {
		e.Position.X = oldPos.X
	}
	e.Position.Z += moveAmount.Z
	if e.isCollidingWithGrid(myIdx, em, p.GetBoundingBox()) || structureManager.CheckCollision(e.Position, e.PrevPosition, rl.Vector3{X: e.Size, Y: e.Size, Z: e.Size}) {
		e.Position.Z = oldPos.Z
	}

	moving := rl.Vector3Distance(oldPos, e.Position) > 0.01
	if e.FootstepSoundTimeLeft > 0 {
		e.FootstepSoundTimeLeft -= dt
	}

	if moving {
		e.AnimationState = animation.StateRunning

		if e.FootstepSoundTimeLeft <= 0 {
			sound.PlaySound3D(e.FootstepSound, e.Position, p.Position, 0.5)
			e.FootstepSoundTimeLeft = e.FootstepSoundTimeUnit
		}
	} else {
		e.FootstepSoundTimeLeft = 0
	}

	e.TargetDirection = e.MoveDirection
	angleRad := math.Atan2(float64(e.TargetDirection.X), float64(e.TargetDirection.Z))
	e.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))
	return bulletCmds
}

func (e *Enemy) Damage(d int32, bulletDir rl.Vector3) {
	e.Health -= d
	e.AnimationState = animation.StateDying
	e.ActionTimeLeft = 0.1
	if rl.Vector3LengthSqr(bulletDir) > 0 {
		e.KnockbackVelocity = rl.Vector3Scale(rl.Vector3Normalize(bulletDir), 60.0)
		e.KnockbackTimeLeft = 0.2
	}
	if !e.IsAlive() {
		e.AnimationState = animation.StateDying
		e.ActionTimeLeft = 10
	}
}

func (e *Enemy) Unload() {
	rl.UnloadModel(e.Model)
	rl.UnloadModelAnimations(e.Animation)
	rl.UnloadSoundAlias(e.FootstepSound)
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
	default:
		panic("unhandled default case")
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

func (e *Enemy) GetBoundingBox() rl.BoundingBox {
	return rl.BoundingBox{
		Min: rl.Vector3{X: e.Position.X - e.Size, Y: e.Position.Y - e.Size, Z: e.Position.Z - e.Size},
		Max: rl.Vector3{X: e.Position.X + e.Size, Y: e.Position.Y + e.Size, Z: e.Position.Z + e.Size},
	}
}

func Soldier(x, z float32) *Enemy {
	return &Enemy{
		Model:                 model.SoldierModel,
		ModelRatio:            0.2,
		Animation:             model.SoldierAnimation,
		Position:              rl.Vector3{X: x, Y: 0, Z: z},
		Size:                  killer.CharSize,
		MoveSpeed:             4,
		Health:                100,
		AttackRange:           10,
		AimTimeLeft:           0.5,
		AimTimeUnit:           0.5,
		FootstepSoundTimeLeft: 0,
		FootstepSoundTimeUnit: 0.4,
		FootstepSound:         sound.FootStep,
		AiType:                Elite,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
	}
}

func Sniper(x, z float32) *Enemy {
	return &Enemy{
		Model:                 model.SoldierModel,
		ModelRatio:            0.2,
		Animation:             model.SoldierAnimation,
		Position:              rl.Vector3{X: x, Y: 0, Z: z},
		Size:                  killer.CharSize,
		MoveSpeed:             1,
		Health:                100,
		AttackRange:           30,
		AimTimeLeft:           5,
		AimTimeUnit:           5,
		FootstepSoundTimeLeft: 0,
		FootstepSoundTimeUnit: 0.4,
		FootstepSound:         sound.FootStep,
		AiType:                Elite,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
	}
}

func Robot(x, z float32) *Enemy {
	return &Enemy{
		Model:                 model.RobotModel,
		ModelRatio:            0.4,
		Animation:             model.RobotAnimation,
		Position:              rl.Vector3{X: x, Y: 0, Z: z},
		Size:                  killer.CharSize,
		MoveSpeed:             8,
		Health:                100,
		AttackRange:           8,
		AimTimeLeft:           0.5,
		AimTimeUnit:           0.5,
		FootstepSoundTimeLeft: 0,
		FootstepSoundTimeUnit: 0.4,
		FootstepSound:         sound.FootStep,
		AiType:                SimpleZombie,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
	}
}
