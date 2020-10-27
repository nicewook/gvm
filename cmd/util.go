package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func fileExist(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal("os.Stat: ", err)
	}
	return !info.IsDir()
}

func copyFile(src, des string) {
	fmt.Printf("copy from %s to %s\n", src, des)
	from, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(des, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

func renameToGo(s string) string {
	// get directory
	dir := filepath.Dir(s)
	return filepath.Join(dir, "go.exe")
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

const (
	green = "green"

	checkBox = "\xE2\x9C\x85"

	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func colorPrint(color string, msg string) error {

	switch color {
	case green:
		fmt.Printf("%s%s%s\n", string(colorGreen), msg, checkBox)
	default:
		return errors.New("not proper color")
	}
	fmt.Print(string(colorReset))
	return nil
}
