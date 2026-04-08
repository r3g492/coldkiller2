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

type Data struct {
	Enemies    []*enemy.Enemy
	Structures []*structure.Structure
}

var Stages []Data

func InitStages() {
	Stages = []Data{
		Type1(),
		Type1(),
		Type1(),
		Type1(),
		Type2(),
		Type2(),
		Type2(),
		Type2(),
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

func GetRandomEnemy(
	radius float32,
	howManyEnemies int,
) []*enemy.Enemy {
	presets := []rl.Vector2{
		{X: -radius, Y: -radius}, {X: 0, Y: -radius}, {X: radius, Y: -radius},

		{X: -radius, Y: 0}, {X: radius, Y: 0},

		{X: -radius, Y: radius}, {X: 0, Y: radius}, {X: radius, Y: radius},
	}

	if howManyEnemies > len(presets) {
		howManyEnemies = len(presets)
	}

	indices := rand.Perm(len(presets))[:howManyEnemies]
	enemies := make([]*enemy.Enemy, howManyEnemies)
	for i, idx := range indices {
		pos := presets[idx]
		enemies[i] = NewEnemy(pos.X, pos.Y)
	}
	return enemies
}

func Type1() Data {
	return Data{
		Enemies:    GetRandomEnemy(8, 1),
		Structures: WallType1(),
	}
}

func WallType1() []*structure.Structure {
	return []*structure.Structure{
		{Position: rl.Vector3{X: -15, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 30}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 15, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 30}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 0, Y: 0, Z: 15}, Size: rl.Vector3{X: 30, Y: 1, Z: 1}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 0, Y: 0, Z: -15}, Size: rl.Vector3{X: 30, Y: 1, Z: 1}, Color: rl.DarkGray},
	}
}

func Type2() Data {
	return Data{
		Enemies:    GetRandomEnemy(15, 2),
		Structures: WallType2(),
	}
}

func WallType2() []*structure.Structure {
	return []*structure.Structure{
		{Position: rl.Vector3{X: -5, Y: 0, Z: -5}, Size: rl.Vector3{X: 2, Y: 2, Z: 2}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 5, Y: 0, Z: -5}, Size: rl.Vector3{X: 2, Y: 2, Z: 2}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: -5, Y: 0, Z: 5}, Size: rl.Vector3{X: 2, Y: 2, Z: 2}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 5, Y: 0, Z: 5}, Size: rl.Vector3{X: 2, Y: 2, Z: 2}, Color: rl.DarkGray},
	}
}
