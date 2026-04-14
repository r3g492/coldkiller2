package input

import rl "github.com/gen2brain/raylib-go/raylib"

type Input struct {
	MoveUp         bool
	MoveDown       bool
	MoveLeft       bool
	MoveRight      bool
	MouseLocation  rl.Vector2
	FireHold       bool
	FirePressed    bool
	FireReleased   bool
	ReloadPressed  bool
	EndGamePressed bool
	DashPressed    bool
}

type KeyMap struct {
	Up        int32
	Down      int32
	Left      int32
	Right     int32
	Fire      rl.MouseButton
	Reload    int32
	ResetGame int32
	EndGame   int32
	Dash      int32
}

func DefaultWASD() KeyMap {
	return KeyMap{
		Up:      rl.KeyW,
		Down:    rl.KeyS,
		Left:    rl.KeyA,
		Right:   rl.KeyD,
		Fire:    rl.MouseLeftButton,
		Reload:  rl.KeyR,
		EndGame: rl.KeyEscape,
		Dash:    rl.KeySpace,
	}
}

func ReadInput(keyMap KeyMap) Input {
	return Input{
		MoveUp:         rl.IsKeyDown(keyMap.Up),
		MoveDown:       rl.IsKeyDown(keyMap.Down),
		MoveLeft:       rl.IsKeyDown(keyMap.Left),
		MoveRight:      rl.IsKeyDown(keyMap.Right),
		MouseLocation:  rl.GetMousePosition(),
		FireHold:       rl.IsMouseButtonDown(keyMap.Fire),
		FirePressed:    rl.IsMouseButtonPressed(keyMap.Fire),
		FireReleased:   rl.IsMouseButtonReleased(keyMap.Fire),
		ReloadPressed:  rl.IsKeyDown(keyMap.Reload),
		EndGamePressed: rl.IsKeyPressed(keyMap.EndGame),
		DashPressed:    rl.IsKeyPressed(keyMap.Dash),
	}
}

func GetKeyName(key int32) string {
	switch key {
	case rl.KeyEscape:
		return "ESC"
	case rl.KeySpace:
		return "SPC"
	case rl.KeyW: // Explicitly handle these if GetKeyName returns lowercase/null
		return "W"
	case rl.KeyA:
		return "A"
	case rl.KeyS:
		return "S"
	case rl.KeyD:
		return "D"
	case rl.KeyR:
		return "R"
	default:
		return "?"
	}
}
