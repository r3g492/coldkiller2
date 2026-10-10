hardcore action shooter in raylib-go  
supports windows, macos, linux  

<img width="720" height="405" alt="output" src="https://github.com/user-attachments/assets/6102b293-1c82-4626-b714-2db04ed721b4" />


steam page: https://store.steampowered.com/app/4690070/Kill_Per_Second/

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

do not ship `steam_appid.txt`; it is only for local testing.
