package model

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	UnitV3Model     rl.Model
	UnitV3Animation []rl.ModelAnimation
	UnitV4Model     rl.Model
	UnitV4Animation []rl.ModelAnimation
)

func Init() {
	UnitV3Model = rl.LoadModel("resources/unit_v3.glb")
	UnitV3Animation = rl.LoadModelAnimations("resources/unit_v3.glb")
	UnitV4Model = rl.LoadModel("resources/unit_v4.glb")
	UnitV4Animation = rl.LoadModelAnimations("resources/unit_v4.glb")
}
