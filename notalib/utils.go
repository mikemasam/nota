package notalib

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func Color(colorCode string) string {
	return fmt.Sprintf("\x1b[38;5;%sm", colorCode)
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true 
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func ResolveHomeDir(loc string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("home path not found")
	}
	return fmt.Sprintf("%s/%s", home, loc)
}
