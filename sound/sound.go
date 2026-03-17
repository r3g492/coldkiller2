package sound

import (
	"coldkiller2/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ShotgunSound   rl.Sound
	ReloadingSound rl.Sound
	Track          rl.Sound
	YouLose        rl.Sound
	ShotNew        rl.Sound
	FootStep       rl.Sound
)

func Init() {
	ShotgunSound = util.LoadSoundFromEmbedded("shotgun-03-38220.mp3")
	rl.SetSoundVolume(ShotgunSound, 0.3)
	ReloadingSound = util.LoadSoundFromEmbedded("1911-reload-6248.mp3")
	rl.SetSoundVolume(ReloadingSound, 0.4)
	Track = util.LoadSoundFromEmbedded("song2.mp3")
	rl.SetSoundVolume(Track, 1.2)
	YouLose = util.LoadSoundFromEmbedded("lost.mp3")
	ShotNew = util.LoadSoundFromEmbedded("shot_new.mp3")
	rl.SetSoundVolume(ShotNew, 0.5)
	FootStep = util.LoadSoundFromEmbedded("data_pion-st1-footstep-sfx-323053.mp3")
}

func PlaySound3D(s rl.Sound, sourcePos rl.Vector3, playerPos rl.Vector3, volumeMult float32) {
	maxHearingDistance := float32(100.0)

	dist := rl.Vector3Distance(playerPos, sourcePos)
	volume := 1.0 - (dist / maxHearingDistance)

	if volume < 0.0 {
		volume = 0.0
	}

	dx := sourcePos.X - playerPos.X

	pan := 0.5 - (dx / maxHearingDistance)

	if pan < 0.0 {
		pan = 0.0
	} else if pan > 1.0 {
		pan = 1.0
	}

	rl.SetSoundVolume(s, volume*volumeMult)
	rl.SetSoundPan(s, pan)
	rl.PlaySound(s)
}
