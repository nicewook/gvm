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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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
	Args: cobra.ExactArgs(1),
	Run:  use,
}

func versionExist(file string) bool {
	filePath := filepath.Join(gopath, "bin", file)
	fmt.Println("want to use version of: ", filePath)
	return fileExist(filePath)
}

func use(cmd *cobra.Command, args []string) {
	// we want to use this version
	useVersion := "go" + args[0]
	useExe := useVersion + ".exe"
	// fmt.Printf("version: %s, exe: %s\n", useVersion, useExe)

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
	filePath := filepath.Join(gopath, "bin", useExe)
	copyFile(filePath, renameToGo(filePath))

	fmt.Println("now we can use ", useVersion)
	getCurVersionCmd2 := exec.Command("go", "version")
	v, err := getCurVersionCmd2.Output()
	if err != nil {
		log.Fatal("getCurVersionCmd2:", err)
	}
	fmt.Println(string(v))
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
