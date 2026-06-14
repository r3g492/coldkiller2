package main

import (
	"coldkiller2/killer"
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// comboWindow is the real-time grace period (in seconds) after a kill before
// the streak breaks. It is measured in real time so slow-time doesn't drain it.
const comboWindow = 2.5

// comboSlowRefill is how much slow-time (seconds) each kill restores, indexed
// by combo tier. Longer streaks pour the meter back faster, so staying
// aggressive sustains slow-mo.
var comboSlowRefill = []float32{0.10, 0.18, 0.30, 0.45}

var (
	comboCount int     // enemies killed in the current uninterrupted streak
	comboTimer float32 // counts down to 0; when it hits 0 the streak resets
	comboPunch float32 // 0..1 scale/shake punch, decays each frame
)

func resetCombo() {
	comboCount = 0
	comboTimer = 0
	comboPunch = 0
}

// updateCombo advances streak state. kills is the number of enemies the player
// killed this frame; dt is real (unscaled) time so the window and punch animate
// at normal speed regardless of slow-time.
func updateCombo(kills int, dt float32, player *killer.Killer) {
	if kills > 0 {
		comboCount += kills
		comboTimer = comboWindow
		comboPunch = 1.0
		if comboCount > currentConfig.BestCombo {
			currentConfig.BestCombo = comboCount
		}
		// A real streak (2+) starts feeding the slow-time meter, scaled by tier.
		if comboCount >= 2 {
			refill := comboSlowRefill[comboTier()] * float32(kills)
			if player.SlowTimeLeft < player.SlowTimeDuration {
				player.SlowTimeLeft += refill
				if player.SlowTimeLeft > player.SlowTimeDuration {
					player.SlowTimeLeft = player.SlowTimeDuration
				}
				player.SlowRefillFlash = 1.0
			}
		}
	} else if comboTimer > 0 {
		comboTimer -= dt
		if comboTimer <= 0 {
			comboCount = 0
		}
	}

	if comboPunch > 0 {
		comboPunch -= dt * 4
		if comboPunch < 0 {
			comboPunch = 0
		}
	}
}

var comboTierColors = []rl.Color{
	rl.NewColor(255, 230, 120, 255), // tier 0: warm yellow
	rl.NewColor(255, 160, 60, 255),  // tier 1: orange
	rl.NewColor(255, 90, 60, 255),   // tier 2: red
	rl.NewColor(255, 60, 200, 255),  // tier 3: magenta
}

func comboTier() int {
	switch {
	case comboCount >= 20:
		return 3
	case comboCount >= 10:
		return 2
	case comboCount >= 5:
		return 1
	default:
		return 0
	}
}

func drawCombo(w int) {
	if comboCount < 2 {
		return
	}

	tier := comboTier()
	color := comboTierColors[tier]

	baseSize := float32(46 + tier*10)
	size := baseSize * (1 + comboPunch*0.25)
	fontSize := int32(size)

	text := fmt.Sprintf("%d", comboCount)
	textW := rl.MeasureText(text, fontSize)

	cx := int32(w) / 2
	y := int32(95)

	shake := comboPunch * float32(2+tier*2)
	var dx, dy int32
	if shake > 0 {
		dx = int32((rand.Float32()*2 - 1) * shake)
		dy = int32((rand.Float32()*2 - 1) * shake)
	}

	// glow halo behind the number
	for i := 3; i >= 1; i-- {
		gc := color
		gc.A = uint8(60 / i)
		rl.DrawText(text, cx-textW/2+int32(i)+dx, y+int32(i)+dy, fontSize, gc)
	}
	rl.DrawText(text, cx-textW/2+dx, y+dy, fontSize, color)

	// "COMBO" label above the number
	label := "COMBO"
	labelSize := int32(16)
	lw := rl.MeasureText(label, labelSize)
	rl.DrawText(label, cx-lw/2+dx, y-labelSize-2+dy, labelSize, color)

	// streak timer bar draining beneath the number
	if comboTimer > 0 {
		barW := float32(textW)
		if barW < 60 {
			barW = 60
		}
		pct := comboTimer / comboWindow
		barX := float32(cx) - barW/2
		barY := float32(y) + size + 6
		rl.DrawRectangleRec(rl.NewRectangle(barX, barY, barW, 5), rl.Fade(rl.DarkGray, 0.5))
		rl.DrawRectangleRec(rl.NewRectangle(barX, barY, barW*pct, 5), color)
	}
}
