# Kill Per Second

top-down action shooter written in Go with raylib-go. released on Steam for windows, macos and linux.

<img width="720" height="405" alt="output" src="https://github.com/user-attachments/assets/6102b293-1c82-4626-b714-2db04ed721b4" />

steam page: https://store.steampowered.com/app/4690070/Kill_Per_Second/

- ~6k lines of Go, 2 direct dependencies (`raylib-go`, `go-steamworks`)
- 100 stages, 8 enemy types, 4 AI behaviours
- single binary per platform with all assets embedded

## design notes

this is a game, not a backend service, but these are the structural choices I made while building it.

### packages and the frame loop

each gameplay area (`enemy`, `bullet`, `blast`, `structure`, `stage`, `killer`) is a package with a `Manager` that owns its state and has `Init / Mutate / Draw3D / Unload`. `main.go` calls them in a fixed order every frame: input → stage → enemies → player → bullets → blasts → visibility → render.

the player and enemies don't create bullets directly. they return `[]BulletCmd` values and `bullet.Manager` creates the bullets. this keeps `enemy` from importing `bullet`, so the package dependencies only go one way.

### world time vs real time

slow-motion runs the world at 0.1× speed. enemies, bullets and blasts use the scaled `worldDt`. the combo timer, the slow-mo meter and UI animations use real `dt`, so slow-mo doesn't stretch the combo window. both values are passed in explicitly (`player.Mutate(ip, worldDt, dt, ...)`).

### per-frame performance

- per-frame buffers are reset with `buf[:0]` and reused instead of reallocated
- dead bullets and blasts are removed by swapping with the last element
- structures and enemies are stored in spatial grids, so collision checks only look at nearby cells
- bullets move more than an enemy's width per frame, so hits are checked against the segment from the previous position to the current one. movement is resolved one axis at a time so characters slide along walls

### stages as data

the 100 stages are stored in `stage/data/stages.json` and embedded with `go:embed`. enemy and structure kinds are strings mapped to factory functions, and unknown kinds are skipped. there's an in-game stage editor (test mode only) that edits the same data and writes it back to the JSON file.

### save file

one JSON file in the user config directory. missing fields fall back to zero values, so older saves still load. the save stores the last cleared stage, and the menu works out whether to show *Retry* or *GAME CLEARED* by comparing it to the number of stages, without a separate flag. progress is saved when a stage is cleared.

### steam and platform code

Steam support is behind a build tag: `steam.go` has the real implementation and `steam_stub.go` has empty functions, so the rest of the code doesn't check for Steam. platform-specific window handling uses the same approach (`platform_darwin.go`, `platform_linux.go`, `console_windows.go`).

achievement unlocks are queued and retried once a second, because Steam rejects them until user stats have loaded. players who cleared the game before achievements were added get the achievement the next time they launch the game. if Steam fails to start, the game runs without achievements.

### smaller things

- the game renders at 1600×900 and is letterboxed to the window. mouse input is converted back to those coordinates
- stage loading is split across frames so the loading screen keeps drawing
- if a stage drags on, extra robots start spawning after 30s, capped at 12 per wave
- upgrades and rebindable keys are defined as tables, so adding one is a single entry

### known shortcuts

- screen state (`showInitMenu`, `intermission`, `paused`, ...) is a set of package-level booleans rather than a proper state machine
- no automated tests. I tested by playing, with a test-mode panel for jumping to any stage
- embedded models and sounds are written to temp files for raylib to load and never cleaned up

## build

linux (on linux):
```
go build -o coldkiller2 .
```

windows (on linux, requires mingw-w64):
```
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags "-H=windowsgui" -o coldkiller2.exe .
```
`-H=windowsgui` builds a GUI exe so no console window opens alongside the game.

macos (on macos, requires xcode command line tools):
```
go build -o coldkiller2 .
```

## steam build

steam builds add `-tags steam` (enables achievements). every build uploaded to steam needs it.

linux (on linux):
```
go build -tags steam -o coldkiller2 .
```

windows (on linux, requires mingw-w64):
```
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -tags steam -ldflags "-H=windowsgui" -o coldkiller2.exe .
```

macos universal binary (on macos):
```
GOARCH=arm64 CGO_ENABLED=1 go build -tags steam -o coldkiller2_arm64 .
GOARCH=amd64 CGO_ENABLED=1 go build -tags steam -o coldkiller2_amd64 .
lipo -create -output coldkiller2 coldkiller2_arm64 coldkiller2_amd64
```

depot contents:
- windows: `coldkiller2.exe` + `steam_api64.dll` (from Steamworks SDK **1.65** `redistributable_bin/win64/`, not in this repo; must match the SDK version go-steamworks targets, or the game crashes on start)
- linux: `coldkiller2` (steam library is embedded)
- macos: `coldkiller2` (steam library is embedded)


## license

source code is licensed under the [MIT License](LICENSE).

assets in `util/images/`, `util/models/` and `util/sounds/` are **not** covered by the MIT License. they are © their respective owners and included only so the game builds. do not reuse or redistribute them separately.
- artwork: CaveTroll, used with permission
- sounds from [Pixabay](https://pixabay.com), used under the [Pixabay Content License](https://pixabay.com/service/license-summary/): `1911-reload-6248.mp3`, `data_pion-st1-footstep-sfx-323053.mp3`, `shot_new_trimmed.mp3`, `shotgun-03-38220.mp3`
- original music by r3g492: `song2.mp3`, `win.mp3`, `lost.mp3`
