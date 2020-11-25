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
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// listallCmd represents the listall command
var listallCmd = &cobra.Command{
	Use:   "listall",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MaximumNArgs(1),
	Run:  listAll,
}

func getRemoteList() []string {
	doc, err := goquery.NewDocument(downloadPage)
	if err != nil {
		log.Fatal(err)
	}

	// use CSS selector found with the browser inspector
	// for each, use index and item
	var (
		count          int
		data           []string
		remoteVersions []string
	)

	doc.Find(".toggle").Each(func(index int, item *goquery.Selection) {
		ver, _ := item.Attr("id")
		// fmt.Printf("id: %s\n", ver)
		data = append(data, ver)
		count++
	})
	// fmt.Println("count:", count)

	doc.Find(".toggleVisible").Each(func(index int, item *goquery.Selection) {
		ver, _ := item.Attr("id")
		// fmt.Printf("id: %s\n", ver)
		data = append(data, ver)
		count++
	})
	fmtV.Println("total SDK count:", count)

	// regex for go version
	re := regexp.MustCompile(`(?m)go\d{0,2}.\d{0,2}.{0,1}\d{0,2}`)
	// https://regex101.com/r/zxxWBl/1
	for _, ver := range data {
		if re.MatchString(ver) {
			remoteVersions = append(remoteVersions, ver)
		}
	}

	// sort
	sortGoSDKList(remoteVersions)
	remoteVersions = append(remoteVersions, systemGo)
	return remoteVersions
}

// calcPrintCount is nested function for printCount
func calcPrintCount(restRows, restList int) int {
	intResult := restList / restRows
	floatResult := float32(restList) / float32(restRows)
	if float32(intResult) < floatResult {
		intResult++
	}
	return intResult
}

func columnPrint(list []string) {

	curVer := getCurGoVersion()
	_, isSystemGo := getCurGoExePath()
	if isSystemGo {
		curVer = systemGo
	}

	fmtV.Printf("current go version: %s\nis system go? %v\n\n", curVer, isSystemGo)
	// column print
	count := len(list)
	totalRows := 30

	for i := 0; i < totalRows; i++ {
		// calculate this time rowCount
		restRows := totalRows - i
		printCount := calcPrintCount(restRows, count)

		for j := 0; j < printCount; j++ {
			ver := list[i+j*30]

			_, found := find(goVerList, ver)
			if found {
				if ver == curVer {
					if err := colorPrint(Red, ver); err != nil {
						log.Fatal(err)
					}
				} else {
					if err := colorPrint(Green, ver); err != nil {
						log.Fatal(err)
					}
				}

			} else {
				verMsg := lefAlignString(ver)
				fmt.Print(verMsg)
			}
		}
		fmt.Println()
		count -= printCount
	}

}

func listAll(cmd *cobra.Command, args []string) {
	remoteList := getRemoteList()
	fmtV.Println("remote go SDK list\n--")

	columnPrint(remoteList)

	// for _, l := range remoteList {
	// 	_, found := find(goVerList, l)
	// 	if found {
	// 		if err := colorPrint(green, l); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	} else {
	// 		fmt.Println(l)
	// 	}
	// }
}

func init() {
	rootCmd.AddCommand(listallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
