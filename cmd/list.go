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
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: `List the locally installed go SDK versions`,
	// Long:  `List the locally installed go SDK versions`,
	Run: list,
}

// getLocalList is
// 1) get the list and sort
// 2) then add systemGo at the end
func getLocalList() []string {
	dirPath := filepath.Join(goPath, "bin")
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal("ReadDir: ", err, dirPath)
	}

	// https://regex101.com/r/zxxWBl/3
	var re = regexp.MustCompile(`(?m)go\d{0,2}.\d{0,2}.{0,1}\d{0,2}.exe`)
	var installedVersions []string
	// https://regex101.com/r/zxxWBl/1
	for _, file := range files {
		name := file.Name()
		if re.MatchString(name) {
			name = strings.TrimRight(name, ".exe")
			// fmt.Println(name)
			installedVersions = append(installedVersions, name)
		}
	}
	sortGoSDKList(installedVersions)
	return append(installedVersions, systemGo)
}

func list(cmd *cobra.Command, args []string) {
	fmt.Println("locally installed go SDK list\n--")
	list := getLocalList()

	curVer := getCurGoVersion()
	_, isSystemGo := getCurGoExePath()
	if isSystemGo {
		curVer = systemGo
	}

	for _, ver := range list {
		if ver == curVer {
			if err := colorPrint(Red, addStar(ver)); err != nil {
				log.Fatal(err)
			}
			fmt.Println()
		} else {
			fmt.Println(ver)
		}
	}
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
