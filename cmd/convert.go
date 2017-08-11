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
	"github.com/spiermar/burn/convert"
	"github.com/spiermar/burn/types"
)

var pretty bool

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [flags] <input>",
	Short: "Convert a performance profile to a JSON",
	Long: `
Convert a performance profile to a JSON.

Examples:
  burn convert examples/out.perf
  burn convert --folded examples/out.perf-folded
	`,
	Run: func(cmd *cobra.Command, args []string) {
		rootNode := types.Node{"root", 0, make(map[string]*types.Node)}
		profile := types.Profile{rootNode, []string{}}

		if foldedStack {
			profile = convert.ParseFolded(args[0])
		} else {
			profile = convert.ParsePerf(args[0])
		}

		if pretty {
			b, err := profile.RootNode.MarshalIndentJSON()
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		} else {
			b, err := profile.RootNode.MarshalJSON()
			if err != nil {
				panic(err)
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
	convertCmd.PersistentFlags().BoolVarP(&foldedStack, "folded", "f", false, "input is a folded stack")
	convertCmd.PersistentFlags().BoolVarP(&pretty, "pretty", "p", false, "json output is pretty printed")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
