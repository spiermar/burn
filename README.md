![burn](lizard.png)

# burn
burn is a CLI tool to convert performance profiles (perf_events, pprof, etc) to a JSON hierarchical data structure that can be consumed by [d3-flame-graph](https://github.com/spiermar/d3-flame-graph). burn can also generate a fully contained HTML flame graph from the same data.

## Getting Started

Just download _burn_ from github using `go get` or the _pre-build_ binary and use it!

### From go

Make sure you have [golang](https://golang.org/) installed and `GOPATH` correctly set.

```bash
$ go get github.com/spiermar/burn
$ burn $GOPATH/src/github.com/spiermar/burn/examples/out.perf
```

### From binary

Binaries are being provided for both `linux` and `darwin`, in `amd64` arch.

#### darwin/amd64

```bash
$ curl -L "https://dl.bintray.com/mspier/binaries/burn/1.0.1/darwin/amd64/burn" -o burn
$ chmod +x burn
$ ./burn <input_file>
```

#### linux/amd64

```bash
$ curl -L "https://dl.bintray.com/mspier/binaries/burn/1.0.1/linux/amd64/burn" -o burn
$ chmod +x burn
$ ./burn <input_file>
```

## Options

```
$ burn convert --help

Convert performance profiles to a hierarchical data structure.

Examples:
  burn convert examples/out.perf
  burn convert --type=folded examples/out.perf-folded
  burn convert --html examples/out.perf
  burn convert --output=flame.json examples/out.perf
  burn convert --html --output=flame.html examples/out.perf
  perf script | burn convert --html

Usage:
  burn convert [flags] (<input>)

Flags:
  -h, --help            help for convert
  -m, --html            output is a html flame graph
      --output string   output file
  -p, --pretty          json output is pretty printed
      --type string     input type (default "perf")

Global Flags:
      --config string       config file (default is $HOME/.burn.yaml)
      --cpuprofile string   write CPU profile to file
      --memprofile string   write heap profile to file
```

## Examples

Input and output examples can be found in the [examples](/examples) directory.

## Building from source

Make sure you have [golang](https://golang.org/) installed and `GOPATH` correctly set.

```bash
$ go get github.com/spiermar/burn
$ cd $GOPATH/src/github.com/spiermar/burn
$ go build
$ ./burn examples/out.perf
```

## Issues

For bugs, questions and discussions please use the [GitHub Issues](https://github.com/spiermar/burn/issues).

## Contributing

We love contributions! But in order to avoid total chaos, we have a few guidelines.

If you found a bug, have questions or feature requests, don't hesitate to open an [issue](https://github.com/spiermar/burn/issues).

If you're working on an issue, please comment on it so we can assign you to it.

If you have code to submit, follow the general pull request format. Fork the repo, make your changes, and submit a [pull request](https://github.com/spiermar/burn/pulls).

## License

Copyright 2017 Martin Spier. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.