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
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var Folded bool

type Node struct {
	Name     string
	Value    int
	Children map[string]*Node
}

type Profile struct {
	Samples Node
	Stack   []string
	Name    string
}

func (n *Node) Add(stackPtr *[]string, index int, value int) {
	n.Value += value
	if index >= 0 {
		head := (*stackPtr)[index]
		childPtr, ok := n.Children[head]
		if !ok {
			childPtr = &(Node{head, 0, make(map[string]*Node)})
			n.Children[head] = childPtr
		}
		childPtr.Add(stackPtr, index-1, value)
	}
}

func (p *Profile) OpenStack(name string) {
	p.Stack = []string{}
	p.Name = name
}

func (p *Profile) CloseStack() {
	p.Stack = append(p.Stack, p.Name)
	p.Samples.Add(&p.Stack, len(p.Stack)-1, 1)
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
		p.Stack = append(p.Stack, name)
	}
}

func Parse(filename string) Profile {
	rootNode := Node{"root", 0, make(map[string]*Node)}
	profile := Profile{rootNode, []string{}, ""}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		var reCommentLine = regexp.MustCompile(`^#`)                                    // Comment line
		var reEventRecordStartLine = regexp.MustCompile(`^(\S.+?)\s+(\d+)\/*(\d+)*\s+`) // Event record start line
		var reStackLine = regexp.MustCompile(`^\s*(\w+)\s*(.+) \((\S*)\)`)              // Stack line
		var reEndStackLine = regexp.MustCompile(`^$`)                                   // End of stack line

		if reCommentLine.MatchString(line) {
			// Do nothing
		} else if matches := reEventRecordStartLine.FindStringSubmatch(line); matches != nil {
			profile.OpenStack(matches[1])
		} else if matches := reStackLine.FindStringSubmatch(line); matches != nil {
			profile.AddFrame(matches[2])
		} else if reEndStackLine.MatchString(line) {
			profile.CloseStack()
		} else {
			panic("Don't know what to do with this line.")
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return profile
}

func ParseFolded(filename string) Profile {
	rootNode := Node{"root", 0, make(map[string]*Node)}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		sep := strings.LastIndex(line, " ")

		s := line[:sep]
		v := line[sep+1:]

		stack := strings.Split(s, ";")
		sort.Sort(sort.Reverse(sort.StringSlice(stack)))

		i, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		rootNode.Add(&stack, len(stack)-1, i)
	}

	if err := scanner.Err(); err != nil {

	}

	profile := Profile{rootNode, []string{}, ""}

	return profile
}

func (n *Node) MarshalJSON() ([]byte, error) {
	v := make([]Node, 0, len(n.Children))
	for _, value := range n.Children {
		v = append(v, *value)
	}

	return json.MarshalIndent(&struct {
		Name     string `json:"name"`
		Value    int    `json:"value"`
		Children []Node `json:"children"`
	}{
		Name:     n.Name,
		Value:    n.Value,
		Children: v,
	}, "", "  ")
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
		rootNode := Node{"root", 0, make(map[string]*Node)}
		profile := Profile{rootNode, []string{}, ""}

		if Folded {
			profile = ParseFolded(args[0])
		} else {
			profile = Parse(args[0])
		}

		b, err := profile.Samples.MarshalJSON()
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
