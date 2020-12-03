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
	"log"
	"os"

	"github.com/spf13/cobra"

	goverbose "github.com/paulvollmer/go-verbose"
)

const (
	downloadPage = "https://golang.org/dl"

	usingVerCfg = "usingVer"
	systemGo    = "system"
)

var (
	cfgFile string
	goRoot  string
	goPath  string

	// cache values
	usingVer  string
	goVerList []string
	verbose   bool

	fmtV *goverbose.Verbose
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gvm",
	Short: `"gvm" is go language version manager for Windows`,
	Long: `"gvm" is go language version manager for Windows. 
You can install/uninstall any version you want to use 
and choose the version what you want to use.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gvm.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// log settings
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	/*
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homeDir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home directory with name ".gvm" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(".gvm")
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	*/
	// Set verbose https://play.golang.org/p/RoRcgJV0pDV
	fmtV = goverbose.New(os.Stdout, verbose)

	// JHS custom config
	goRoot = os.Getenv("GOROOT")
	goPath = os.Getenv("GOPATH")
	fmtV.Println("GOROOT:", goRoot)
	fmtV.Println("GOPATH:", goPath)

	// init go version list installed in local
	goVerList = getLocalList()
}
