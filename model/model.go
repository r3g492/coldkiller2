package model

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	UnitV4Model     rl.Model
	UnitV4Animation []rl.ModelAnimation
)

func Init() {
	UnitV4Model = rl.LoadModel("resources/unit_v4.glb")
	UnitV4Animation = rl.LoadModelAnimations("resources/unit_v4.glb")
}
