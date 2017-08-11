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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spiermar/burn/convert"
	"github.com/spiermar/burn/types"
)

var pretty bool
var folded bool
var html bool
var output string

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [flags] <input>",
	Short: "Convert performance profiles to a hierarchical data structure",
	Long: `
Convert performance profiles to a hierarchical data structure.

Examples:
  burn convert examples/out.perf
  burn convert --folded examples/out.perf-folded
  burn convert --html examples/out.perf
  burn convert --output=flame.json examples/out.perf
  burn convert --html --output=flame.html examples/out.perf
	`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath := string(args[0])

		rootNode := types.Node{"root", 0, make(map[string]*types.Node)}
		profile := types.Profile{rootNode, []string{}}

		if folded {
			profile = convert.ParseFolded(filePath)
		} else {
			profile = convert.ParsePerf(filePath)
		}

		b := []byte{}

		if pretty {
			err := (error)(nil)
			b, err = profile.RootNode.MarshalIndentJSON()
			if err != nil {
				panic(err)
			}
		} else {
			err := (error)(nil)
			b, err = profile.RootNode.MarshalJSON()
			if err != nil {
				panic(err)
			}
		}

		wr := os.Stdout

		if output != "" {
			err := (error)(nil)
			wr, err = os.Create(output)
			if err != nil {
				panic(err)
			}
			defer wr.Close()
		}

		if html {
			sep := strings.LastIndex(filePath, "/")
			filename := filePath[sep+1:]
			convert.GenerateHtml(wr, filename, string(b))
		} else {
			fmt.Fprintln(wr, string(b))
		}

		wr.Sync()
	},
}

func init() {
	RootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")
	convertCmd.PersistentFlags().BoolVarP(&folded, "folded", "f", false, "input is a folded stack")
	convertCmd.PersistentFlags().BoolVarP(&pretty, "pretty", "p", false, "json output is pretty printed")
	convertCmd.PersistentFlags().BoolVarP(&html, "html", "m", false, "output is a html flame graph")
	convertCmd.PersistentFlags().StringVar(&output, "output", "", "output file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
