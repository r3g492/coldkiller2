package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/structure"
	"math"
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

type EnemyKind int

const (
	KindRobot EnemyKind = iota
	KindSoldier
)

type EnemySpec struct {
	Kind  EnemyKind
	Count int
}

const enemySpawnSpacing = float32(8)

func GetRandomEnemy(radius float32, structures []*structure.Structure, existing []*enemy.Enemy, specs ...EnemySpec) []*enemy.Enemy {
	count := int(2 * math.Pi * float64(radius) / float64(enemySpawnSpacing))
	if count < 1 {
		count = 1
	}
	angleStep := 2 * math.Pi / float64(count)
	angleOffset := rand.Float64() * 2 * math.Pi // random rotation so positions vary each init

	enemySize := rl.Vector3{X: killer.CharSize, Y: killer.CharSize, Z: killer.CharSize}
	presets := make([]rl.Vector2, 0, count)
	for i := 0; i < count; i++ {
		angle := angleOffset + float64(i)*angleStep
		pos := rl.Vector3{
			X: float32(math.Cos(angle)) * radius,
			Y: 0,
			Z: float32(math.Sin(angle)) * radius,
		}
		overlaps := false
		for _, s := range structures {
			if s.CheckCollision(pos, pos, enemySize) {
				overlaps = true
				break
			}
		}
		if !overlaps {
			for _, e := range existing {
				if rl.Vector3Distance(pos, e.Position) < enemySpawnSpacing {
					overlaps = true
					break
				}
			}
		}
		if !overlaps {
			presets = append(presets, rl.Vector2{X: pos.X, Y: pos.Z})
		}
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
		factory := enemy.Robot
		if s.Kind == KindSoldier {
			factory = enemy.Soldier
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
	structs := WallType1()
	return Data{
		Enemies:    GetRandomEnemy(17, structs, nil, EnemySpec{KindSoldier, 1}),
		Structures: structs,
	}
}

func Type2() Data {
	structs := WallType1()
	return Data{
		Enemies:    GetRandomEnemy(17, structs, nil, EnemySpec{KindRobot, 1}),
		Structures: structs,
	}
}

func Type3() Data {
	structs := WallType1()
	return Data{
		Enemies:    GetRandomEnemy(17, structs, nil, EnemySpec{KindRobot, 2}),
		Structures: structs,
	}
}

func Type4() Data {
	structs := WallType1()
	return Data{
		Enemies:    GetRandomEnemy(17, structs, nil, EnemySpec{KindRobot, 3}),
		Structures: structs,
	}
}

func Type5() Data {
	structs := WallType2()
	return Data{
		Enemies: GetRandomEnemy(
			17, structs, nil,
			EnemySpec{KindRobot, 3},
			EnemySpec{KindSoldier, 1},
		),
		Structures: structs,
	}
}

func Type6() Data {
	structs := WallType3()
	enemies := GetRandomEnemy(
		20, structs, nil,
		EnemySpec{KindRobot, 1},
		EnemySpec{KindSoldier, 1},
	)
	enemies = append(
		enemies,
		GetRandomEnemy(
			25, structs, enemies,
			EnemySpec{KindRobot, 3},
			EnemySpec{KindSoldier, 1},
		)...,
	)

	return Data{
		Enemies:    enemies,
		Structures: structs,
	}
}

func WallType1() []*structure.Structure {
	return []*structure.Structure{
		{Position: rl.Vector3{X: -20, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 40}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 20, Y: 0, Z: 0}, Size: rl.Vector3{X: 1, Y: 1, Z: 40}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 0, Y: 0, Z: 20}, Size: rl.Vector3{X: 40, Y: 1, Z: 1}, Color: rl.DarkGray},
		{Position: rl.Vector3{X: 0, Y: 0, Z: -20}, Size: rl.Vector3{X: 40, Y: 1, Z: 1}, Color: rl.DarkGray},
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
