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

// see how download works https://go.googlesource.com/dl/+/refs/heads/master/internal/version/version.go
// git repo: git clone https://go.googlesource.com/dl
func uninstall(cmd *cobra.Command, args []string) {

	noArgumentDisplayHelp(cmd, args)

	// version format check
	var validVersions []string
	for _, ver := range args {
		if isGoVersionString("go" + ver) {
			validVersions = append(validVersions, "go"+ver)
		} else if ver == systemGo {
			fmt.Printf("gvm will not uninstall %s go version\n", makeColorString(colorGreen, systemGo))
		} else {
			fmt.Printf("%s is not proper go version format\n", makeColorString(colorGreen, ver))
		}
	}

	if len(validVersions) == 0 {
		return
	}

	// check if installed version
	var uninstallVersions []string
	for _, ver := range validVersions {
		if isInstalledVersion(ver) {
			fmtV.Printf("%s is installed version. so, it will be removed\n", makeColorString(colorGreen, ver))
			uninstallVersions = append(uninstallVersions, ver)
		} else {
			fmt.Printf("%s is not installed version\n", makeColorString(colorGreen, ver))
		}
	}

	// Start uninstall
	fmt.Printf("--\nStart to remove versions: %v\n", uninstallVersions)

	for _, ver := range uninstallVersions {
		fmt.Printf("--\nStart to remove: %v\n", ver)

		verExe := ver + ".exe"
		// check if it's used version
		// then change to use system version
		if usingVer == ver {
			useVersion(systemGo)
		}

		// remove gopath\bin\go<version>.exe
		filePath := filepath.Join(goPath, "bin", verExe)
		if err := os.Remove(filePath); err != nil {
			log.Println(err)
		} else {
			fmt.Printf("Execution file: %s is removed\n", filePath)
		}

		// find home\sdk\<goversion> folder and remove ex) C:\Users\hsjeong\sdk\go1.13.3
		sdkPath, err := goroot(ver)
		if err != nil {
			log.Fatalf("%s: %v", ver, err)
		}

		// remove folder of the home\sdk\<version folder>
		if err := os.RemoveAll(sdkPath); err != nil {
			fmt.Printf("fail to remove sdk at %s: %v\n", sdkPath, err)
		} else {
			fmt.Printf("SDK: %s folder is removed or no such folder exist\n", sdkPath)
		}
	}
}

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall installed go SDK(s)",
	Long: `Uninstall installed go SDK(s)
You can uninstall one or more versions at once

ex) $ gvm uninstall 1.13.1 1.13.2`,
	// Args: cobra.MinimumNArgs(1),
	Run: uninstall,
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
