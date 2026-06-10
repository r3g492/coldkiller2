package model

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"coldkiller2/util"
)

var (
	BombRobotModel       rl.Model
	BombRobotAnimation   []rl.ModelAnimation
	GunRobotModel        rl.Model
	GunRobotAnimation    []rl.ModelAnimation
	ChargeRobotModel     rl.Model
	ChargeRobotAnimation []rl.ModelAnimation
	PlayerModel          rl.Model
	PlayerAnimation      []rl.ModelAnimation
)

func Init() {
	BombRobotModel, BombRobotAnimation = util.LoadModelFromEmbedded("bombRobot.glb")
	GunRobotModel, GunRobotAnimation = util.LoadModelFromEmbedded("gunRobot.glb")
	ChargeRobotModel, ChargeRobotAnimation = util.LoadModelFromEmbedded("chargeRobot.glb")
	PlayerModel, PlayerAnimation = util.LoadModelFromEmbedded("playerRobot.glb")
}
