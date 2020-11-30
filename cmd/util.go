package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	ver "github.com/hashicorp/go-version"
)

func fileExist(filePath string) bool {
	fmtV.Printf("check file exist: %s\n", filePath)
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
	Green = "green"
	Red   = "red"

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

func lefAlignString(msg string) string {
	return fmt.Sprintf("%-20v", msg)
}

func addStar(msg string) string {
	return msg + " *"
}

func colorPrint(color string, msg string) error {
	leftAlignMsg := lefAlignString(msg)

	switch color {
	case Red:
		fmt.Printf("%s%s", string(colorRed), leftAlignMsg)
	case Green:
		fmt.Printf("%s%s", string(colorGreen), leftAlignMsg)
	default:
		return errors.New("not proper color")
	}
	fmt.Print(string(colorReset))
	return nil
}

func colorPrintLeftAlign(color string, msg string) error {

	return colorPrint(color, lefAlignString(msg))
}

func goroot(version string) (string, error) {
	home, err := homedir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	return filepath.Join(home, "sdk", version), nil
}

func homedir() (string, error) {
	if dir := os.Getenv("USERPROFILE"); dir != "" {
		return dir, nil
	}
	return "", errors.New("can't find user home directory; %USERPROFILE% is empty")
}

func sortGoSDKList(list []string) {

	sort.Slice(list, func(i, j int) bool {
		va := strings.TrimPrefix(list[i], "go")
		vA, err1 := ver.NewVersion(va)
		if err1 != nil {
			fmtV.Println(err1)
			return false
		}

		vb := strings.TrimPrefix(list[j], "go")
		vB, err2 := ver.NewVersion(vb)
		if err2 != nil {
			fmtV.Println(err2)
			return false
		}

		return vA.LessThan(vB)
	})
}

const (
	getGoExePathCmdStr = "where go.exe"
	getCurGoVerCmdStr  = "go.exe version"
)

func getCurGoVersion() string {
	getCurGoVerCmd := makeExecCmd(getCurGoVerCmdStr)
	b, err := getCurGoVerCmd.Output()
	if err != nil {
		log.Fatal("getCurGoVerCmd: ", err)
	}
	versionOutput := strings.Split(string(b), " ")
	return versionOutput[2]
}

func getCurGoExePath() (string, bool) {
	getCurGoExePathCmd := makeExecCmd(getGoExePathCmdStr)
	b, err := getCurGoExePathCmd.Output()
	if err != nil {
		log.Fatal("getCurGoExePathCmd: ", err)
	}
	curGoExePath := strings.TrimSpace(string(b)) // needed to remove space

	var isSystemGo bool
	fmtV.Printf("current go.exe path is %s, GOROOT is %s\n", curGoExePath, goRoot)
	if strings.Contains(
		strings.ToLower(curGoExePath),
		strings.ToLower(goRoot),
	) {
		isSystemGo = true
	}
	return curGoExePath, isSystemGo
}

func makeExecCmd(cmdStr string) *exec.Cmd {
	cmdSlice := strings.Split(cmdStr, " ")
	if len(cmdSlice) <= 0 {
		return exec.Command(cmdSlice[0])
	}
	return exec.Command(cmdSlice[0], cmdSlice[1:]...)
}
