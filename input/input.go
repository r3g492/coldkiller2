package input

import rl "github.com/gen2brain/raylib-go/raylib"

type Input struct {
	MoveUp         bool
	MoveDown       bool
	MoveLeft       bool
	MoveRight      bool
	MouseLocation  rl.Vector2
	FireDown       bool
	FirePressed    bool
	FireReleased   bool
	ReloadPressed  bool
	EndGamePressed bool
	SlowTimeDown   bool
}

type KeyMap struct {
	Up        int32          `json:"up"`
	Down      int32          `json:"down"`
	Left      int32          `json:"left"`
	Right     int32          `json:"right"`
	Fire      rl.MouseButton `json:"fire"`
	Reload    int32          `json:"reload"`
	ResetGame int32          `json:"reset_game"`
	EndGame   int32          `json:"end_game"`
	SlowTime  int32          `json:"slow_time"`
}

type Bindable struct {
	Name string
	Get  func(*KeyMap) int32
	Set  func(*KeyMap, int32)
}

func Bindables() []Bindable {
	return []Bindable{
		{"Up", func(k *KeyMap) int32 { return k.Up }, func(k *KeyMap, v int32) { k.Up = v }},
		{"Left", func(k *KeyMap) int32 { return k.Left }, func(k *KeyMap, v int32) { k.Left = v }},
		{"Down", func(k *KeyMap) int32 { return k.Down }, func(k *KeyMap, v int32) { k.Down = v }},
		{"Right", func(k *KeyMap) int32 { return k.Right }, func(k *KeyMap, v int32) { k.Right = v }},
		{"Reload", func(k *KeyMap) int32 { return k.Reload }, func(k *KeyMap, v int32) { k.Reload = v }},
		{"Slow", func(k *KeyMap) int32 { return k.SlowTime }, func(k *KeyMap, v int32) { k.SlowTime = v }},
	}
}

func DefaultWASD() KeyMap {
	return KeyMap{
		Up:       rl.KeyW,
		Down:     rl.KeyS,
		Left:     rl.KeyA,
		Right:    rl.KeyD,
		Fire:     rl.MouseLeftButton,
		Reload:   rl.KeyR,
		EndGame:  rl.KeyEscape,
		SlowTime: rl.KeySpace,
	}
}

func ReadInput(keyMap KeyMap) Input {
	return Input{
		MoveUp:         rl.IsKeyDown(keyMap.Up),
		MoveDown:       rl.IsKeyDown(keyMap.Down),
		MoveLeft:       rl.IsKeyDown(keyMap.Left),
		MoveRight:      rl.IsKeyDown(keyMap.Right),
		MouseLocation:  rl.GetMousePosition(),
		FireDown:       rl.IsMouseButtonDown(keyMap.Fire),
		FirePressed:    rl.IsMouseButtonPressed(keyMap.Fire),
		FireReleased:   rl.IsMouseButtonReleased(keyMap.Fire),
		ReloadPressed:  rl.IsKeyPressed(keyMap.Reload),
		EndGamePressed: rl.IsKeyPressed(keyMap.EndGame),
		SlowTimeDown:   rl.IsKeyDown(keyMap.SlowTime),
	}
}

func GetKeyName(key int32) string {
	switch key {
	case rl.KeyEscape:
		return "ESC"
	case rl.KeySpace:
		return "SPC"
	case rl.KeyEnter, rl.KeyKpEnter:
		return "ENT"
	case rl.KeyTab:
		return "TAB"
	case rl.KeyLeftShift, rl.KeyRightShift:
		return "SHFT"
	case rl.KeyLeftControl, rl.KeyRightControl:
		return "CTRL"
	case rl.KeyLeftAlt, rl.KeyRightAlt:
		return "ALT"
	case rl.KeyUp:
		return "UP"
	case rl.KeyDown:
		return "DN"
	case rl.KeyLeft:
		return "LT"
	case rl.KeyRight:
		return "RT"
	}
	if key >= rl.KeyA && key <= rl.KeyZ {
		return string(rune('A' + (key - rl.KeyA)))
	}
	if key >= rl.KeyZero && key <= rl.KeyNine {
		return string(rune('0' + (key - rl.KeyZero)))
	}
	return "?"
}
