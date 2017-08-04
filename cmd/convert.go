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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var Folded bool

type Node struct {
	Name     string
	Value    int
	Children map[string]Node
}

type Profile struct {
	Samples Node
	Stack   []string
	Name    string
}

func (n *Node) Add(frames []string, value int) {
	n.Value += value
	if len(frames) > 0 {
		head := frames[0]
		child, ok := n.Children[head]
		if !ok {
			child = Node{head, 0, make(map[string]Node)}
		}
		child.Add(frames[1:], value)
		n.Children[head] = child
	}
}

func (p *Profile) OpenStack(name string) {
	p.Stack = []string{}
	p.Name = name
}

func (p *Profile) CloseStack() {
	p.Stack = append([]string{p.Name}, p.Stack...)
	p.Samples.Add(p.Stack, 1)
	p.Stack = []string{}
	p.Name = ""
}

func (p *Profile) AddFrame(name string) {
	re, _ := regexp.Compile(`^\(`) // Skip process names
	if !re.MatchString(name) {
		name = strings.Replace(name, ";", ":", -1) // replace ; with :
		name = strings.Replace(name, "<", "", -1)  // remove '<'
		name = strings.Replace(name, ">", "", -1)  // remove '>'
		name = strings.Replace(name, "\\", "", -1) // remove '\'
		name = strings.Replace(name, "\"", "", -1) // remove '"'
		if index := strings.Index(name, "("); index != -1 {
			name = name[:index] // delete everything after '('
		}
		p.Stack = append([]string{name}, p.Stack...)
	}
}

func Parse(filename string) Profile {
	profile := Profile{Node{"root", 0, make(map[string]Node)}, []string{}, ""}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		var rComment = regexp.MustCompile(`^#`)                                    // Comment line
		var rEventRecordStart = regexp.MustCompile(`^(\S.+?)\s+(\d+)\/*(\d+)*\s+`) // Event record start
		var rStackLine = regexp.MustCompile(`^\s*(\w+)\s*(.+) \((\S*)\)`)          // Stack line
		var rEndStack = regexp.MustCompile(`^$`)                                   // End of stack

		switch {
		case rComment.MatchString(line):
			break
		case rEventRecordStart.MatchString(line):
			matches := rEventRecordStart.FindStringSubmatch(line)
			profile.OpenStack(matches[1])
			break
		case rStackLine.MatchString(line):
			matches := rStackLine.FindStringSubmatch(line)
			profile.AddFrame(matches[2])
			break
		case rEndStack.MatchString(line):
			profile.CloseStack()
			break
		default:
			panic("Don't know what to do with this line.")
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return profile
}

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
		profile := Parse(args[0])
		b, err := json.MarshalIndent(profile.Samples, "", "  ")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	},
}

func init() {
	RootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().BoolVarP(&Folded, "folded", "f", false, "Input is a folded stack.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
