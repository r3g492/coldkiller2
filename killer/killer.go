package killer

import (
	"coldkiller2/input"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Killer struct {
	model                 rl.Model
	ModelDirection        rl.Vector3
	animation             []rl.ModelAnimation
	animationIdx          int
	animationCurrentFrame int32
	animationFrameCounter float32
	animationFrameSpeed   float32
	direction             rl.Vector3
	Position              rl.Vector3
	Size                  float32
}

func Init() *Killer {
	// playerAnimation setting
	playerModel := rl.LoadModel("resources/robot.glb")
	playerAnimation := rl.LoadModelAnimations("resources/robot.glb")
	return &Killer{
		model:                 playerModel,
		ModelDirection:        rl.Vector3{X: 0, Y: 0, Z: 0},
		animation:             playerAnimation,
		animationIdx:          0,
		animationCurrentFrame: 0,
		animationFrameCounter: 0,
		animationFrameSpeed:   0.1,
		direction:             rl.Vector3{X: 0, Y: 0, Z: 0},
		Position:              rl.Vector3{X: 0, Y: 0, Z: 0},
		Size:                  2,
	}
}

func (k *Killer) Unload() {
	rl.UnloadModel(k.model)
	rl.UnloadModelAnimations(k.animation)
}

func (k *Killer) Draw3D() {
	anim := k.animation[k.animationIdx]
	rl.UpdateModelAnimation(k.model, anim, k.animationCurrentFrame)

	rl.PushMatrix()
	rl.Translatef(k.Position.X, k.Position.Y, k.Position.Z)
	rl.DrawCubeWires(rl.Vector3{X: 0, Y: 0, Z: 0}, k.Size*2, k.Size*2, k.Size*2, rl.Purple)

	yaw := float32(math.Atan2(float64(k.ModelDirection.X), float64(k.ModelDirection.Z))) * (180.0 / math.Pi)
	xzDist := math.Sqrt(float64(k.ModelDirection.X*k.ModelDirection.X) + float64(k.ModelDirection.Z*k.ModelDirection.Z))
	pitch := float32(math.Atan2(float64(k.ModelDirection.Y), xzDist)) * (180.0 / math.Pi)
	rl.Rotatef(yaw, 0, 1, 0)
	rl.Rotatef(-pitch, 1, 0, 0)

	rl.DrawModel(k.model, rl.NewVector3(0, -k.Size, 0), 0.7, rl.White)
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

	if input.MoveUp {
		k.Position.Z -= 0.1
	}

	if input.MoveDown {
		k.Position.Z += 0.1
	}

	if input.MoveLeft {
		k.Position.X -= 0.1
	}

	if input.MoveRight {
		k.Position.X += 0.1
	}
}
