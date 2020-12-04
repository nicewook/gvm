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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MaximumNArgs(1),
	Run:  use,
}

func versionExist(file string) bool {
	filePath := filepath.Join(goPath, "bin", file)
	fmt.Println("want to use version of: ", filePath)
	return fileExist(filePath)
}

func useSystemGo() {

	systemGoPath := filepath.Join(goRoot, "bin")
	// if err := filepath.Walk(systemGoPath, func(path string, info os.FileInfo, err error) error {
	// 	files = append(files, path)
	// 	return nil
	// }); err != nil {
	// 	log.Fatal(err)
	// }

	//
	files, err := ioutil.ReadDir(systemGoPath)
	if err != nil {
		log.Fatal(err)
	}
	var fNames []string
	for _, f := range files {
		fNames = append(fNames, f.Name())
	}
	fmt.Println("files: ", fNames)

	for _, v := range fNames {
		v = strings.TrimRight(v, ".exe")
		if v == "go" {
			fmt.Println("you already use system Go SDK")
			return
		} else if isGoVersionString(v) {
			fmt.Println("trying to use system Go SDK")
			// rename
			fmt.Println("v: ", v)
			curGoExe := filepath.Join(systemGoPath, v+".exe")
			goExe := filepath.Join(systemGoPath, "go.exe")
			if err := os.Rename(curGoExe, goExe); err != nil {
				log.Fatal(err)
			}

			// remove gopath go.exe
			if err := os.Remove(filepath.Join(goPath, "bin", "go.exe")); err != nil {
				fmt.Println("fail to remove go.exe in GOPATH: ", err)
			} else {
				fmt.Println("go.exe in GOPATH removed")
				fmt.Println("now using system Go SDK")
			}
			return
		}
	}
	log.Fatal("no Go SDK in GOROOT")
}

func useVersion(version string) { // ex) version == 1.15.2 (without "go")

	if version == systemGo {
		fmt.Println("use system version")
		useSystemGo()
		return
	}
	useVersion := "go" + version
	useExe := useVersion + ".exe"
	fmtV.Printf("wanted version: %s, exe: %s\n", useVersion, useExe)

	if !versionExist(useExe) {
		log.Fatal("%s is not installed version", useExe)
	}

	// get current go.exe info
	getPathCmd := exec.Command("where", "go.exe")
	getCurVersionCmd := exec.Command("go", "version")

	b, err := getPathCmd.Output()
	if err != nil {
		log.Fatal("getPathCmd: ", err)
	}
	curFilePath := strings.TrimSpace(string(b)) // needed to remove space
	fmtV.Printf("current version: %sAAA\n", curFilePath)

	// if exist then rename it. ex) go.exe -> go1.14.1.exe
	if fileExist(curFilePath) {
		// then we need file to rename
		b, err := getCurVersionCmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		versionOutput := strings.Split(string(b), " ")
		curVersionExe := versionOutput[2] + ".exe"
		fmt.Println("curVersionExe: ", curVersionExe)

		// rename
		dir := filepath.Dir(curFilePath)
		newFilePath := filepath.Join(dir, curVersionExe)
		if err := os.Rename(curFilePath, newFilePath); err != nil {
			log.Fatal("os.Rename: ", err)
		}
		if fileExist(newFilePath) {
			fmt.Println("rename succeeded")
		}

		// temp code for restore
		// if err := os.Rename(newFilePath, curFilePath); err != nil {
		// 	log.Fatal("os.Rename: ", err)
		// }
	}

	// then copy the go<required-version>.exe to go.exe
	filePath := filepath.Join(goPath, "bin", useExe)
	copyFile(filePath, renameToGo(filePath))

	fmt.Println("now we can use ", useVersion)
	// save usingVer
	usingVer = useVersion // ex) go1.13.2
	viper.Set(curUsingVerCfg, usingVer)

	getCurVersionCmd2 := exec.Command("go", "version")
	v, err := getCurVersionCmd2.Output()
	if err != nil {
		log.Fatal("getCurVersionCmd2:", err)
	}
	fmt.Println("changed version: ", string(v))
}

func use(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		curVer := getCurGoVersion()
		if _, isSystem := getCurGoExePath(); isSystem {
			fmt.Printf("Currently using %s, %s\n", systemGo, curVer)
			return
		}
		fmt.Println("Currently using", curVer)
		return
	}
	useVersion(args[0])
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
