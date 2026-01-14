package sound

import (
	"coldkiller2/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ShotgunSound   rl.Sound
	ReloadingSound rl.Sound
)

func Init() {
	ShotgunSound = util.LoadSoundFromEmbedded("shotgun-03-38220.mp3")
	ReloadingSound = util.LoadSoundFromEmbedded("1911-reload-6248.mp3")
}
