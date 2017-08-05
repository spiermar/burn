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

type Profile struct {
	Samples Node
	Stack   []string
	Name    string
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