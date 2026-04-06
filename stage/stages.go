package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/model"
	"coldkiller2/sound"
	"coldkiller2/structure"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var aimPracticePresets = []rl.Vector2{
	{X: -7, Y: -7}, {X: 0, Y: -7}, {X: 7, Y: -7},

	{X: -7, Y: 0}, {X: 7, Y: 0},

	{X: -7, Y: 7}, {X: 0, Y: 7}, {X: 7, Y: 7},
}

type Data struct {
	Enemies    []*enemy.Enemy
	Structures []*structure.Structure
}

var Stages []Data

func InitStages() {
	Stages = []Data{
		NewAimPracticeStage(),
		NewAimPracticeStage(),
		NewAimPracticeStage(),
		NewAimPracticeStage(),
		NewAimPracticeStage(),
		NewAimPracticeStage(),
		NewAimPracticeStage(),
		NewAimPracticeStage(),
		NewAimPracticeStage(),
		NewAimPracticeStage(),
	}
}

func NewEnemy(x, z float32) *enemy.Enemy {
	if x < -9 {
		x = -9
	}
	if x > 9 {
		x = 9
	}
	if z < -9 {
		z = -9
	}
	if z > 9 {
		z = 9
	}

	return &enemy.Enemy{
		Model:                 model.UnitV4Model,
		Animation:             model.UnitV4Animation,
		Position:              rl.Vector3{X: x, Y: 0, Z: z},
		Size:                  killer.CharSize,
		MoveSpeed:             4,
		Health:                100,
		AttackRange:           10,
		AimTimeLeft:           1,
		AimTimeUnit:           1,
		FootstepSoundTimeLeft: 0.4,
		FootstepSoundTimeUnit: 0.4,
		FootstepSound:         sound.FootStep,
		AiType:                enemy.Elite,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
	}
}

func GetRandomEnemy() *enemy.Enemy {
	randomIndex := rand.Intn(len(aimPracticePresets))
	pos := aimPracticePresets[randomIndex]

	return NewEnemy(pos.X, pos.Y)
}

func NewAimPracticeStage() Data {
	return Data{
		Enemies:    []*enemy.Enemy{GetRandomEnemy()},
		Structures: DefaultWalls(),
	}
}

func DefaultWalls() []*structure.Structure {
	return []*structure.Structure{
		{Position: rl.Vector3{X: -10, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 10, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 0, Y: 0, Z: 10}, Size: rl.Vector3{X: 20, Y: 1, Z: 1}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 0, Y: 0, Z: -10}, Size: rl.Vector3{X: 20, Y: 1, Z: 1}, Color: rl.DarkGray},
	}
}
