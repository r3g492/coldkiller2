package background

import (
	"coldkiller2/killer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Star struct {
	Position rl.Vector3
	Size     float32
	Color    rl.Color
}

var starTex rl.Texture2D
var stars []rl.Vector3

func InitEnvironment() {
	starTex = GenStarTexture()
	for i := 0; i < 300; i++ {
		stars = append(stars, rl.Vector3{
			X: float32(rl.GetRandomValue(-400, 400)),
			Y: float32(rl.GetRandomValue(20, 200)),
			Z: float32(rl.GetRandomValue(-400, 400)),
		})
	}
}

func DrawCleanEnvironment(p *killer.Killer) {
	rl.BeginMode3D(p.Camera)

	for _, sPos := range stars {
		rl.DrawBillboard(p.Camera, starTex, sPos, 0.5, rl.SkyBlue)
	}

	rl.DrawCube(rl.NewVector3(10, -1.5, 20), 0.1, 0.1, 0.1, rl.Gray)

	gridSize := float32(10.0)
	offsetX := float32(int(p.Position.X/gridSize)) * gridSize
	offsetZ := float32(int(p.Position.Z/gridSize)) * gridSize

	rl.PushMatrix()
	rl.Translatef(offsetX, -2.0, offsetZ)
	rl.DrawGrid(100, gridSize)
	rl.PopMatrix()

	rl.DrawCircleV(rl.NewVector2(p.Position.X, p.Position.Z), 150, rl.Fade(rl.Black, 0.5))

	rl.EndMode3D()
}

func DrawMonolith(pos rl.Vector3, width float32, height float32) {
	rl.DrawCube(pos, width, height, width, rl.NewColor(5, 5, 5, 255))
	rl.DrawCubeWires(pos, width, height, width, rl.DarkPurple)
}

func GenStarTexture() rl.Texture2D {
	img := rl.GenImageGradientRadial(16, 16, 0.5, rl.White, rl.Blank)
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return tex
}
