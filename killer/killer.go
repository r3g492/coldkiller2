package killer

import (
	"coldkiller2/input"
	"coldkiller2/util"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Killer struct {
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
	AttackTimeLeft        float32
}

func Init() *Killer {
	playerModel := rl.LoadModel("resources/robot.glb")
	playerAnimation := rl.LoadModelAnimations("resources/robot.glb")
	playerPosition := rl.Vector3{X: 0, Y: 0, Z: 0}
	shotGunSound := util.LoadSoundFromEmbedded("shotgun-03-38220.mp3")
	return &Killer{
		Model:         playerModel,
		ModelAngleDeg: 0,
		Animation:     playerAnimation,
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
		AnimationIdx:          2,
		AnimationCurrentFrame: 0,
		AnimationFrameCounter: 0,
		AnimationFrameSpeed:   24,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:              playerPosition,
		Size:                  2,
		MoveSpeed:             20.0,
		Camera: rl.Camera3D{
			Position:   rl.Vector3Add(playerPosition, rl.NewVector3(0.0, 10.0, 0.0)),
			Target:     playerPosition,
			Up:         rl.NewVector3(0.0, 0.0, -1),
			Fovy:       30.0,
			Projection: rl.CameraOrthographic,
		},
		ShotGunSound:   shotGunSound,
		AttackTimeLeft: 0,
	}
}

func (k *Killer) Unload() {
	rl.UnloadModel(k.Model)
	rl.UnloadModelAnimations(k.Animation)
}

func (k *Killer) Draw3D() {
	anim := k.Animation[k.AnimationIdx]
	rl.UpdateModelAnimation(k.Model, anim, k.AnimationCurrentFrame)
	rl.PushMatrix()
	rl.Translatef(k.Position.X, k.Position.Y, k.Position.Z)
	rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, k.Size*2, k.Size*2, k.Size*2, rl.Green)
	rl.Rotatef(-60, 1, 0, 0)
	rl.Rotatef(k.ModelAngleDeg, 0, 1, 0)
	rl.DrawModel(k.Model, rl.NewVector3(0, -k.Size, 0), 0.7, rl.White)
	rl.PopMatrix()
	rl.DrawRay(rl.NewRay(k.Position, k.TargetDirection), rl.Green)
}

func (k *Killer) Mutate(input input.Input, dt float32) []BulletCmd {
	var bulletCmds []BulletCmd
	mouseMovement(input, k)
	if input.PunchPressed {
		k.AnimationCurrentFrame = 0
		k.AnimationFrameSpeed = 96
	}
	if input.PunchHold {
		k.AnimationIdx = 7
		return bulletCmds
	}

	attack := false
	if k.AttackTimeLeft <= 0 {
		bulletCmds, attack = k.attack(input)
		if attack {
			k.AttackTimeLeft = 0.2
			k.AnimationIdx = 5
			k.AnimationFrameSpeed = 150
			k.AnimationCurrentFrame = 0
		}
	}
	move := false
	if !attack && k.AttackTimeLeft <= 0 {
		move = k.movement(input, dt)
		k.Camera = rl.Camera3D{
			Position:   rl.Vector3Add(k.Position, rl.NewVector3(0.0, 10.0, 0.0)),
			Target:     k.Position,
			Up:         rl.NewVector3(0.0, 0.0, -1),
			Fovy:       30.0,
			Projection: rl.CameraOrthographic,
		}
		if move {
			k.AnimationIdx = 6
			k.AnimationFrameSpeed = 100
		}
	}

	k.AttackTimeLeft -= dt

	if k.AttackTimeLeft <= 0 && !attack && !move {
		k.AnimationIdx = 2
		k.AnimationFrameSpeed = 24
	}

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

func (k *Killer) movement(input input.Input, dt float32) bool {
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
	move := rl.Vector3Length(moveAmount) > 0
	if move {
		k.Position = rl.Vector3Add(k.Position, moveAmount)
	}
	return move
}

func (k *Killer) attack(input input.Input) ([]BulletCmd, bool) {
	var bulletCmds []BulletCmd
	if input.PunchReleased {
		rl.PlaySound(k.ShotGunSound)
		angleRad := math.Atan2(float64(k.TargetDirection.X), float64(k.TargetDirection.Z))
		k.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))
		fireDir := rl.Vector3Normalize(k.TargetDirection)
		spawnPos := rl.Vector3Add(k.Position, rl.Vector3{X: 0, Y: 0, Z: 0})
		spawnPos = rl.Vector3Add(spawnPos, rl.Vector3Scale(fireDir, 1.5))
		bulletCmds = append(bulletCmds, BulletCmd{spawnPos, fireDir})
		return bulletCmds, true
	}
	return []BulletCmd{}, false
}

func (k *Killer) endOfHold() []BulletCmd {
	var bulletCmds []BulletCmd
	rl.PlaySound(k.ShotGunSound)
	angleRad := math.Atan2(float64(k.TargetDirection.X), float64(k.TargetDirection.Z))
	k.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))
	fireDir := rl.Vector3Normalize(k.TargetDirection)
	spawnPos := rl.Vector3Add(k.Position, rl.Vector3{X: 0, Y: 0, Z: 0})
	spawnPos = rl.Vector3Add(spawnPos, rl.Vector3Scale(fireDir, 1.5))
	bulletCmds = append(bulletCmds, BulletCmd{spawnPos, fireDir})
	return bulletCmds
}

func (k *Killer) PlanAnimate(dt float32) {
	k.AnimationFrameCounter += k.AnimationFrameSpeed * dt
	anim := k.Animation[k.AnimationIdx]
	for k.AnimationFrameCounter >= 1.0 {
		k.AnimationCurrentFrame++
		k.AnimationFrameCounter -= 1.0

		if k.AnimationIdx == 7 && k.AnimationCurrentFrame >= anim.FrameCount-5 {
			k.AnimationCurrentFrame = anim.FrameCount - 5
			return
		}

		if k.AnimationCurrentFrame >= anim.FrameCount {
			k.AnimationCurrentFrame = 0
		}
	}
}
