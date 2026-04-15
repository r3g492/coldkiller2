package main

import (
	"coldkiller2/blast"
	"coldkiller2/bullet"
	"coldkiller2/enemy"
	"coldkiller2/killer"
	"coldkiller2/stage"
	"coldkiller2/structure"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func log(mouseLocation rl.Vector2, dt float32, player *killer.Killer) {
	if time.Since(lastLog) >= 1000*time.Millisecond {
		msg := fmt.Sprintf("mouseLocation=%v, dt=%v", mouseLocation, dt)
		fmt.Println(msg)
		fmt.Println(player.MoveDirection)
		lastLog = time.Now()
	}
}

func unloadGame(
	bulletManager *bullet.Manager,
	blastManager *blast.Manager,
	structureManager *structure.Manager,
	player *killer.Killer,
	enemyManager *enemy.Manager,
	stageManager *stage.Manager,
) {
	bulletManager.Unload()
	blastManager.Unload()
	structureManager.Unload()
	player.Unload()
	enemyManager.Unload()
	stageManager.Unload()
}

func generateFloorTile() rl.Texture2D {
	const sz = int32(64)

	base := rl.NewColor(48, 48, 56, 255)
	bevelLight := rl.NewColor(104, 104, 112, 255)
	bevelDark := rl.NewColor(24, 24, 32, 255)
	innerLight := rl.NewColor(64, 64, 72, 255)
	innerDark := rl.NewColor(32, 32, 40, 255)
	checker := rl.NewColor(52, 52, 61, 255)

	img := rl.GenImageColor(int(sz), int(sz), base)

	for y := int32(2); y < sz-2; y++ {
		for x := int32(2); x < sz-2; x++ {
			if (x+y)%4 < 2 {
				rl.ImageDrawPixel(img, x, y, checker)
			}
		}
	}

	for i := int32(0); i < sz; i++ {
		rl.ImageDrawPixel(img, i, 0, bevelLight)
		rl.ImageDrawPixel(img, 0, i, bevelLight)
		rl.ImageDrawPixel(img, i, sz-1, bevelDark)
		rl.ImageDrawPixel(img, sz-1, i, bevelDark)
	}

	for i := int32(1); i < sz-1; i++ {
		rl.ImageDrawPixel(img, i, 1, innerLight)
		rl.ImageDrawPixel(img, 1, i, innerLight)
		rl.ImageDrawPixel(img, i, sz-2, innerDark)
		rl.ImageDrawPixel(img, sz-2, i, innerDark)
	}

	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	rl.SetTextureWrap(tex, rl.WrapRepeat)
	return tex
}

func initNewGame(
	bulletManager *bullet.Manager,
	blastManager *blast.Manager,
	structureManager *structure.Manager,
	player *killer.Killer,
	enemyManager *enemy.Manager,
	stageManager *stage.Manager,
) {
	bulletManager.Init()
	blastManager.Init()
	structureManager.Init()
	player.Init()
	enemyManager.Init(player)
	stageManager.Init(
		structureManager,
		enemyManager,
		player,
	)
	stage.InitStages()
}
