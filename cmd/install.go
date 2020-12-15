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
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func alreadyInstalled(gover string) bool {

	if gover == getSystemGoVer() {
		gover = systemGo
	}

	_, found := find(goVerList, gover)
	if found {
		return true
	}
	return false
}

func isExistRemote(gover string) bool {
	for _, ver := range getRemoteList() {
		if gover == ver {
			return true
		}
	}
	return false
}

func installOneVersion(version string) {

	installVersion := "go" + version
	installURL := "golang.org/dl/" + installVersion
	downloadExe := goPath + "\\bin\\" + installVersion + ".exe"

	// check regex of the version name
	if isGoVersionString(installVersion) == false {
		fmt.Printf("%s is not proper go version format\n", makeColorString(colorRed, version))
		os.Exit(0)
	} else {
		fmtV.Printf("%s is good go version format\n", makeColorString(colorGreen, version))
	}
	// check the version exist or already downloaded
	if alreadyInstalled(installVersion) {
		fmt.Printf("%s is already installed", makeColorString(colorRed, installVersion))
		os.Exit(0)
	}
	if isExistRemote(installVersion) == false {
		fmt.Printf("%s is not existing version\n", makeColorString(colorRed, installVersion))
		os.Exit(0)
	}

	fmtV.Printf("start to install %s\n--\n", installVersion)
	fmtV.Printf("URL to download: %s\n", installURL)
	fmtV.Printf("download and execute %s\nfirst execution, it will downloaded go SDK, after then it will execute downloaded go SDK version\n--\n", downloadExe)

	// command wants to run
	// refer to: https://www.ardanlabs.com/blog/2020/04/modules-06-vendoring.html
	/*
		go get golang.org/dl/go1.13.10
		go1.13.10 download
	*/
	getCmd := exec.Command("go", "get", installURL)
	getCmd.Stdout = os.Stdout
	getCmd.Stderr = os.Stderr

	downloadCmd := exec.Command(downloadExe, "download")
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr

	if err := getCmd.Run(); err != nil {
		log.Fatal("getCmd:", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		var (
			timeoutSecond = 30
			elapsedSecond = 0
		)
		for !fileExist(downloadExe) {
			fmtV.Printf("not yet installed %s.exe file for download go SDK\n", installVersion)
			time.Sleep(1000 * time.Millisecond)

			elapsedSecond++
			if elapsedSecond >= timeoutSecond {
				log.Fatalf("Download file timeout. %d seconds", timeoutSecond)
			}
		}
		fmtV.Printf("%s.exe installed successfully\n", installVersion)
		wg.Done()
	}()
	wg.Wait()

	fmt.Printf("start downloading go SDK by executing %s.exe\n", installVersion)
	if err := downloadCmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func install(cmd *cobra.Command, args []string) {

	noArgumentDisplayHelp(cmd, args)

	for _, ver := range args {
		installOneVersion(ver)
	}
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Download and install your desired go version",
	Long: `Download and install any go version you want.
Just specify the version number without "go" prefix.

ex) $ gvm install 1.15.5 // install go v1.15.5

In detail, 
First, It installs go<version>.exe to "GOPATH\bin"
Second, It execute go<version>.exe, which will install SDK for the version. 
After then, when we execute go<version>.exe, it will run the installed go version.`,
	Run: install,
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
