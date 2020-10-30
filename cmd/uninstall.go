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
	"path/filepath"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run:  uninstall,
}

// see how download works https://go.googlesource.com/dl/+/refs/heads/master/internal/version/version.go
// git repo: git clone https://go.googlesource.com/dl
func uninstall(cmd *cobra.Command, args []string) {
	fmt.Println("args: ", args)
	uninstallVers := args
	_ = uninstallVers

	fmt.Println("uninstall called to uninstall: ", uninstallVers)

	for _, v := range uninstallVers {
		ver := "go" + v
		verExe := ver + ".exe"
		// check if it's used version
		// then change to use system version
		if usingVer == v {
			useVersion("system")
		}

		// remove gopath\bin\go<version>.exe
		filePath := filepath.Join(goPath, "bin", verExe)
		if err := os.Remove(filePath); err != nil {
			log.Println(err)
		} else {
			fmt.Println(filePath, " is removed")
		}

		// find home\sdk\<goversion> folder and remove ex) C:\Users\hsjeong\sdk\go1.13.3
		sdkPath, err := goroot(ver)
		if err != nil {
			log.Fatalf("%s: %v", ver, err)
		}

		// remove folder of the home\sdk\<version folder>
		if err := os.RemoveAll(sdkPath); err != nil {
			log.Printf("fail to remove sdk path: %v\n", err)
		} else {
			fmt.Println(sdkPath, " folder is removed or no such folder exist")
		}
	}
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
