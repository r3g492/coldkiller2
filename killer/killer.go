package killer

import (
	"coldkiller2/bullet"
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

func Init() *Killer {
	playerModel := rl.LoadModel("resources/robot.glb")
	playerAnimation := rl.LoadModelAnimations("resources/robot.glb")
	playerPosition := rl.Vector3{X: 0, Y: 0, Z: 0}
	shotGunSound := util.LoadSoundFromEmbedded("shotgun-03-38220.mp3")
	return &Killer{
		Model:                 playerModel,
		ModelAngleDeg:         0,
		Animation:             playerAnimation,
		AnimationIdx:          0,
		AnimationCurrentFrame: 0,
		AnimationFrameCounter: 0,
		AnimationFrameSpeed:   0.1,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:              playerPosition,
		Size:                  2,
		MoveSpeed:             10.0,
		Camera: rl.Camera3D{
			Position:   rl.Vector3Add(playerPosition, rl.NewVector3(0.0, 10.0, 0.0)),
			Target:     playerPosition,
			Up:         rl.NewVector3(0.0, 0.0, -1),
			Fovy:       30.0,
			Projection: rl.CameraOrthographic,
		},
		ShotGunSound:   shotGunSound,
		ActionTimeLeft: 0,
		Bullets:        make([]bullet.Bullet, 0, 100),
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
	rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, k.Size*2, k.Size*2, k.Size*2, rl.Purple)
	rl.Rotatef(270, 1, 0, 0)
	rl.Rotatef(k.ModelAngleDeg, 0, 1, 0)
	rl.DrawModel(k.Model, rl.NewVector3(0, -k.Size, 0), 0.7, rl.White)
	rl.PopMatrix()
	rl.DrawRay(rl.NewRay(k.Position, k.TargetDirection), rl.Green)
	k.DrawBullets()
}

func (k *Killer) Mutate(input input.Input, dt float32) {
	k.UpdateBullets(dt)

	if k.ActionTimeLeft > 0 {
		k.ActionTimeLeft -= dt
	}
	if k.ActionTimeLeft <= 0 {
		k.movement(input, dt)
	}

	k.Camera = rl.Camera3D{
		Position:   rl.Vector3Add(k.Position, rl.NewVector3(0.0, 10.0, 0.0)),
		Target:     k.Position,
		Up:         rl.NewVector3(0.0, 0.0, -1),
		Fovy:       30.0,
		Projection: rl.CameraOrthographic,
	}
}

func (k *Killer) movement(input input.Input, dt float32) {
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
	if input.Fire {
		rl.PlaySound(k.ShotGunSound)
		move = false

		// 1. Calculate Rotation
		angleRad := math.Atan2(float64(k.TargetDirection.X), float64(k.TargetDirection.Z))
		k.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))

		fireDir := rl.Vector3Normalize(k.TargetDirection)

		spawnPos := rl.Vector3Add(k.Position, rl.Vector3{X: 0, Y: 0, Z: 0})
		spawnPos = rl.Vector3Add(spawnPos, rl.Vector3Scale(fireDir, 1.5))

		k.Bullets = append(k.Bullets, bullet.NewBullet(spawnPos, fireDir))

		k.ActionTimeLeft = 0.2
	}
	if move {
		k.Position = rl.Vector3Add(k.Position, moveAmount)
		angleRad := math.Atan2(float64(k.MoveDirection.X), float64(k.MoveDirection.Z))
		k.ModelAngleDeg = float32(angleRad * (180.0 / math.Pi))
	}
	mouseLocation := input.MouseLocation
	ray := rl.GetScreenToWorldRay(mouseLocation, k.Camera)
	targetOnXzPlane := rl.Vector3{
		X: ray.Position.X,
		Y: 0,
		Z: ray.Position.Z,
	}
	k.TargetDirection = rl.Vector3Subtract(targetOnXzPlane, k.Position)
}

func (k *Killer) UpdateBullets(dt float32) {
	for i := 0; i < len(k.Bullets); i++ {
		movement := rl.Vector3Scale(k.Bullets[i].Direction, k.Bullets[i].Speed*dt)
		k.Bullets[i].Position = rl.Vector3Add(k.Bullets[i].Position, movement)

		k.Bullets[i].LifeTime -= dt

		if k.Bullets[i].LifeTime <= 0 || !k.Bullets[i].Active {
			k.Bullets[i] = k.Bullets[len(k.Bullets)-1]
			k.Bullets = k.Bullets[:len(k.Bullets)-1]
			i--
		}
	}
}

func (k *Killer) DrawBullets() {
	for _, b := range k.Bullets {
		rl.DrawSphere(b.Position, b.Radius, rl.Yellow)
	}
}
