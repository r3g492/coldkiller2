package util

import (
	"embed"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed sounds/*
//go:embed images/*
//go:embed models/*
var resFS embed.FS

func LoadTextureFromEmbedded(filename string) rl.Texture2D {
	data, err := resFS.ReadFile("images/" + filename)
	if err != nil {
		log.Fatalf("failed to read embedded image %s: %v", filename, err)
	}
	ext := ".png"
	img := rl.LoadImageFromMemory(ext, data, int32(len(data)))
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return tex
}

func LoadModelFromEmbedded(filename string) (rl.Model, []rl.ModelAnimation) {
	data, err := resFS.ReadFile("models/" + filename)
	if err != nil {
		log.Fatalf("failed to read embedded model %s: %v", filename, err)
	}
	tmpFile, err := os.CreateTemp("", "*.glb")
	if err != nil {
		log.Fatalf("failed to create temporary file for %s: %v", filename, err)
	}
	_, err = tmpFile.Write(data)
	if err != nil {
		log.Fatalf("failed to write to temporary file for %s: %v", filename, err)
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()
	return rl.LoadModel(tmpPath), rl.LoadModelAnimations(tmpPath)
}

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
