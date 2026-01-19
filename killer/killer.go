package killer

import (
	"coldkiller2/animation"
	"coldkiller2/input"
	"coldkiller2/sound"
	"fmt"
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

	MoveDirection   rl.Vector3
	TargetDirection rl.Vector3
	Position        rl.Vector3
	Size            float32
	MoveSpeed       float32
	Camera          rl.Camera3D
	ActionTimeLeft  float32
	MaxActionTime   float32
	Health          int32
	AmmoCapacity    int32
	Ammo            int32
}

func Init() *Killer {
	playerModel := rl.LoadModel("resources/unit_v3.glb")
	playerAnimation := rl.LoadModelAnimations("resources/unit_v3.glb")
	playerPosition := rl.Vector3{X: 0, Y: 0, Z: 0}
	return &Killer{
		Model:           playerModel,
		ModelAngleDeg:   0,
		Animation:       playerAnimation,
		MoveDirection:   rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection: rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:        playerPosition,
		Size:            1.0,
		MoveSpeed:       20.0,
		Camera: rl.Camera3D{
			Position:   rl.Vector3Add(playerPosition, rl.NewVector3(0.0, 10.0, 0.0)),
			Target:     playerPosition,
			Up:         rl.NewVector3(0.0, 0.0, -1),
			Fovy:       30.0,
			Projection: rl.CameraOrthographic,
		},
		ActionTimeLeft: 0,
		Health:         100,
		AmmoCapacity:   6,
		Ammo:           6,
	}
}

func (k *Killer) Unload() {
	rl.UnloadModel(k.Model)
	rl.UnloadModelAnimations(k.Animation)
}

func (k *Killer) Draw3D() {
	rl.PushMatrix()
	rl.Translatef(
		k.Position.X,
		k.Position.Y,
		k.Position.Z,
	)
	rl.Rotatef(-60, 1, 0, 0)
	rl.Rotatef(k.ModelAngleDeg, 0, 1, 0)
	rl.DrawModel(k.Model, rl.NewVector3(0, -k.Size, 0), 0.45, rl.Green)
	rl.PopMatrix()

	rl.PushMatrix()
	rl.Translatef(
		k.Position.X,
		k.Position.Y,
		k.Position.Z,
	)
	// rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, k.Size*2, k.Size*2, k.Size*2, rl.Green)
	rl.PopMatrix()
}

func (k *Killer) DrawUI() {
	uiWorldPos := rl.Vector3{X: k.Position.X, Y: k.Position.Y + 3.0, Z: k.Position.Z}
	screenPos := rl.GetWorldToScreen(uiWorldPos, k.Camera)

	ammoText := fmt.Sprintf("%d / %d", k.Ammo, k.AmmoCapacity)
	fontSize := int32(20)
	textWidth := rl.MeasureText(ammoText, fontSize)
	rl.DrawText(ammoText, int32(screenPos.X)-textWidth/2, int32(screenPos.Y), fontSize, rl.White)

	if k.ActionTimeLeft > 0 && k.MaxActionTime > 0 {
		barWidth := float32(60)
		barHeight := float32(8)
		pct := k.ActionTimeLeft / k.MaxActionTime
		fillWidth := pct * barWidth

		barX := screenPos.X - barWidth/2
		barY := screenPos.Y + 25

		rl.DrawRectangleRec(rl.NewRectangle(barX, barY, barWidth, barHeight), rl.DarkGray)
		rl.DrawRectangleRec(rl.NewRectangle(barX, barY, fillWidth, barHeight), rl.Yellow)
		rl.DrawRectangleLinesEx(rl.NewRectangle(barX, barY, barWidth, barHeight), 1, rl.Black)
	}
}

func (k *Killer) Mutate(input input.Input, dt float32, obstacles []rl.BoundingBox) []BulletCmd {
	var bulletCmds []BulletCmd

	if k.IsAlive() {
		mouseMovement(input, k)
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
		var reloadTime float32 = 0.7
		k.ActionTimeLeft = reloadTime
		k.MaxActionTime = reloadTime
		k.AnimationState = animation.StateReloading
		k.AnimationCurrentFrame = 0
	}

	if !attack && k.ActionTimeLeft <= 0 {
		moving := k.movement(input, dt, obstacles)
		k.Camera = rl.Camera3D{
			Position:   rl.Vector3Add(k.Position, rl.NewVector3(0.0, 10.0, 0.0)),
			Target:     k.Position,
			Up:         rl.NewVector3(0.0, 0.0, -1),
			Fovy:       30.0,
			Projection: rl.CameraOrthographic,
		}
		if moving {
			k.AnimationState = animation.StateRunning
		} else {
			k.AnimationState = animation.StateIdle
		}
	}

	k.ActionTimeLeft -= dt
	return bulletCmds
}

func mouseMovement(input input.Input, k *Killer) {
	mouseLocation := input.MouseLocation
	ray := rl.GetScreenToWorldRay(mouseLocation, k.Camera)
	targetOnXzPlane := rl.Vector3{
		X: ray.Position.X,
		Y: 0,
		Z: ray.Position.Z,
	}
	k.TargetDirection = rl.Vector3Subtract(targetOnXzPlane, k.Position)
	angleRad := math.Atan2(float64(k.TargetDirection.X), float64(k.TargetDirection.Z))
	k.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))
}

func (k *Killer) movement(input input.Input, dt float32, obstacles []rl.BoundingBox) bool {
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
	if rl.Vector3LengthSqr(k.MoveDirection) > 0 {
		k.MoveDirection = rl.Vector3Normalize(k.MoveDirection)
	}
	moveAmount := rl.Vector3Scale(k.MoveDirection, k.MoveSpeed*dt)
	if rl.Vector3Length(moveAmount) > 0 {
		oldPos := k.Position
		k.Position.X += moveAmount.X
		if k.isColliding(obstacles) {
			k.Position.X = oldPos.X
		}
		k.Position.Z += moveAmount.Z
		if k.isColliding(obstacles) {
			k.Position.Z = oldPos.Z
		}
		return k.Position != oldPos
	}
	return false
}

func (k *Killer) attack(input input.Input) ([]BulletCmd, bool) {
	var bulletCmds []BulletCmd
	if input.PunchPressed && k.Ammo > 0 {
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
	k.AnimationState = animation.StateDying
	var shotTime float32 = 0.1
	k.ActionTimeLeft = shotTime
	k.MaxActionTime = shotTime
	if k.Health <= 0 {
		k.AnimationState = animation.StateDying
		var dyingTime float32 = 1
		k.ActionTimeLeft = dyingTime
		k.MaxActionTime = dyingTime
	}
}

func (k *Killer) IsAlive() bool {
	return k.Health > 0
}
