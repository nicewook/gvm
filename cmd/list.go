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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// getLocalList is
// 1) get the list and sort
// 2) then add systemGo at the end
func getLocalList() []string {
	dirPath := filepath.Join(goPath, "bin")
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal("ReadDir: ", err, dirPath)
	}

	var installedVersions []string
	// https://regex101.com/r/zxxWBl/1
	for _, file := range files {
		name := file.Name()
		name = strings.TrimRight(name, ".exe")
		if isGoVersionString(name) && name != "go" {
			installedVersions = append(installedVersions, name)
		}
	}
	sortGoSDKList(installedVersions)
	return append(installedVersions, systemGo)
}

func list(cmd *cobra.Command, args []string) {

	curVer := getCurGoVersion()
	_, isSystemGo := getCurGoExePath()
	if isSystemGo {
		curVer = systemGo
	}

	fmt.Println("locally installed Go SDK(s) list\n--")
	list := getLocalList()

	for _, ver := range list {
		if ver == curVer {
			fmt.Printf("%s %s\n", makeColorString(colorRed, ver), makeColorString(colorYellow, "*"))
		} else {
			fmt.Println(ver)
		}
	}
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: `List the locally installed Go SDK versions`,
	Long: `List the locally installed Go SDK versions.
It shows the currently using version as red with the asterisk.`,
	Run: list,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
