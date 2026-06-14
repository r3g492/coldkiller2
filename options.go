package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"coldkiller2/input"
	"coldkiller2/killer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type WindowMode int

const (
	WindowModeBorderless WindowMode = iota
	WindowModeWindowed
)

type Resolution struct {
	Width  int
	Height int
	Label  string
}

var presetResolutions = []Resolution{
	{1280, 720, "1280x720"},
	{1920, 1080, "1920x1080"},
	{2560, 1440, "2560x1440"},
}

type PlayerStats struct {
	MoveSpeed            float32 `json:"move_speed"`
	AmmoCapacity         int32   `json:"ammo_capacity"`
	Range                float32 `json:"range"`
	SlowTimeDuration     float32 `json:"slow_time_duration"`
	SlowTimeRechargeRate float32 `json:"slow_time_recharge_rate"`
	ReloadTimeUnit       float32 `json:"reload_time_unit"`
	Health               int32   `json:"health"`
	Ammo                 int32   `json:"ammo"`
}

type GameConfig struct {
	WindowMode  WindowMode `json:"window_mode"`
	ResWidth    int        `json:"res_width"`
	ResHeight   int        `json:"res_height"`
	GameCleared bool       `json:"game_cleared"`
	BestKps     float32    `json:"best_kps"`
	BestCombo   int        `json:"best_combo"`

	StageLevel  int         `json:"stage_level"`
	PlayerStats PlayerStats `json:"player_stats"`

	KeyBindings *input.KeyMap `json:"key_bindings,omitempty"`
}

var currentConfig = GameConfig{
	WindowMode: WindowModeBorderless,
	ResWidth:   1280,
	ResHeight:  720,
}

const (
	saveDirName  = "KillPerSecondGameSaves"
	saveFileName = "save.dat"
)

func savePath() string {
	base, err := os.UserConfigDir()
	if err != nil {
		return saveFileName
	}
	return filepath.Join(base, saveDirName, saveFileName)
}

func loadConfig() {
	data, err := os.ReadFile(savePath())
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &currentConfig)
}

func saveConfig() {
	path := savePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return
	}
	data, _ := json.Marshal(currentConfig)
	_ = os.WriteFile(path, data, 0644)
}

func saveKeyBindings() {
	km := keyMap
	currentConfig.KeyBindings = &km
	saveConfig()
}

func saveProgress(stageLevel int, p *killer.Killer) {
	currentConfig.StageLevel = stageLevel
	currentConfig.PlayerStats = PlayerStats{
		MoveSpeed:            p.MoveSpeed,
		AmmoCapacity:         p.AmmoCapacity,
		Range:                p.Range,
		SlowTimeDuration:     p.SlowTimeDuration,
		SlowTimeRechargeRate: p.SlowTimeRechargeRate,
		ReloadTimeUnit:       p.ReloadTimeUnit,
		Health:               p.Health,
		Ammo:                 p.Ammo,
	}
	saveConfig()
}

func loadProgress(p *killer.Killer) {
	s := currentConfig.PlayerStats
	p.MoveSpeed = s.MoveSpeed
	p.AmmoCapacity = s.AmmoCapacity
	p.Range = s.Range
	p.SlowTimeDuration = s.SlowTimeDuration
	p.SlowTimeRechargeRate = s.SlowTimeRechargeRate
	p.ReloadTimeUnit = s.ReloadTimeUnit
	p.SlowTimeLeft = s.SlowTimeDuration
}

func applyWindowMode(mode WindowMode, resW, resH int) {
	if rl.IsWindowFullscreen() {
		rl.ToggleFullscreen()
	}
	if rl.IsWindowState(rl.FlagBorderlessWindowedMode) {
		rl.ToggleBorderlessWindowed()
	}

	monitor := rl.GetCurrentMonitor()
	mw := rl.GetMonitorWidth(monitor)
	mh := rl.GetMonitorHeight(monitor)

	switch mode {
	case WindowModeBorderless:
		if !rl.IsWindowState(rl.FlagWindowUndecorated) {
			rl.SetWindowState(rl.FlagWindowUndecorated)
		}
		rl.SetWindowSize(mw, mh)
		rl.SetWindowPosition(0, 0)
		rl.ToggleBorderlessWindowed()
		hideSystemUI()

	case WindowModeWindowed:
		if rl.IsWindowState(rl.FlagWindowUndecorated) {
			rl.ClearWindowState(rl.FlagWindowUndecorated)
		}
		rl.SetWindowSize(resW, resH)
		rl.SetWindowPosition((mw-resW)/2, (mh-resH)/2)
		restoreSystemUI()
	}

	currentConfig.WindowMode = mode
	if mode == WindowModeWindowed {
		currentConfig.ResWidth = resW
		currentConfig.ResHeight = resH
	}
	saveConfig()
}

func doOptionsMenu(w, h int) bool {
	if rl.IsCursorHidden() {
		rl.EnableCursor()
	}

	beginFrame()
	rl.ClearBackground(rl.Black)

	titleText := "Options"
	titleSz := int32(70)
	titleW := rl.MeasureText(titleText, titleSz)
	t := float32(rl.GetTime())
	hue := t * 60
	hue -= float32(int(hue/360)) * 360
	rl.DrawText(titleText, int32(w)/2-titleW/2, 100, titleSz, rl.ColorFromHSV(hue, 0.8, 1.0))

	if windowModeConfigurable {
		mp := virtualMousePosition()

		const sectionSz = int32(20)
		sectionColor := rl.NewColor(160, 160, 160, 255)

		wmLabel := "Window Mode"
		wmLabelW := rl.MeasureText(wmLabel, sectionSz)
		rl.DrawText(wmLabel, int32(w)/2-wmLabelW/2, 220, sectionSz, sectionColor)

		const modeCount = 2
		modeBtnW := float32(240)
		modeBtnH := float32(64)
		modeSpacing := float32(16)
		modeTotalW := float32(modeCount)*modeBtnW + float32(modeCount-1)*modeSpacing
		modeStartX := float32(w)/2 - modeTotalW/2
		modeStartY := float32(250)

		modeLabels := [modeCount]string{"Borderless", "Windowed"}
		for i := 0; i < modeCount; i++ {
			x := modeStartX + float32(i)*(modeBtnW+modeSpacing)
			rect := rl.Rectangle{X: x, Y: modeStartY, Width: modeBtnW, Height: modeBtnH}
			active := int(currentConfig.WindowMode) == i

			baseCol := rl.NewColor(40, 40, 50, 255)
			borderCol := rl.DarkGray
			textCol := rl.Gray
			if active {
				baseCol = rl.NewColor(30, 80, 30, 255)
				borderCol = rl.NewColor(80, 200, 80, 255)
				textCol = rl.White
			}
			if rl.CheckCollisionPointRec(mp, rect) {
				baseCol.R += 20
				baseCol.G += 20
				baseCol.B += 20
			}

			rl.DrawRectangleRec(rect, rl.Fade(baseCol, 0.8))
			rl.DrawRectangleLinesEx(rect, 2, borderCol)

			sz := int32(18)
			lw := rl.MeasureText(modeLabels[i], sz)
			rl.DrawText(modeLabels[i], int32(rect.X+rect.Width/2)-lw/2, int32(rect.Y+rect.Height/2)-sz/2, sz, textCol)

			if rl.CheckCollisionPointRec(mp, rect) && rl.IsMouseButtonReleased(rl.MouseLeftButton) && !active {
				applyWindowMode(WindowMode(i), currentConfig.ResWidth, currentConfig.ResHeight)
			}
		}

		resEnabled := currentConfig.WindowMode == WindowModeWindowed
		resAlpha := float32(0.35)
		if resEnabled {
			resAlpha = 1.0
		}
		resLabelCol := rl.NewColor(
			uint8(float32(160)*resAlpha),
			uint8(float32(160)*resAlpha),
			uint8(float32(160)*resAlpha),
			255,
		)

		resLabel := "Resolution  (windowed only)"
		resLabelW := rl.MeasureText(resLabel, sectionSz)
		rl.DrawText(resLabel, int32(w)/2-resLabelW/2, 358, sectionSz, resLabelCol)

		monitor := rl.GetCurrentMonitor()
		nativeW := rl.GetMonitorWidth(monitor)
		nativeH := rl.GetMonitorHeight(monitor)

		resolutions := make([]Resolution, 0, len(presetResolutions)+1)
		nativeInPresets := false
		for _, r := range presetResolutions {
			resolutions = append(resolutions, r)
			if r.Width == nativeW && r.Height == nativeH {
				nativeInPresets = true
			}
		}
		if !nativeInPresets {
			resolutions = append(resolutions, Resolution{nativeW, nativeH, "Native"})
		}

		resBtnW := float32(190)
		resBtnH := float32(54)
		resSpacing := float32(14)
		resTotalW := float32(len(resolutions))*resBtnW + float32(len(resolutions)-1)*resSpacing
		resStartX := float32(w)/2 - resTotalW/2
		resStartY := float32(390)

		for _, r := range resolutions {
			idx := r
			active := resEnabled && currentConfig.ResWidth == idx.Width && currentConfig.ResHeight == idx.Height

			x := resStartX
			resStartX += resBtnW + resSpacing

			rect := rl.Rectangle{X: x, Y: resStartY, Width: resBtnW, Height: resBtnH}

			baseCol := rl.NewColor(40, 40, 50, 255)
			borderCol := rl.DarkGray
			textCol := rl.Gray
			if active {
				baseCol = rl.NewColor(30, 80, 30, 255)
				borderCol = rl.NewColor(80, 200, 80, 255)
				textCol = rl.White
			}
			if resEnabled && rl.CheckCollisionPointRec(mp, rect) {
				baseCol.R += 20
				baseCol.G += 20
				baseCol.B += 20
			}

			fillAlpha := float32(0.8) * resAlpha
			rl.DrawRectangleRec(rect, rl.Fade(baseCol, fillAlpha))
			bCol := rl.NewColor(
				uint8(float32(borderCol.R)*resAlpha),
				uint8(float32(borderCol.G)*resAlpha),
				uint8(float32(borderCol.B)*resAlpha),
				255,
			)
			rl.DrawRectangleLinesEx(rect, 2, bCol)

			sz := int32(16)
			tCol := rl.NewColor(
				uint8(float32(textCol.R)*resAlpha),
				uint8(float32(textCol.G)*resAlpha),
				uint8(float32(textCol.B)*resAlpha),
				255,
			)
			lw := rl.MeasureText(idx.Label, sz)
			rl.DrawText(idx.Label, int32(rect.X+rect.Width/2)-lw/2, int32(rect.Y+rect.Height/2)-sz/2, sz, tCol)

			if resEnabled && rl.CheckCollisionPointRec(mp, rect) && rl.IsMouseButtonReleased(rl.MouseLeftButton) && !active {
				applyWindowMode(WindowModeWindowed, idx.Width, idx.Height)
			}
		}

	} else {
		const noteSz = int32(20)
		note := "Display: Borderless (only supported mode on this platform)"
		noteW := rl.MeasureText(note, noteSz)
		rl.DrawText(note, int32(w)/2-noteW/2, 240, noteSz, rl.NewColor(160, 160, 160, 255))
	}

	drawKeyBindings(w)

	backRect := rl.Rectangle{
		X:      float32(w)/2 - btnWidth/2,
		Y:      770,
		Width:  btnWidth,
		Height: btnHeight,
	}
	back := false
	if rebindTarget >= 0 {
		if rl.IsKeyPressed(rl.KeyEscape) {
			rebindTarget = -1
		}
		drawButton(backRect, "Back", rl.DarkGray, rl.Gray, rl.White)
	} else {
		back = drawButton(backRect, "Back", rl.DarkGray, rl.Gray, rl.White) || rl.IsKeyPressed(rl.KeyEscape)
	}

	endFrame()
	return back
}

var rebindTarget = -1

func drawKeyBindings(w int) {
	mp := virtualMousePosition()
	bindables := input.Bindables()

	const sectionSz = int32(20)
	sectionColor := rl.NewColor(160, 160, 160, 255)
	kbLabel := "Key Bindings"
	kbLabelW := rl.MeasureText(kbLabel, sectionSz)
	rl.DrawText(kbLabel, int32(w)/2-kbLabelW/2, 600, sectionSz, sectionColor)

	// Capture a key press for the active rebind target.
	if rebindTarget >= 0 {
		if key := rl.GetKeyPressed(); key != 0 && key != rl.KeyEscape {
			bindables[rebindTarget].Set(&keyMap, key)
			rebindTarget = -1
			saveKeyBindings()
		}
	}

	btnW := float32(150)
	btnH := float32(64)
	spacing := float32(16)
	totalW := float32(len(bindables))*btnW + float32(len(bindables)-1)*spacing
	startX := float32(w)/2 - totalW/2
	startY := float32(636)

	for i, b := range bindables {
		x := startX + float32(i)*(btnW+spacing)
		rect := rl.Rectangle{X: x, Y: startY, Width: btnW, Height: btnH}
		waiting := rebindTarget == i

		baseCol := rl.NewColor(40, 40, 50, 255)
		borderCol := rl.DarkGray
		textCol := rl.Gray
		if waiting {
			baseCol = rl.NewColor(80, 60, 20, 255)
			borderCol = rl.NewColor(220, 180, 60, 255)
			textCol = rl.White
		}
		if rl.CheckCollisionPointRec(mp, rect) {
			baseCol.R += 20
			baseCol.G += 20
			baseCol.B += 20
		}

		rl.DrawRectangleRec(rect, rl.Fade(baseCol, 0.8))
		rl.DrawRectangleLinesEx(rect, 2, borderCol)

		nameSz := int32(14)
		nameW := rl.MeasureText(b.Name, nameSz)
		rl.DrawText(b.Name, int32(rect.X+rect.Width/2)-nameW/2, int32(rect.Y)+10, nameSz, sectionColor)

		keyLabel := input.GetKeyName(b.Get(&keyMap))
		if waiting {
			keyLabel = "..."
		}
		keySz := int32(24)
		keyW := rl.MeasureText(keyLabel, keySz)
		rl.DrawText(keyLabel, int32(rect.X+rect.Width/2)-keyW/2, int32(rect.Y+rect.Height)-keySz-8, keySz, textCol)

		if rl.CheckCollisionPointRec(mp, rect) && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			if waiting {
				rebindTarget = -1
			} else {
				rebindTarget = i
			}
		}
	}
}
