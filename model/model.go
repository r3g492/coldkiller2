package model

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	RobotModel       rl.Model
	RobotAnimation   []rl.ModelAnimation
	SoldierModel     rl.Model
	SoldierAnimation []rl.ModelAnimation
)

func Init() {
	RobotModel = rl.LoadModel("resources/robot.glb")
	RobotAnimation = rl.LoadModelAnimations("resources/robot.glb")
	SoldierModel = rl.LoadModel("resources/soldier.glb")
	SoldierAnimation = rl.LoadModelAnimations("resources/soldier.glb")
}
