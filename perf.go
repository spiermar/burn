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

package stacko

func ParsePerf(filename string) Profile {
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