package util

import (
	"embed"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed sounds/*
var resFS embed.FS

func LoadSoundFromEmbedded(filename string) rl.Sound {
	data, err := resFS.ReadFile("sounds/" + filename)
	if err != nil {
		log.Fatalf("failed to read embedded sound %s: %v", filename, err)
	}
	tmpFile, err := os.CreateTemp("", "*.mp3")
	if err != nil {
		log.Fatalf("failed to create temporary file for %s: %v", filename, err)
	}
	_, err = tmpFile.Write(data)
	if err != nil {
		log.Fatalf("failed to write to temporary file for %s: %v", filename, err)
	}
	_ = tmpFile.Close()
	snd := rl.LoadSound(tmpFile.Name())
	return snd
}
