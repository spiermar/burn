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
	"html/template"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spiermar/stacko/folded"
	"github.com/spiermar/stacko/perf"
	"github.com/spiermar/stacko/types"
)

var output string

const tpl = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <link rel="stylesheet" type="text/css" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
	<link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/gh/spiermar/d3-flame-graph@1.0.1/dist/d3.flameGraph.min.css">
	<style>
    /* Space out content a bit */
    body {
      padding-top: 20px;
      padding-bottom: 20px;
    }

    /* Custom page header */
    .header {
      padding-bottom: 20px;
      padding-right: 15px;
      padding-left: 15px;
      border-bottom: 1px solid #e5e5e5;
    }

    /* Make the masthead heading the same height as the navigation */
    .header h3 {
      margin-top: 0;
      margin-bottom: 0;
      line-height: 40px;
    }

    /* Customize container */
    .container {
      max-width: 990px;
    }
    </style>

    <title>{{.Name}}</title>

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
    <div class="container">
      <div class="header clearfix">
        <nav>
          <div class="pull-right">
            <form class="form-inline" id="form">
              <a class="btn" href="javascript: resetZoom();">Reset zoom</a>
              <a class="btn" href="javascript: clear();">Clear</a>
              <div class="form-group">
                <input type="text" class="form-control" id="term">
              </div>
              <a class="btn btn-primary" href="javascript: search();">Search</a>
            </form>
          </div>
        </nav>
        <h3 class="text-muted">{{.Name}}</h3>
      </div>
      <div id="chart">
      </div>
      <hr>
      <div id="details">
      </div>
    </div>

    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/d3/4.10.0/d3.min.js"></script>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/d3-tip/0.7.1/d3-tip.min.js"></script>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/lodash.js/4.17.4/lodash.min.js"></script>
	<script type="text/javascript" src="https://cdn.jsdelivr.net/gh/spiermar/d3-flame-graph@1.0.1/dist/d3.flameGraph.min.js"></script>
	<script type="text/javascript">
		var data = {{.Stack}};
	</script>
	<script type="text/javascript">
    var flameGraph = d3.flameGraph()
      .height(540)
      .width(960)
      .cellHeight(18)
      .transitionDuration(750)
      .transitionEase(d3.easeCubic)
      .sort(true)
      //Example to sort in reverse order
      //.sort(function(a,b){ return d3.descending(a.name, b.name);})
      .title("")
      .onClick(onClick);


    // Example on how to use custom tooltips using d3-tip.
    var tip = d3.tip()
      .direction("s")
      .offset([8, 0])
      .attr('class', 'd3-flame-graph-tip')
      .html(function(d) { return "name: " + d.data.name + ", value: " + d.data.value; });

    flameGraph.tooltip(tip);

    // Example on how to use custom labels
    // var label = function(d) {
    //  return "name: " + d.name + ", value: " + d.value;
    // }

    // flameGraph.label(label);

    d3.select("#chart")
      .datum(data)
      .call(flameGraph);

    document.getElementById("form").addEventListener("submit", function(event){
      event.preventDefault();
      search();
    });

    function search() {
      var term = document.getElementById("term").value;
      flameGraph.search(term);
    }

    function clear() {
      document.getElementById('term').value = '';
      flameGraph.clear();
    }

    function resetZoom() {
      flameGraph.resetZoom();
    }

    function onClick(d) {
      console.info("Clicked on " + d.data.name);
    }
    </script>
  </body>
</html>`

// htmlCmd represents the html command
var htmlCmd = &cobra.Command{
	Use:   "html",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := string(args[0])

		rootNode := types.Node{"root", 0, make(map[string]*types.Node)}
		profile := types.Profile{rootNode, []string{}}

		if foldedStack {
			profile = folded.ParseFolded(filename)
		} else {
			profile = perf.ParsePerf(filename)
		}

		b, err := profile.RootNode.MarshalJSON()
		if err != nil {
			panic(err)
		}

		t, err := template.New("flamegraph").Parse(tpl)
		if err != nil {
			panic(err)
		}

		sep := strings.LastIndex(filename, "/")

		data := struct {
			Name  string
			Stack template.JS
		}{
			Name:  filename[sep+1:],
			Stack: template.JS(b),
		}

		if output == "" {
			output = filename + ".html"
		}

		f, err := os.Create(output)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		err = t.Execute(f, data)
		if err != nil {
			panic(err)
		}

		f.Sync()
	},
}

func init() {
	RootCmd.AddCommand(htmlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// htmlCmd.PersistentFlags().String("foo", "", "A help for foo")
	htmlCmd.PersistentFlags().BoolVarP(&foldedStack, "folded", "f", false, "Input is a folded stack.")
	htmlCmd.PersistentFlags().StringVar(&output, "output", "", "Output file.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// htmlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
