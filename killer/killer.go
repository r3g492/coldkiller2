package killer

import (
	"coldkiller2/animation"
	"coldkiller2/input"
	"coldkiller2/model"
	"coldkiller2/sound"
	"coldkiller2/structure"
	"coldkiller2/util"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Killer struct {
	Model         rl.Model
	ModelAngleDeg float32

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
	Camera                rl.Camera3D
	ActionTimeLeft        float32
	MaxActionTime         float32
	Health                int32
	AmmoCapacity          int32
	Ammo                  int32
	FootstepSoundTimeLeft float32
	FootstepSoundTimeUnit float32
	float32
	FootstepSound    rl.Sound
	HitFlashTimer    float32
	DashTimeLeft     float32
	DashCooldown     float32
	DashPushTimeLeft float32
	DashDirection    rl.Vector3
	CameraOffset     rl.Vector3
}

const ModelRatio = 0.2
const CharSize = 0.72

func Create() *Killer {
	playerModel := model.KillerModel
	playerAnimation := model.KillerAnimation
	playerPosition := rl.Vector3{X: 0, Y: 0, Z: 0}
	return &Killer{
		Model:           playerModel,
		ModelAngleDeg:   0,
		Animation:       playerAnimation,
		MoveDirection:   rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection: rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:        playerPosition,
		Size:            CharSize,
		MoveSpeed:       5,
		Camera: rl.Camera3D{
			Position:   rl.Vector3Add(playerPosition, rl.NewVector3(0.0, 10.0, 0.0)),
			Target:     playerPosition,
			Up:         rl.NewVector3(0.0, 0.0, -1),
			Fovy:       30.0,
			Projection: rl.CameraOrthographic,
		},
		ActionTimeLeft:        0,
		AmmoCapacity:          6,
		Ammo:                  6,
		FootstepSoundTimeUnit: 0.4,
		FootstepSoundTimeLeft: 0.4,
		FootstepSound:         rl.LoadSoundAlias(sound.FootStep),
	}
}

func (k *Killer) Init() {
	k.Position = rl.Vector3{X: 0, Y: 0, Z: 0}
	k.Health = 100
	k.Ammo = k.AmmoCapacity
	k.MoveDirection = rl.Vector3{X: 0, Y: 0, Z: 0}
	k.TargetDirection = rl.Vector3{X: 0, Y: 0, Z: 0}
	k.ModelAngleDeg = 0
	k.ActionTimeLeft = 0
	k.DashTimeLeft = 0
	k.DashCooldown = 0
}

func (k *Killer) Unload() {
	rl.UnloadModel(k.Model)
	if len(k.Animation) > 0 {
		rl.UnloadModelAnimations(k.Animation)
	}
}

func (k *Killer) Draw3D() {
	rl.DrawCylinder(
		rl.Vector3{X: k.Position.X, Y: -1, Z: k.Position.Z + CharSize*0.4},
		CharSize*0.4, CharSize*0.4, 0.01, 16,
		rl.NewColor(0, 0, 0, 40),
	)

	rl.PushMatrix()
	rl.Translatef(
		k.Position.X,
		k.Position.Y,
		k.Position.Z,
	)
	rl.Rotatef(-30, 1, 0, 0)
	rl.Rotatef(k.ModelAngleDeg, 0, 1, 0)
	rl.DrawModel(k.Model, rl.NewVector3(0, -k.Size, 0), ModelRatio, rl.DarkGreen)
	if k.IsAlive() {
		// rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, k.Size*2, k.Size*2, k.Size*2, rl.Green)
	}
	rl.PopMatrix()

	rl.PushMatrix()
	rl.Translatef(
		k.Position.X,
		k.Position.Y,
		k.Position.Z,
	)
	rl.PopMatrix()
}

func (k *Killer) DrawUi() {
	uiWorldPos := rl.Vector3{X: k.Position.X, Y: k.Position.Y + 3.0, Z: k.Position.Z}
	screenPos := rl.GetWorldToScreenEx(uiWorldPos, k.Camera, util.VirtualWidth, util.VirtualHeight)

	var _ = int32(20)
	barWidth := float32(60)
	barHeight := float32(8)
	barX := screenPos.X - barWidth/2

	if k.AmmoCapacity > 0 {
		ammoPct := float32(k.Ammo) / float32(k.AmmoCapacity)
		ammoFillWidth := ammoPct * barWidth
		ammoBarY := screenPos.Y + 25

		rl.DrawRectangleRec(rl.NewRectangle(barX, ammoBarY, barWidth, barHeight), rl.Fade(rl.DarkGray, 0.6))

		ammoColor := rl.SkyBlue
		if ammoPct < 0.2 {
			ammoColor = rl.Red
		}

		rl.DrawRectangleRec(rl.NewRectangle(barX, ammoBarY, ammoFillWidth, barHeight), ammoColor)
		rl.DrawRectangleLinesEx(rl.NewRectangle(barX, ammoBarY, barWidth, barHeight), 1, rl.Black)
	}

	if k.ActionTimeLeft > 0 && k.MaxActionTime > 0 {
		pct := k.ActionTimeLeft / k.MaxActionTime
		fillWidth := pct * barWidth
		barY := screenPos.Y + 37

		rl.DrawRectangleRec(rl.NewRectangle(barX, barY, barWidth, barHeight), rl.Fade(rl.DarkGray, 0.6))
		rl.DrawRectangleRec(rl.NewRectangle(barX, barY, fillWidth, barHeight), rl.Yellow)
		rl.DrawRectangleLinesEx(rl.NewRectangle(barX, barY, barWidth, barHeight), 1, rl.Black)
	}

	/*fpsText := fmt.Sprintf("%d", rl.GetFPS())
	rl.DrawText(fpsText, int32(screenPos.X)-20, int32(screenPos.Y)+60, 1, rl.Yellow)*/
}
func (k *Killer) Mutate(
	input input.Input,
	dt float32,
	obstacles []rl.BoundingBox,
	structureManager *structure.Manager,
) []BulletCmd {
	var bulletCmds []BulletCmd

	if k.IsAlive() {
		mouseMovement(input, k)
	}
	if rl.Vector3LengthSqr(k.TargetDirection) > 0 {
		aimDir := rl.Vector3Normalize(k.TargetDirection)
		targetOffset := rl.Vector3Scale(aimDir, 2.5)
		k.CameraOffset.X += (targetOffset.X - k.CameraOffset.X) * 4.0 * dt
		k.CameraOffset.Z += (targetOffset.Z - k.CameraOffset.Z) * 4.0 * dt
	}
	attack := false
	if k.ActionTimeLeft <= 0 {
		bulletCmds, attack = k.attack(input)
		if attack {
			var attackTime float32 = 0.1
			k.ActionTimeLeft = attackTime
			k.MaxActionTime = attackTime
			k.AnimationState = animation.StateAttacking
			k.AnimationCurrentFrame = 0
		}
	}

	if input.ReloadPressed && k.ActionTimeLeft <= 0 {
		rl.PlaySound(sound.ReloadingSound)
		k.Ammo = k.AmmoCapacity
		var reloadTime float32 = 0.4
		k.ActionTimeLeft = reloadTime
		k.MaxActionTime = reloadTime
		k.AnimationState = animation.StateReloading
		k.AnimationCurrentFrame = 0
	}

	if !attack && k.ActionTimeLeft <= 0 {
		moving := k.movement(input, dt, obstacles, structureManager)

		if k.FootstepSoundTimeLeft > 0 {
			k.FootstepSoundTimeLeft -= dt
		}

		if moving {
			k.AnimationState = animation.StateRunning

			if k.FootstepSoundTimeLeft <= 0 {
				sound.PlaySound3D(k.FootstepSound, k.Position, k.Position, 0.3)
				k.FootstepSoundTimeLeft = k.FootstepSoundTimeUnit
			}
		} else {
			k.AnimationState = animation.StateIdle
			k.FootstepSoundTimeLeft = 0
		}
	}

	camTarget := rl.Vector3Add(k.Position, k.CameraOffset)
	k.Camera = rl.Camera3D{
		Position:   rl.Vector3Add(camTarget, rl.NewVector3(0.0, 10.0, 0.0)),
		Target:     camTarget,
		Up:         rl.NewVector3(0.0, 0.0, -1),
		Fovy:       30.0,
		Projection: rl.CameraOrthographic,
	}

	k.ActionTimeLeft -= dt
	if k.HitFlashTimer > 0 {
		k.HitFlashTimer -= dt
	}
	if k.DashTimeLeft > 0 {
		k.DashTimeLeft -= dt
	}
	if k.DashPushTimeLeft > 0 {
		k.DashPushTimeLeft -= dt
	}
	if k.DashCooldown > 0 {
		k.DashCooldown -= dt
	}
	return bulletCmds
}

func mouseMovement(input input.Input, k *Killer) {
	mouseLocation := input.MouseLocation
	ray := rl.GetScreenToWorldRayEx(mouseLocation, k.Camera, util.VirtualWidth, util.VirtualHeight)
	targetOnXzPlane := rl.Vector3{
		X: ray.Position.X,
		Y: 0,
		Z: ray.Position.Z,
	}
	k.TargetDirection = rl.Vector3Subtract(targetOnXzPlane, k.Position)
	angleRad := math.Atan2(float64(k.TargetDirection.X), float64(k.TargetDirection.Z))
	k.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))
}

func (k *Killer) movement(
	input input.Input,
	dt float32,
	obstacles []rl.BoundingBox,
	structureManager *structure.Manager,
) bool {
	k.PrevPosition = k.Position
	k.MoveDirection = rl.Vector3{}
	if input.MoveUp {
		k.MoveDirection.Z -= 1
	}
	if input.MoveDown {
		k.MoveDirection.Z += 1
	}
	if input.MoveLeft {
		k.MoveDirection.X -= 1
	}
	if input.MoveRight {
		k.MoveDirection.X += 1
	}
	isMoving := rl.Vector3LengthSqr(k.MoveDirection) > 0
	if isMoving {
		k.MoveDirection = rl.Vector3Normalize(k.MoveDirection)
	}
	if input.DashPressed && isMoving && k.DashCooldown <= 0 {
		rl.PlaySound(sound.Dash)
		k.DashTimeLeft = 0.3
		k.DashPushTimeLeft = 0.4
		k.DashCooldown = 1.0
		k.DashDirection = k.MoveDirection
	}
	speed := k.MoveSpeed
	if k.DashTimeLeft > 0 {
		speed = k.MoveSpeed * 3.5
		k.MoveDirection = k.DashDirection
	}
	moveAmount := rl.Vector3Scale(k.MoveDirection, speed*dt)
	if rl.Vector3Length(moveAmount) > 0 {
		oldPos := k.Position
		k.Position.X += moveAmount.X
		if k.isColliding(obstacles) || structureManager.CheckCollision(k.Position, k.PrevPosition, rl.Vector3{X: k.Size, Y: k.Size, Z: k.Size}) {
			k.Position.X = oldPos.X
		}
		k.Position.Z += moveAmount.Z
		if k.isColliding(obstacles) || structureManager.CheckCollision(k.Position, k.PrevPosition, rl.Vector3{X: k.Size, Y: k.Size, Z: k.Size}) {
			k.Position.Z = oldPos.Z
		}
		return k.Position != oldPos
	}
	return false
}

func (k *Killer) attack(input input.Input) ([]BulletCmd, bool) {
	var bulletCmds []BulletCmd
	if input.FireHold && k.Ammo > 0 && k.DashTimeLeft <= 0 {
		rl.PlaySound(sound.ShotgunSound)
		angleRad := math.Atan2(float64(k.TargetDirection.X), float64(k.TargetDirection.Z))
		k.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))
		dir := rl.Vector3Normalize(k.TargetDirection)
		spawnPos := rl.Vector3Add(k.Position, rl.Vector3{X: 0, Y: 0, Z: 0})
		bulletCmds = append(bulletCmds, BulletCmd{spawnPos, dir, 200})
		k.Ammo--
		return bulletCmds, true
	}
	return bulletCmds, false
}

func (k *Killer) ResolveAnimation() {
	if len(k.Animation) == 0 {
		return
	}
	switch k.AnimationState {
	case animation.StateIdle:
		k.setAnim(0, 24, true)
	case animation.StateRunning:
		k.setAnim(1, 180, true)
	case animation.StateAttacking:
		k.setAnim(2, 150, false)
	case animation.StateDying:
		k.setAnim(3, 300, false)
	case animation.StateReloading:
		k.setAnim(2, 150, false)
	default:
		panic("unhandled default case")
	}
}

func (k *Killer) setAnim(idx int, speed float32, loop bool) {
	if k.AnimationIdx != idx {
		k.AnimationIdx = idx
		k.AnimationCurrentFrame = 0
		k.AnimationFrameCounter = 0
	}
	k.AnimationFrameSpeed = speed
	k.AnimationReplay = loop
}

func (k *Killer) PlanAnimate(dt float32) {
	if len(k.Animation) == 0 {
		return
	}
	k.AnimationFrameCounter += k.AnimationFrameSpeed * dt
	anim := k.Animation[k.AnimationIdx]
	for k.AnimationFrameCounter >= 1.0 {
		k.AnimationCurrentFrame++
		k.AnimationFrameCounter -= 1.0
		if k.AnimationReplay == false && k.AnimationCurrentFrame >= anim.FrameCount-5 {
			k.AnimationCurrentFrame = anim.FrameCount - 5
			return
		}
	}
}

func (k *Killer) Animate() {
	if len(k.Animation) == 0 {
		return
	}
	anim := k.Animation[k.AnimationIdx]
	rl.UpdateModelAnimation(k.Model, anim, k.AnimationCurrentFrame)
}

func (k *Killer) isColliding(obstacles []rl.BoundingBox) bool {
	myBox := k.GetBoundingBox()
	for _, box := range obstacles {
		if rl.CheckCollisionBoxes(myBox, box) {
			return true
		}
	}
	return false
}

func (k *Killer) GetBoundingBox() rl.BoundingBox {
	return rl.BoundingBox{
		Min: rl.Vector3{X: k.Position.X - k.Size, Y: k.Position.Y - k.Size, Z: k.Position.Z - k.Size},
		Max: rl.Vector3{X: k.Position.X + k.Size, Y: k.Position.Y + k.Size, Z: k.Position.Z + k.Size},
	}
}

func (k *Killer) Damage(d int32) {
	k.Health -= d
	k.HitFlashTimer = 0.35
	rl.SetSoundVolume(sound.ShotNew, 0.8)
	rl.PlaySound(sound.ShotNew)
	k.AnimationState = animation.StateDying
	var shotTime float32 = 0.1
	k.ActionTimeLeft = shotTime
	k.MaxActionTime = shotTime
	if !k.IsAlive() {
		k.AnimationState = animation.StateDying
		var dyingTime float32 = 1
		k.ActionTimeLeft = dyingTime
		k.MaxActionTime = dyingTime
	}
}

func (k *Killer) DrawHitFlash() {
	if k.HitFlashTimer <= 0 {
		return
	}
	alpha := uint8(k.HitFlashTimer / 0.35 * 120)
	w := rl.GetScreenWidth()
	h := rl.GetScreenHeight()
	rl.DrawRectangle(0, 0, int32(w), int32(h), rl.NewColor(220, 30, 30, alpha))
}

func (k *Killer) IsAlive() bool {
	return k.Health > 0
}
