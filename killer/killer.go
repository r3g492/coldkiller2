package killer

import (
	"coldkiller2/input"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Killer struct {
	Model                 rl.Model
	ModelDirection        rl.Vector3
	Animation             []rl.ModelAnimation
	AnimationIdx          int
	AnimationCurrentFrame int32
	AnimationFrameCounter float32
	AnimationFrameSpeed   float32
	Direction             rl.Vector3
	Position              rl.Vector3
	Size                  float32
	MoveSpeed             float32
}

func Init() *Killer {
	// playerAnimation setting
	playerModel := rl.LoadModel("resources/robot.glb")
	playerAnimation := rl.LoadModelAnimations("resources/robot.glb")
	return &Killer{
		Model:                 playerModel,
		ModelDirection:        rl.Vector3{X: 0, Y: 0, Z: 0},
		Animation:             playerAnimation,
		AnimationIdx:          0,
		AnimationCurrentFrame: 0,
		AnimationFrameCounter: 0,
		AnimationFrameSpeed:   0.1,
		Direction:             rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:              rl.Vector3{X: 0, Y: 0, Z: 0},
		Size:                  2,
		MoveSpeed:             10.0,
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

	yaw := float32(math.Atan2(float64(k.ModelDirection.X), float64(k.ModelDirection.Z))) * (180.0 / math.Pi)
	xzDist := math.Sqrt(float64(k.ModelDirection.X*k.ModelDirection.X) + float64(k.ModelDirection.Z*k.ModelDirection.Z))
	pitch := float32(math.Atan2(float64(k.ModelDirection.Y), xzDist)) * (180.0 / math.Pi)
	rl.Rotatef(yaw, 0, 1, 0)
	rl.Rotatef(-pitch, 1, 0, 0)

	rl.DrawModel(k.Model, rl.NewVector3(0, -k.Size, 0), 0.7, rl.White)
	rl.PopMatrix()
}

func (k *Killer) GetCamera() rl.Camera3D {
	return rl.Camera3D{
		Position:   rl.Vector3Add(k.Position, rl.NewVector3(0.0, 10.0, 0.0)),
		Target:     k.Position,
		Up:         rl.NewVector3(0.0, 0.0, -1),
		Fovy:       30.0,
		Projection: rl.CameraOrthographic,
	}
}

func (k *Killer) Control(input input.Input, dt float32) {
	k.Direction = rl.Vector3{}
	if input.MoveUp {
		k.Direction.Z -= 1
	}
	if input.MoveDown {
		k.Direction.Z += 1
	}
	if input.MoveLeft {
		k.Direction.X -= 1
	}
	if input.MoveRight {
		k.Direction.X += 1
	}
	if rl.Vector3LengthSqr(k.Direction) > 0 {
		k.Direction = rl.Vector3Normalize(k.Direction)
	}
	moveAmount := rl.Vector3Scale(k.Direction, k.MoveSpeed*dt)
	k.Position = rl.Vector3Add(k.Position, moveAmount)
}
