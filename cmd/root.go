// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cpuProfile string
var memProfile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "burn command [flags] <input>",
	Short: `A tool to convert performance profiles to hierarchical data structures that 
can be visualized as flame graphs`,
	Long: `
burn 🔥 is a CLI tool to convert performance profiles to hierarchical data
structures that can be visualized as flame graphs.

burn can also generate a self contained html flame graphs from the same data.

Examples:
  burn convert examples/out.perf
  burn convert --folded examples/out.perf-folded
  burn html examples/out.perf
	`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cpuProfile != "" {
			f, err := os.Create(cpuProfile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if cpuProfile != "" {
			pprof.StopCPUProfile()
		}
		if memProfile != "" {
			f, err := os.Create(memProfile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.burn.yaml)")
	RootCmd.PersistentFlags().StringVar(&cpuProfile, "cpuprofile", "", "write CPU profile to file")
	RootCmd.PersistentFlags().StringVar(&memProfile, "memprofile", "", "write heap profile to file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".burn" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".burn")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
