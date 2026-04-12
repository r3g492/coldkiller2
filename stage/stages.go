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
		Type3(),
		Type3(),
		Type3(),
		Type3(),
		Type4(),
		Type4(),
		Type4(),
		Type4(),
		Type5(),
		Type5(),
		Type5(),
		Type5(),
		Type6(),
		Type6(),
		Type6(),
		Type6(),
	}
}

func Soldier(x, z float32) *enemy.Enemy {
	return &enemy.Enemy{
		Model:                 model.UnitV4Model,
		ModelRatio:            0.2,
		Animation:             model.UnitV4Animation,
		Position:              rl.Vector3{X: x, Y: 0, Z: z},
		Size:                  killer.CharSize,
		MoveSpeed:             4,
		Health:                100,
		AttackRange:           10,
		AimTimeLeft:           1,
		AimTimeUnit:           1,
		FootstepSoundTimeLeft: 0,
		FootstepSoundTimeUnit: 0.4,
		FootstepSound:         sound.FootStep,
		AiType:                enemy.Elite,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
	}
}

func Dog(x, z float32) *enemy.Enemy {
	return &enemy.Enemy{
		Model:                 model.UnitV3Model,
		ModelRatio:            0.4,
		Animation:             model.UnitV3Animation,
		Position:              rl.Vector3{X: x, Y: 0, Z: z},
		Size:                  killer.CharSize,
		MoveSpeed:             8,
		Health:                100,
		AttackRange:           2,
		AimTimeLeft:           0.5,
		AimTimeUnit:           0.5,
		FootstepSoundTimeLeft: 0,
		FootstepSoundTimeUnit: 0.4,
		FootstepSound:         sound.FootStep,
		AiType:                enemy.SimpleZombie,
		MoveDirection:         rl.Vector3{X: 0, Y: 0, Z: 0},
		TargetDirection:       rl.Vector3{X: 0, Y: 0, Z: 0},
	}
}

type EnemyKind int

const (
	KindDog EnemyKind = iota
	KindSoldier
)

type EnemySpec struct {
	Kind  EnemyKind
	Count int
}

func GetRandomEnemy(radius float32, specs ...EnemySpec) []*enemy.Enemy {
	presets := []rl.Vector2{
		{X: -radius, Y: -radius}, {X: 0, Y: -radius}, {X: radius, Y: -radius},

		{X: -radius, Y: 0}, {X: radius, Y: 0},

		{X: -radius, Y: radius}, {X: 0, Y: radius}, {X: radius, Y: radius},
	}

	total := 0
	for _, s := range specs {
		total += s.Count
	}
	if total > len(presets) {
		total = len(presets)
	}

	indices := rand.Perm(len(presets))[:total]
	enemies := make([]*enemy.Enemy, 0, total)
	idx := 0
	for _, s := range specs {
		factory := Dog
		if s.Kind == KindSoldier {
			factory = Soldier
		}
		for range s.Count {
			if idx >= len(indices) {
				break
			}
			pos := presets[indices[idx]]
			enemies = append(enemies, factory(pos.X, pos.Y))
			idx++
		}
	}
	return enemies
}

func Type1() Data {
	return Data{
		Enemies:    GetRandomEnemy(8, EnemySpec{KindSoldier, 1}),
		Structures: WallType1(),
	}
}

func Type2() Data {
	return Data{
		Enemies:    GetRandomEnemy(8, EnemySpec{KindDog, 1}),
		Structures: WallType1(),
	}
}

func Type3() Data {
	return Data{
		Enemies:    GetRandomEnemy(8, EnemySpec{KindDog, 2}),
		Structures: WallType1(),
	}
}

func Type4() Data {
	return Data{
		Enemies:    GetRandomEnemy(8, EnemySpec{KindDog, 3}),
		Structures: WallType1(),
	}
}

func Type5() Data {
	return Data{
		Enemies: GetRandomEnemy(
			15,
			EnemySpec{KindDog, 1},
			EnemySpec{KindSoldier, 1},
		),
		Structures: WallType2(),
	}
}

func Type6() Data {
	var enemies []*enemy.Enemy
	enemies = GetRandomEnemy(
		15,
		EnemySpec{KindDog, 1},
		EnemySpec{KindSoldier, 1},
	)
	enemies = append(
		enemies,
		GetRandomEnemy(
			25,
			EnemySpec{KindDog, 0},
			EnemySpec{KindSoldier, 2},
		)...,
	)

	return Data{
		Enemies:    enemies,
		Structures: WallType3(),
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

func WallType2() []*structure.Structure {
	return []*structure.Structure{
		{Position: rl.Vector3{X: -5, Y: 0, Z: -5}, Size: rl.Vector3{X: 2, Y: 2, Z: 2}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 5, Y: 0, Z: -5}, Size: rl.Vector3{X: 2, Y: 2, Z: 2}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: -5, Y: 0, Z: 5}, Size: rl.Vector3{X: 2, Y: 2, Z: 2}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 5, Y: 0, Z: 5}, Size: rl.Vector3{X: 2, Y: 2, Z: 2}, Color: rl.DarkGray},
	}
}

func WallType3() []*structure.Structure {
	return []*structure.Structure{
		{Position: rl.Vector3{X: -8, Y: 0, Z: -8}, Size: rl.Vector3{X: 3, Y: 3, Z: 3}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 8, Y: 0, Z: -8}, Size: rl.Vector3{X: 3, Y: 3, Z: 3}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: -8, Y: 0, Z: 8}, Size: rl.Vector3{X: 3, Y: 3, Z: 3}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 8, Y: 0, Z: 8}, Size: rl.Vector3{X: 3, Y: 3, Z: 3}, Color: rl.DarkGray},
	}
}
