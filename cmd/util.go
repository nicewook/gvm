package cmd

import (
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
