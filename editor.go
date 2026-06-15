package main

import (
	"coldkiller2/stage"
	"coldkiller2/structure"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// editing is true while the in-game stage editor is open (test mode only).
var editing = false

type editorTool int

const (
	toolEnemy editorTool = iota
	toolStructure
	toolStart
)

var editorStartColor = rl.Gold

type editorState struct {
	stageIdx      int
	tool          editorTool
	kindIdx       int        // index into stage.EnemyKinds
	structKindIdx int        // index into stage.StructureKinds
	structSize    rl.Vector3 // size of the structure being placed (resizable)
	camTarget     rl.Vector3
	camFovy       float32 // orthographic view height; smaller = more zoomed in
	saveToast     float32 // seconds left to show the "SAVED" message
}

var editor = editorState{
	structSize: rl.Vector3{X: 2, Y: 2, Z: 2},
	camFovy:    60,
}

// selectStructure switches to the wall tool and resets the working size to the
// chosen kind's default.
func (e *editorState) selectStructure(idx int) {
	e.tool = toolStructure
	e.structKindIdx = idx
	e.structSize, _ = stage.StructureDef(stage.StructureKinds[idx])
}

var editorKindColors = map[stage.EnemyKind]rl.Color{
	stage.KindRobot:      rl.RayWhite,
	stage.KindSoldier:    rl.Green,
	stage.KindSniper:     rl.Red,
	stage.KindCharger:    rl.Orange,
	stage.KindSuperRobot: rl.Purple,
	stage.KindRival:      rl.Blue,
}

const editorPanelWidth = 360

// paletteButton is one clickable tool swatch in the left panel.
type paletteButton struct {
	rect    rl.Rectangle
	tool    editorTool
	kindIdx int // -1 for the wall tool
	label   string
	color   rl.Color
}

// editorPalette returns the clickable tool buttons (enemy kinds, structure
// kinds, player start). The geometry is shared by input handling and drawing so
// they stay in sync.
func editorPalette() []paletteButton {
	const (
		x   = float32(20)
		w   = float32(320)
		h   = float32(34)
		gap = float32(6)
	)
	startY := float32(96)
	row := 0

	buttons := make([]paletteButton, 0, len(stage.EnemyKinds)+len(stage.StructureKinds)+1)
	place := func(tool editorTool, idx int, label string, color rl.Color) {
		buttons = append(buttons, paletteButton{
			rect:    rl.Rectangle{X: x, Y: startY + float32(row)*(h+gap), Width: w, Height: h},
			tool:    tool,
			kindIdx: idx,
			label:   label,
			color:   color,
		})
		row++
	}

	for i, k := range stage.EnemyKinds {
		place(toolEnemy, i, string(k), editorKindColors[k])
	}
	for i, k := range stage.StructureKinds {
		_, color := stage.StructureDef(k)
		place(toolStructure, i, string(k), color)
	}
	place(toolStart, -1, "player start", editorStartColor)
	return buttons
}

func (b paletteButton) selected() bool {
	if b.tool != editor.tool {
		return false
	}
	switch b.tool {
	case toolEnemy:
		return b.kindIdx == editor.kindIdx
	case toolStructure:
		return b.kindIdx == editor.structKindIdx
	default:
		return true
	}
}

// openEditor enters the editor, clamping to a valid stage and showing the cursor.
func openEditor() {
	editing = true
	if editor.stageIdx >= len(stage.Stages) {
		editor.stageIdx = 0
	}
}

// editorCamera returns the top-down orthographic camera for the current view.
func (e *editorState) camera() rl.Camera3D {
	return rl.Camera3D{
		Position:   rl.Vector3{X: e.camTarget.X, Y: 60, Z: e.camTarget.Z},
		Target:     e.camTarget,
		Up:         rl.Vector3{X: 0, Y: 0, Z: -1},
		Fovy:       e.camFovy,
		Projection: rl.CameraOrthographic,
	}
}

// cursorWorld maps the (virtual) mouse position onto the y=0 ground plane.
func (e *editorState) cursorWorld(mouse rl.Vector2) rl.Vector3 {
	ray := rl.GetScreenToWorldRayEx(mouse, e.camera(), VirtualWidth, VirtualHeight)
	// Top-down orthographic: the ray originates at the cursor's world XZ and
	// points straight down, so the ground hit is just its XZ.
	return rl.Vector3{X: ray.Position.X, Y: 0, Z: ray.Position.Z}
}

func updateEditor(dt float32, mouse rl.Vector2) {
	if rl.IsCursorHidden() {
		rl.EnableCursor()
	}

	// Pan with WASD, zoom with the wheel.
	pan := editor.camFovy * 0.9 * dt
	if rl.IsKeyDown(rl.KeyW) {
		editor.camTarget.Z -= pan
	}
	if rl.IsKeyDown(rl.KeyS) {
		editor.camTarget.Z += pan
	}
	if rl.IsKeyDown(rl.KeyA) {
		editor.camTarget.X -= pan
	}
	if rl.IsKeyDown(rl.KeyD) {
		editor.camTarget.X += pan
	}
	if wheel := rl.GetMouseWheelMove(); wheel != 0 {
		editor.camFovy -= wheel * 4
		editor.camFovy = rl.Clamp(editor.camFovy, 10, 160)
	}

	// Resize the wall being placed (arrow keys, wall tool only).
	if editor.tool == toolStructure {
		if rl.IsKeyPressed(rl.KeyRight) {
			editor.structSize.X++
		}
		if rl.IsKeyPressed(rl.KeyLeft) && editor.structSize.X > 1 {
			editor.structSize.X--
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			editor.structSize.Z++
		}
		if rl.IsKeyPressed(rl.KeyDown) && editor.structSize.Z > 1 {
			editor.structSize.Z--
		}
	}

	// Stage navigation and creation.
	if rl.IsKeyPressed(rl.KeyQ) && editor.stageIdx > 0 {
		editor.stageIdx--
	}
	if rl.IsKeyPressed(rl.KeyE) && editor.stageIdx < len(stage.Stages)-1 {
		editor.stageIdx++
	}
	if rl.IsKeyPressed(rl.KeyN) {
		stage.Stages = append(stage.Stages, stage.StageSpec{})
		editor.stageIdx = len(stage.Stages) - 1
	}

	// Save to the source JSON.
	if rl.IsKeyPressed(rl.KeyF5) {
		if err := stage.Save(); err != nil {
			fmt.Println("stage save failed:", err)
		} else {
			editor.saveToast = 2.0
		}
	}
	if editor.saveToast > 0 {
		editor.saveToast -= dt
	}

	// Clicks over the left panel are UI (tool selection); clicks over the world
	// place or delete items.
	overPanel := mouse.X < editorPanelWidth
	cursor := editor.cursorWorld(mouse)
	spec := &stage.Stages[editor.stageIdx]

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		if overPanel {
			for _, b := range editorPalette() {
				if rl.CheckCollisionPointRec(mouse, b.rect) {
					switch b.tool {
					case toolEnemy:
						editor.tool = toolEnemy
						editor.kindIdx = b.kindIdx
					case toolStructure:
						editor.selectStructure(b.kindIdx)
					case toolStart:
						editor.tool = toolStart
					}
					break
				}
			}
		} else {
			switch editor.tool {
			case toolEnemy:
				spec.Enemies = append(spec.Enemies, stage.EnemySpec{
					Kind: stage.EnemyKinds[editor.kindIdx],
					X:    cursor.X,
					Z:    cursor.Z,
				})
			case toolStructure:
				spec.Structures = append(spec.Structures, stage.StructureSpec{
					Kind:     stage.StructureKinds[editor.structKindIdx],
					Position: rl.Vector3{X: cursor.X, Y: 0, Z: cursor.Z},
					Size:     editor.structSize,
				})
			case toolStart:
				spec.Start = stage.Point{X: cursor.X, Z: cursor.Z}
			}
		}
	}
	if rl.IsMouseButtonPressed(rl.MouseRightButton) && !overPanel && editor.tool != toolStart {
		editor.deleteNearest(cursor)
	}
}

func (e *editorState) deleteNearest(p rl.Vector3) {
	spec := &stage.Stages[e.stageIdx]
	const pickRadius = float32(4)
	if e.tool == toolEnemy {
		best, bestD := -1, pickRadius
		for i, en := range spec.Enemies {
			d := rl.Vector2Distance(rl.Vector2{X: en.X, Y: en.Z}, rl.Vector2{X: p.X, Y: p.Z})
			if d < bestD {
				best, bestD = i, d
			}
		}
		if best >= 0 {
			spec.Enemies = append(spec.Enemies[:best], spec.Enemies[best+1:]...)
		}
		return
	}
	best, bestD := -1, pickRadius
	for i, s := range spec.Structures {
		d := rl.Vector2Distance(rl.Vector2{X: s.Position.X, Y: s.Position.Z}, rl.Vector2{X: p.X, Y: p.Z})
		if d < bestD {
			best, bestD = i, d
		}
	}
	if best >= 0 {
		spec.Structures = append(spec.Structures[:best], spec.Structures[best+1:]...)
	}
}

func drawEditor(mouse rl.Vector2) {
	cam := editor.camera()
	cursor := editor.cursorWorld(mouse)
	spec := stage.Stages[editor.stageIdx]

	beginFrame()
	rl.ClearBackground(rl.NewColor(20, 20, 26, 255))

	rl.BeginMode3D(cam)
	rl.DrawGrid(100, 2)

	for _, s := range spec.Structures {
		_, color := stage.StructureDef(s.Kind)
		st := structure.Structure{Position: s.Position, Size: s.Size, Direction: s.Direction, Color: color}
		st.Draw3D()
	}

	for _, en := range spec.Enemies {
		drawEnemyMarker(rl.Vector3{X: en.X, Y: 0, Z: en.Z}, editorKindColors[en.Kind])
	}

	drawStartMarker(rl.Vector3{X: spec.Start.X, Y: 0, Z: spec.Start.Z}, editorStartColor)

	// Cursor preview for the active tool.
	switch editor.tool {
	case toolEnemy:
		preview := editorKindColors[stage.EnemyKinds[editor.kindIdx]]
		preview.A = 150
		drawEnemyMarker(cursor, preview)
	case toolStructure:
		_, color := stage.StructureDef(stage.StructureKinds[editor.structKindIdx])
		center := rl.Vector3{X: cursor.X, Y: editor.structSize.Y / 2, Z: cursor.Z}
		rl.DrawCube(center, editor.structSize.X, editor.structSize.Y, editor.structSize.Z, rl.Fade(color, 0.5))
		rl.DrawCubeWires(center, editor.structSize.X, editor.structSize.Y, editor.structSize.Z, rl.SkyBlue)
	case toolStart:
		preview := editorStartColor
		preview.A = 150
		drawStartMarker(cursor, preview)
	}
	rl.EndMode3D()

	drawEditorHUD(spec, mouse)
	endFrame()
}

func drawEnemyMarker(pos rl.Vector3, col rl.Color) {
	rl.DrawCylinder(pos, 0.6, 0.6, 1.4, 12, col)
	rl.DrawCylinderWires(pos, 0.6, 0.6, 1.4, 12, rl.Black)
}

func drawStartMarker(pos rl.Vector3, col rl.Color) {
	rl.DrawCylinder(pos, 0.3, 0.3, 2.4, 10, col)
	rl.DrawSphere(rl.Vector3{X: pos.X, Y: 2.8, Z: pos.Z}, 0.6, col)
}

func drawEditorHUD(spec stage.StageSpec, mouse rl.Vector2) {
	const pad = int32(20)
	rl.DrawRectangle(0, 0, editorPanelWidth, VirtualHeight, rl.NewColor(0, 0, 0, 150))

	rl.DrawText(fmt.Sprintf("STAGE EDITOR   %d / %d", editor.stageIdx+1, len(stage.Stages)), pad, 20, 24, rl.Gold)
	rl.DrawText(fmt.Sprintf("enemies: %d    walls: %d", len(spec.Enemies), len(spec.Structures)), pad, 54, 18, rl.LightGray)

	// Clickable tool palette: a color swatch + label per enemy/structure kind,
	// plus the player start.
	palette := editorPalette()
	for _, b := range palette {
		bg := rl.NewColor(40, 40, 48, 255)
		switch {
		case b.selected():
			bg = rl.NewColor(70, 90, 120, 255)
		case rl.CheckCollisionPointRec(mouse, b.rect):
			bg = rl.NewColor(55, 55, 66, 255)
		}
		rl.DrawRectangleRec(b.rect, bg)
		if b.selected() {
			rl.DrawRectangleLinesEx(b.rect, 2, rl.Gold)
		}
		swatch := rl.Rectangle{X: b.rect.X + 8, Y: b.rect.Y + 8, Width: b.rect.Height - 16, Height: b.rect.Height - 16}
		rl.DrawRectangleRec(swatch, b.color)
		rl.DrawRectangleLinesEx(swatch, 1, rl.Black)
		rl.DrawText(b.label, int32(b.rect.X)+int32(b.rect.Height)+6, int32(b.rect.Y)+11, 18, rl.White)
	}

	last := palette[len(palette)-1].rect
	y := int32(last.Y+last.Height) + 16
	line := func(s string, c rl.Color) {
		rl.DrawText(s, pad, y, 16, c)
		y += 24
	}

	if editor.tool == toolStructure {
		kind := stage.StructureKinds[editor.structKindIdx]
		line(fmt.Sprintf("%s  %.0f x %.0f  (arrows)", kind, editor.structSize.X, editor.structSize.Z), rl.SkyBlue)
		y += 4
	}
	for _, h := range []string{
		"L-Click  place      R-Click  delete",
		"WASD  pan           Wheel  zoom",
		"Q / E  prev / next stage",
		"N  new stage        F5  save",
		"ESC  back to menu",
	} {
		line(h, rl.Gray)
	}

	if editor.saveToast > 0 {
		rl.DrawText("SAVED", VirtualWidth/2-50, 30, 34, rl.Green)
	}
}
