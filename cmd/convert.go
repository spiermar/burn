// Copyright © 2017 Martin Spier <spiermar@gmail.com>
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

	"github.com/spf13/cobra"
	"github.com/spiermar/stacko/folded"
	"github.com/spiermar/stacko/perf"
	"github.com/spiermar/stacko/types"
)

var Folded bool
var Pretty bool

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootNode := types.Node{"root", 0, make(map[string]*types.Node)}
		profile := types.Profile{rootNode, []string{}, ""}

		if Folded {
			profile = folded.ParseFolded(args[0])
		} else {
			profile = perf.ParsePerf(args[0])
		}

		if Pretty {
			b, err := profile.Samples.MarshalIndentJSON()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(b))
		} else {
			b, err := profile.Samples.MarshalJSON()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(b))
		}
	},
}

func init() {
	RootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().BoolVarP(&Folded, "folded", "f", false, "Input is a folded stack.")
	RootCmd.PersistentFlags().BoolVarP(&Pretty, "pretty", "p", false, "JSON output is pretty printed.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
