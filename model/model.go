package model

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"coldkiller2/util"
)

var (
	RobotModel       rl.Model
	RobotAnimation   []rl.ModelAnimation
	SoldierModel     rl.Model
	SoldierAnimation []rl.ModelAnimation
	KillerModel      rl.Model
	KillerAnimation  []rl.ModelAnimation
)

func Init() {
	RobotModel, RobotAnimation = util.LoadModelFromEmbedded("robot.glb")
	SoldierModel, SoldierAnimation = util.LoadModelFromEmbedded("soldier.glb")
	KillerModel, KillerAnimation = util.LoadModelFromEmbedded("killer.glb")
}
