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
	ReloadingSound = util.LoadSoundFromEmbedded("1911-reload-6248.mp3")
	Track = util.LoadSoundFromEmbedded("spinning-head-271171.mp3")
	YouLose = util.LoadSoundFromEmbedded("you-lose-game-sound-230514.mp3")
}
