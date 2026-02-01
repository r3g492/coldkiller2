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
)

func Init() {
	ShotgunSound = util.LoadSoundFromEmbedded("shotgun-03-38220.mp3")
	rl.SetSoundVolume(ShotgunSound, 0.2)
	ReloadingSound = util.LoadSoundFromEmbedded("1911-reload-6248.mp3")
	rl.SetSoundVolume(ReloadingSound, 0.4)
	Track = util.LoadSoundFromEmbedded("song2.mp3")
	YouLose = util.LoadSoundFromEmbedded("you-lose-game-sound-230514.mp3")
}
