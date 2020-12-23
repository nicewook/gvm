/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func useSystemGo() {

	systemGoPath := filepath.Join(goRoot, "bin")
	files, err := ioutil.ReadDir(systemGoPath)
	if err != nil {
		log.Fatal(err)
	}
	var fNames []string
	for _, f := range files {
		fNames = append(fNames, f.Name())
	}

	for _, fn := range fNames {
		ver := strings.TrimRight(fn, ".exe")
		if isGoVersionString(ver) {
			fmtV.Printf("trying to use system Go SDK: %s\n", ver)

			// check if currently using system go
			curGoExePath, isSystemGo := getCurGoExePath()
			if isSystemGo {
				fmt.Printf("already using system Go. %s\n", makeColorString(colorGreen, curGoExePath))
				return
			}

			// if using go.exe is not system go,
			// copy to go<version>.exe and remove current using go.exe
			curGoExeVersion := getCurGoExeVersion()
			fmt.Printf("We used %s, %s\n", curGoExePath, makeColorString(colorGreen, strings.TrimSuffix(curGoExeVersion, ".exe")))
			copyFile(curGoExePath, renameToGoVersion(curGoExePath, curGoExeVersion))
			removeFile(curGoExePath)

			// copy
			curGoExe := filepath.Join(systemGoPath, ver+".exe")
			goExe := filepath.Join(systemGoPath, "go.exe")
			copyFile(curGoExe, goExe)

			fmt.Printf("now we are using %s\n", makeColorString(colorGreen, getCurGoExeVersion()))
			return
		}
	}
	log.Fatal("no Go SDK in GOROOT")
}

// func getCurGoExePath() string {
// 	getGoExePathCmd := exec.Command("where", "go.exe")

// 	b, err := getGoExePathCmd.Output()
// 	if err != nil {
// 		log.Fatalf("getGoExePathCmd: %v", err)
// 	}
// 	curGoExePath := strings.TrimSpace(string(b)) // needed to remove space
// 	fmtV.Printf("current version: %s\n", curGoExePath)
// 	return curGoExePath
// }

func getCurGoExeVersion() string {
	getCurVersionCmd := exec.Command("go", "version")

	b, err := getCurVersionCmd.Output()
	if err != nil {
		log.Fatalf("getCurVersionCmd: %v", err)
	}

	versionOutput := strings.Split(string(b), " ")
	curVersionExe := versionOutput[2] + ".exe"
	fmtV.Printf("current go.exe version: %v\n", curVersionExe)
	return curVersionExe
}

func useVersion(version string) { // ex) version == 1.15.2 (without "go")

	if version == systemGo {
		fmt.Println("use system version")
		useSystemGo()
		return
	}

	useVersion := "go" + version
	useExe := useVersion + ".exe"

	// check regex of the version name
	if isGoVersionString(useVersion) == false {
		fmt.Printf("%s is not proper go version format\n", makeColorString(colorRed, version))
		os.Exit(0)
	} else {
		fmtV.Printf("%s is good go version format\n", makeColorString(colorGreen, version))
	}
	// check the version exist or already downloaded
	if alreadyInstalled(useVersion) == false {
		fmt.Printf("%s is not installed version", makeColorString(colorRed, version))
		os.Exit(0)
	}

	fmtV.Printf("wanted version: %s, exe: %s\n", useVersion, useExe)

	// check current go.exe exist - we will remove it
	curGoExePath, _ := getCurGoExePath()
	if fileExist(curGoExePath) == false {
		fmt.Println("cannot find currently using go.exe file")
		os.Exit(0)
		return
	}

	// check check the version to use
	// then copy the go<required-version>.exe to go.exe
	useExeFullPath := filepath.Join(goPath, "bin", useExe)
	if fileExist(useExeFullPath) == false {
		fmt.Printf("cannot find go version that we want to use: %v\n", useExeFullPath)
		os.Exit(0)
		return
	}

	// so we are ready to change go version
	curGoExeVersion := getCurGoExeVersion()
	fmt.Printf("We used %s, %s\n", curGoExePath, makeColorString(colorGreen, strings.TrimSuffix(curGoExeVersion, ".exe")))
	copyFile(curGoExePath, renameToGoVersion(curGoExePath, curGoExeVersion))
	removeFile(curGoExePath)

	copyFile(useExeFullPath, renameToGo(useExeFullPath))

	fmt.Printf("now we are using %s\n", makeColorString(colorGreen, getCurGoExeVersion()))
}

func getAllGoExePath() []string {
	cmd := exec.Command("where", "go.exe")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	defer cmd.Wait()

	var goExeFiles []string
	buff := bufio.NewScanner(stdout)
	for buff.Scan() {
		goExeFiles = append(goExeFiles, buff.Text())
	}

	fmtV.Println("current go.exe files.", goExeFiles)
	return goExeFiles
}

func onlyOneGoExeAllowed() {

	goExeFiles := getAllGoExePath()
	if len(goExeFiles) <= 1 {
		fmtV.Printf("go.exe count: %s, no other go.exe\n", len(goExeFiles))
		return
	}
	fmtV.Printf("go.exe count: %s\n", len(goExeFiles))

	// check systemGo exist
	var systemGoExist bool
	for _, v := range goExeFiles {
		if strings.Contains(v, goRoot) {
			systemGoExist = true
		}
	}

	// if systemGoExist, remove all go.exe in goPath
	// or leave the first go.exe and remove others all
	for i, v := range goExeFiles {
		if strings.Contains(v, goRoot) == false {
			if i == 0 && systemGoExist == false { // one go.exe should exist
				continue
			}
			if err := os.Remove(v); err != nil {
				log.Fatal(err)
			}
		}
	}

	goExeFiles = getAllGoExePath() // see the result
}

func use(cmd *cobra.Command, args []string) {

	onlyOneGoExeAllowed()
	if len(args) <= 0 {
		curVer := getCurGoVersion()
		if _, isSystem := getCurGoExePath(); isSystem {
			fmt.Printf("Currently using %s, %s\n", makeColorString(colorGreen, systemGo), makeColorString(colorGreen, curVer))
			return
		}
		fmt.Printf("Currently using %s\n", makeColorString(colorGreen, curVer))
		return
	}
	useVersion(args[0])
}

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Changes the Go SDK version to your desired version.",
	Long: `Changes the Go SDK version to your desired version.

ex) $ gvm use 1.13.2`,
	Args: cobra.MaximumNArgs(1),
	Run:  use,
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
