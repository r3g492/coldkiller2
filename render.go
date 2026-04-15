package main

import (
	"coldkiller2/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	VirtualWidth  = util.VirtualWidth
	VirtualHeight = util.VirtualHeight
)

var renderTarget rl.RenderTexture2D

func initRender() {
	renderTarget = rl.LoadRenderTexture(VirtualWidth, VirtualHeight)
	rl.SetTextureFilter(renderTarget.Texture, rl.FilterBilinear)
}

func unloadRender() {
	rl.UnloadRenderTexture(renderTarget)
}

func letterbox() (scale, offsetX, offsetY float32) {
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())
	scale = screenW / VirtualWidth
	if sy := screenH / VirtualHeight; sy < scale {
		scale = sy
	}
	offsetX = (screenW - VirtualWidth*scale) / 2
	offsetY = (screenH - VirtualHeight*scale) / 2
	return
}

func beginFrame() {
	rl.BeginTextureMode(renderTarget)
}

func endFrame() {
	rl.EndTextureMode()

	scale, offsetX, offsetY := letterbox()

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	src := rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(renderTarget.Texture.Width),
		Height: -float32(renderTarget.Texture.Height),
	}
	dest := rl.Rectangle{
		X:      offsetX,
		Y:      offsetY,
		Width:  VirtualWidth * scale,
		Height: VirtualHeight * scale,
	}
	rl.DrawTexturePro(renderTarget.Texture, src, dest, rl.Vector2{}, 0, rl.White)
	rl.EndDrawing()
}

func virtualMousePosition() rl.Vector2 {
	scale, offsetX, offsetY := letterbox()
	m := rl.GetMousePosition()
	v := rl.Vector2{
		X: (m.X - offsetX) / scale,
		Y: (m.Y - offsetY) / scale,
	}
	if v.X < 0 {
		v.X = 0
	} else if v.X > VirtualWidth {
		v.X = VirtualWidth
	}
	if v.Y < 0 {
		v.Y = 0
	} else if v.Y > VirtualHeight {
		v.Y = VirtualHeight
	}
	return v
}
