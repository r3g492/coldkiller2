package stage

import (
	"coldkiller2/enemy"
	"coldkiller2/structure"
	"embed"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed data/stages.json
var stageFS embed.FS

const stageDataFile = "data/stages.json"

// EnemyKind selects which enemy factory a placed enemy is built from.
type EnemyKind string

const (
	KindRobot      EnemyKind = "robot"
	KindSoldier    EnemyKind = "soldier"
	KindSniper     EnemyKind = "sniper"
	KindCharger    EnemyKind = "charger"
	KindSuperRobot EnemyKind = "super_robot"
	KindRival      EnemyKind = "rival"
	KindRedRival   EnemyKind = "red_rival"
	KindGoldRival  EnemyKind = "gold_rival"
)

// EnemyKinds lists every placeable kind in a stable order (used by the editor).
var EnemyKinds = []EnemyKind{KindRobot, KindSoldier, KindSniper, KindCharger, KindSuperRobot, KindRival, KindRedRival, KindGoldRival}

var enemyFactories = map[EnemyKind]func(x, z float32) *enemy.Enemy{
	KindRobot:      enemy.Robot,
	KindSoldier:    enemy.Soldier,
	KindSniper:     enemy.Sniper,
	KindCharger:    enemy.ChargerRobot,
	KindSuperRobot: enemy.SuperRobot,
	KindRival:      enemy.Rival,
	KindRedRival:   enemy.RedRival,
	KindGoldRival:  enemy.GoldRival,
}

// Point is an XZ position on the ground plane.
type Point struct {
	X float32 `json:"x"`
	Z float32 `json:"z"`
}

// EnemySpec is one authored enemy placement.
type EnemySpec struct {
	Kind EnemyKind `json:"kind"`
	X    float32   `json:"x"`
	Z    float32   `json:"z"`
}

// StructureKind selects a structure preset (default size + color).
type StructureKind string

const (
	StructWall   StructureKind = "wall"
	StructBlock  StructureKind = "block"
	StructPillar StructureKind = "pillar"
	StructCrate  StructureKind = "crate"
	StructBunker StructureKind = "bunker"
)

// StructureKinds lists every placeable structure preset in a stable order.
var StructureKinds = []StructureKind{StructWall, StructBlock, StructPillar, StructCrate, StructBunker}

type structureType struct {
	Size  rl.Vector3
	Color rl.Color
}

var structureTypes = map[StructureKind]structureType{
	StructWall:   {Size: rl.Vector3{X: 6, Y: 2, Z: 1}, Color: rl.DarkGray},
	StructBlock:  {Size: rl.Vector3{X: 2, Y: 2, Z: 2}, Color: rl.Gray},
	StructPillar: {Size: rl.Vector3{X: 1, Y: 4, Z: 1}, Color: rl.Beige},
	StructCrate:  {Size: rl.Vector3{X: 1.5, Y: 1.5, Z: 1.5}, Color: rl.Brown},
	StructBunker: {Size: rl.Vector3{X: 5, Y: 3, Z: 5}, Color: rl.DarkGreen},
}

// StructureDef returns the default size and color for a structure kind, falling
// back to a plain gray block for unknown kinds.
func StructureDef(kind StructureKind) (rl.Vector3, rl.Color) {
	if t, ok := structureTypes[kind]; ok {
		return t.Size, t.Color
	}
	return rl.Vector3{X: 2, Y: 2, Z: 2}, rl.DarkGray
}

// StructureSpec is one authored wall/obstacle.
type StructureSpec struct {
	Kind      StructureKind `json:"kind"`
	Position  rl.Vector3    `json:"position"`
	Size      rl.Vector3    `json:"size"`
	Direction rl.Vector3    `json:"direction"`
}

// StageSpec is the authored content of a single stage.
type StageSpec struct {
	Start      Point           `json:"start"`
	Enemies    []EnemySpec     `json:"enemies"`
	Structures []StructureSpec `json:"structures"`
}

// Stages holds every authored stage in play order. It is loaded by InitStages
// and edited in place by the in-game editor.
var Stages []StageSpec

// InitStages loads the authored stages from the embedded JSON data.
func InitStages() {
	data, err := stageFS.ReadFile(stageDataFile)
	if err != nil {
		log.Fatalf("stage: failed to read %s: %v", stageDataFile, err)
	}
	if err := json.Unmarshal(data, &Stages); err != nil {
		log.Fatalf("stage: failed to parse %s: %v", stageDataFile, err)
	}
}

// Save writes the current Stages back to the on-disk source JSON. The path is
// resolved relative to this source file so it works regardless of the working
// directory during development. Shipped builds run off the embedded copy and are
// not expected to call Save.
func Save() error {
	out, err := json.MarshalIndent(Stages, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sourceDataPath(), out, 0o644)
}

func sourceDataPath() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), stageDataFile)
}

// Instantiate builds live enemies and structures from the authored spec.
func (s StageSpec) Instantiate() ([]*enemy.Enemy, []*structure.Structure) {
	enemies := make([]*enemy.Enemy, 0, len(s.Enemies))
	for _, es := range s.Enemies {
		if factory, ok := enemyFactories[es.Kind]; ok {
			enemies = append(enemies, factory(es.X, es.Z))
		}
	}
	structures := make([]*structure.Structure, 0, len(s.Structures))
	for _, ss := range s.Structures {
		_, color := StructureDef(ss.Kind)
		structures = append(structures, &structure.Structure{
			Position:  ss.Position,
			Size:      ss.Size,
			Direction: ss.Direction,
			Color:     color,
		})
	}
	return enemies, structures
}
