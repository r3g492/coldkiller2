package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/model"
	"coldkiller2/sound"
	"coldkiller2/structure"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Data struct {
	Enemies    []*enemy.Enemy
	Structures []*structure.Structure
}

var Stages []Data

func InitStages() {
	Stages = []Data{
		// --- STAGE 0 ---
		{
			Enemies: []*enemy.Enemy{
				{
					Model:                 model.UnitV4Model,
					ModelAngleDeg:         0,
					Animation:             model.UnitV4Animation,
					MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
					TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
					Position:              rl.Vector3{X: 5, Y: 0, Z: 5},
					Size:                  killer.CharSize,
					MoveSpeed:             4,
					ActionTimeLeft:        0,
					Health:                100,
					ShouldBeDeleted:       false,
					AttackRange:           10,
					AimTimeLeft:           1,
					AimTimeUnit:           1,
					FootstepSoundTimeLeft: 0.4,
					FootstepSoundTimeUnit: 0.4,
					FootstepSound:         sound.FootStep,
					AiType:                enemy.Elite,
				},
			},
			Structures: []*structure.Structure{
				{Position: rl.Vector3{X: -10, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 0, Y: 0, Z: 0}, Color: rl.DarkGray},
				{Position: rl.Vector3{X: 10, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 0, Y: 0, Z: 0}, Color: rl.DarkGray},
				{Position: rl.Vector3{X: 0, Y: 0, Z: 10}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 1, Y: 0, Z: 0}, Color: rl.DarkGray},
				{Position: rl.Vector3{X: 0, Y: 0, Z: -10}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 1, Y: 0, Z: 0}, Color: rl.DarkGray},
			},
		},

		// --- STAGE 1 ---
		{
			Enemies: []*enemy.Enemy{
				{
					Model:                 model.UnitV4Model,
					ModelAngleDeg:         0,
					Animation:             model.UnitV4Animation,
					Position:              rl.Vector3{X: 5, Y: 0, Z: 5},
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
				},
				{
					Model:                 model.UnitV4Model,
					ModelAngleDeg:         0,
					Animation:             model.UnitV4Animation,
					Position:              rl.Vector3{X: -5, Y: 0, Z: 5},
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
				},
			},
			Structures: []*structure.Structure{
				{Position: rl.Vector3{X: -10, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 0, Y: 0, Z: 0}, Color: rl.DarkGray},
				{Position: rl.Vector3{X: 10, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 0, Y: 0, Z: 0}, Color: rl.DarkGray},
				{Position: rl.Vector3{X: 0, Y: 0, Z: 10}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 1, Y: 0, Z: 0}, Color: rl.DarkGray},
				{Position: rl.Vector3{X: 0, Y: 0, Z: -10}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 1, Y: 0, Z: 0}, Color: rl.DarkGray},
			},
		},

		// --- STAGE 2 ---
		{
			Enemies: []*enemy.Enemy{
				{
					Model:                 model.UnitV4Model,
					Position:              rl.Vector3{X: 5, Y: 0, Z: 5},
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
				},
				{
					Model:                 model.UnitV4Model,
					Position:              rl.Vector3{X: -5, Y: 0, Z: 5},
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
				},
			},
			Structures: []*structure.Structure{
				{Position: rl.Vector3{X: -10, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 0, Y: 0, Z: 0}, Color: rl.DarkGray},
				{Position: rl.Vector3{X: 10, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 0, Y: 0, Z: 0}, Color: rl.DarkGray},
				{Position: rl.Vector3{X: 0, Y: 0, Z: 10}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 1, Y: 0, Z: 0}, Color: rl.DarkGray},
				{Position: rl.Vector3{X: 0, Y: 0, Z: -10}, Size: rl.Vector3{X: 1, Y: 1, Z: 20}, Direction: rl.Vector3{X: 1, Y: 0, Z: 0}, Color: rl.DarkGray},
			},
		},
	}
}
