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
	Shot1          rl.Sound
	Shot2          rl.Sound
	Shot3          rl.Sound
	Shot4          rl.Sound
	Shot5          rl.Sound
	Shot6          rl.Sound
	Shot7          rl.Sound
	Shot8          rl.Sound
	ShotNew        rl.Sound
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
}
