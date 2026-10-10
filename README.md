hardcore action shooter in raylib-go  
supports windows, macos, linux  

<img width="720" height="405" alt="output" src="https://github.com/user-attachments/assets/6102b293-1c82-4626-b714-2db04ed721b4" />


steam page: https://store.steampowered.com/app/4690070/Kill_Per_Second/

## build (on linux)

linux:
```
go build -o coldkiller2 .
```

windows (requires mingw-w64):
```
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags "-H=windowsgui" -o coldkiller2.exe .
```
`-H=windowsgui` builds a GUI exe so no console window opens alongside the game.
