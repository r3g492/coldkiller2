package stage

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func savePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), "save.dat"), nil
}

func SaveProgress(highestBeaten int) {
	path, err := savePath()
	if err != nil {
		return
	}
	os.WriteFile(path, []byte(strconv.Itoa(highestBeaten)), 0644)
}

func LoadProgress() int {
	path, err := savePath()
	if err != nil {
		return 0
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return n
}
