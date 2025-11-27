package input

import rl "github.com/gen2brain/raylib-go/raylib"

type Input struct {
	MoveUp        bool
	MoveDown      bool
	MoveLeft      bool
	MoveRight     bool
	MouseLocation rl.Vector2
	Fire          bool
}

type KeyMap struct {
	Up    int32
	Down  int32
	Left  int32
	Right int32
	Fire  rl.MouseButton
}

func DefaultWASD() KeyMap {
	return KeyMap{
		Up:    rl.KeyW,
		Down:  rl.KeyS,
		Left:  rl.KeyA,
		Right: rl.KeyD,
		Fire:  rl.MouseLeftButton,
	}
}

func ReadInput(keyMap KeyMap) Input {
	return Input{
		MoveUp:        rl.IsKeyDown(keyMap.Up),
		MoveDown:      rl.IsKeyDown(keyMap.Down),
		MoveLeft:      rl.IsKeyDown(keyMap.Left),
		MoveRight:     rl.IsKeyDown(keyMap.Right),
		MouseLocation: rl.GetMousePosition(),
		Fire:          rl.IsMouseButtonPressed(keyMap.Fire),
	}
}
