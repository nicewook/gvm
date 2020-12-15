package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	ver "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

func fileExist(filePath string) bool {
	var result bool

	fmtV.Printf("check %s exists: ", filePath)
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			result = false
			fmtV.Println(result)
			return result
		}
		log.Fatal("os.Stat: ", err)
	}
	result = !info.IsDir()
	fmtV.Println(result)
	return result
}

func copyFile(src, des string) {
	fmtV.Printf("copy from %s to %s\n", src, des)
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

func removeFile(curGoExePath string) {
	fmtV.Printf("remove %s\n", curGoExePath)
	if err := os.Remove(curGoExePath); err != nil {
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

	switch color {
	case Red:
		fmt.Printf("%s%s", string(colorRed), msg)
	case Green:
		fmt.Printf("%s%s", string(colorGreen), msg)
	default:
		return errors.New("not proper color")
	}
	fmt.Print(string(colorReset))
	return nil
}

func makeColorString(color string, msg string) string {
	return fmt.Sprintf("%s%s%s", color, msg, colorReset)
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

func getSystemGoExeName() string {
	dirPath := filepath.Join(goRoot, "bin")
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal("ReadDir: ", err, dirPath)
	}
	for _, file := range files {
		name := file.Name()
		if name == "go.exe" {
			return "go.exe"
		}
		name = strings.TrimRight(name, ".exe")
		if isGoVersionString(name) {
			return name + ".exe"
		}
	}
	log.Fatal("no system go execution file found")
	return ""
}

func getSystemGoVer() string {
	getSystemGoVerCmd := makeExecCmd(fmt.Sprintf("%s version", getSystemGoExeName()))
	b, err := getSystemGoVerCmd.Output()
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
	fmtV.Printf("currently using %s\n", curGoExePath)
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

func noArgumentDisplayHelp(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(0)
	}
}

// regex for go version
const semVerRegex string = `go?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

func isGoVersionString(version string) bool {
	re := regexp.MustCompile("^" + semVerRegex + "$")
	return re.MatchString(version)
}

func isInstalledVersion(version string) bool {
	for _, installedVer := range getLocalList() {
		if installedVer == version {
			return true
		}
	}
	return false
}
