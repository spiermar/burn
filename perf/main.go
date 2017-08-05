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

package perf

import (
	"bufio"
	"os"
	"regexp"

	"github.com/looplab/fsm"
	"github.com/spiermar/stacko/types"
)

func ParsePerf(filename string) types.Profile {
	rootNode := types.Node{"root", 0, make(map[string]*types.Node)}
	profile := types.Profile{rootNode, []string{}, ""}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	state := fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "read_comment", Src: []string{"start"}, Dst: "comment"},
			{Name: "open_stack", Src: []string{"start"}, Dst: "event"},
			{Name: "open_stack", Src: []string{"comment"}, Dst: "event"},
			{Name: "read_stack", Src: []string{"event"}, Dst: "open"},
			{Name: "close_stack", Src: []string{"open"}, Dst: "closed"},
			{Name: "open_stack", Src: []string{"closed"}, Dst: "event"},
			{Name: "finish", Src: []string{"closed"}, Dst: "end"},
		},
		fsm.Callbacks{
			"enter_event": func(e *fsm.Event) {

			},
			"enter_open": func(e *fsm.Event) {

			},
			"enter_closed": func(e *fsm.Event) {

			},
		},
	)

	reCommentLine := regexp.MustCompile(`^#`)                                    // Comment line
	reEventRecordStartLine := regexp.MustCompile(`^(\S.+?)\s+(\d+)\/*(\d+)*\s+`) // Event record start line
	reStackLine := regexp.MustCompile(`^\s*(\w+)\s*(.+) \((\S*)\)`)              // Stack line
	reEndStackLine := regexp.MustCompile(`^$`)                                   // End of stack line

	for scanner.Scan() {
		line := scanner.Text()
		current := state.Current()

		switch current {
		case "start":
			if reCommentLine.MatchString(line) {
				err := state.Event("read_comment")
				if err != nil {
					panic(err)
				}
			} else if matches := reEventRecordStartLine.FindStringSubmatch(line); matches != nil {
				err := state.Event("open_stack")
				if err != nil {
					panic(err)
				}
				profile.OpenStack(matches[1])
			} else {
				panic("Invalid format.")
			}
		case "comment":
			if reCommentLine.MatchString(line) {
				// Do nothing
			} else if matches := reEventRecordStartLine.FindStringSubmatch(line); matches != nil {
				err := state.Event("open_stack")
				if err != nil {
					panic(err)
				}
				profile.OpenStack(matches[1])
			} else {
				panic("Invalid format.")
			}
		case "event":
			if matches := reStackLine.FindStringSubmatch(line); matches != nil {
				err := state.Event("read_stack")
				if err != nil {
					panic(err)
				}
				profile.AddFrame(matches[2])
			} else {
				panic("Invalid format.")
			}
		case "open":
			if matches := reStackLine.FindStringSubmatch(line); matches != nil {
				profile.AddFrame(matches[2])
			} else if reEndStackLine.MatchString(line) {
				err := state.Event("close_stack")
				if err != nil {
					panic(err)
				}
				profile.CloseStack()
			} else {
				panic("Invalid format.")
			}
		case "closed":
			if matches := reEventRecordStartLine.FindStringSubmatch(line); matches != nil {
				err := state.Event("open_stack")
				if err != nil {
					panic(err)
				}
				profile.OpenStack(matches[1])
			} else {
				err := state.Event("finish")
				if err != nil {
					panic(err)
				}
			}
		case "end":
			break
		default:
			panic("Invalid state.")
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return profile
}
