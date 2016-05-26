/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cmds

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/fabric8io/gitcontroller/util"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

type Result string

const (
	Success Result = "✔"
	Failure Result = "✘"

	// cmd flags
	Namespace = "namespace"
	Selector  = "selector"

	DataDir = "repos"
)

func showBanner() {
	ct.ChangeColor(ct.Blue, false, ct.None, false)
	fmt.Println(fabric8AsciiArt)
	ct.ResetColor()
}

func createListOpts(selector string) (*api.ListOptions, error) {
	listOpts := api.ListOptions{}
	if len(selector) > 0 {
		sel, err := labels.Parse(selector)
		if err != nil {
			return nil, err
		}
		util.Info("Using label selector: ")
		util.Successf("%v", sel)
		util.Info("\n")

		listOpts.LabelSelector = sel
	}
	return &listOpts, nil
}

func printError(err error) {
	if err != nil {
		util.Failuref("%v", err)
	}
	util.Blank()
}

const fabric8AsciiArt = `             [38;5;25m▄[38;5;25m▄▄[38;5;25m▄[38;5;25m▄[38;5;25m▄[38;5;235m▄[39m         [00m
             [48;5;25;38;5;25m█[48;5;235;38;5;235m█[48;5;235;38;5;235m█[48;5;25;38;5;25m█[48;5;25;38;5;25m█[48;5;25;38;5;25m█[48;5;235;38;5;235m█[49;39m         [00m
     [48;5;233;38;5;235m▄[48;5;235;38;5;25m▄[38;5;25m▄[38;5;25m▄[38;5;24m▄[38;5;25m▄[48;5;233;38;5;235m▄[49;39m [48;5;25;38;5;25m▄[48;5;235;38;5;24m▄[48;5;235;38;5;24m▄[48;5;25;38;5;25m▄[48;5;25;38;5;25m▄[48;5;25;38;5;25m▄[48;5;235;38;5;235m█[49;39m         [00m
     [48;5;235;38;5;235m█[48;5;24;38;5;24m█[48;5;25;38;5;25m█[48;5;24;38;5;24m█[48;5;235;38;5;235m█[48;5;25;38;5;25m█[48;5;235;38;5;235m█[49;39m [38;5;235m▀[38;5;235m▀▀▀▀▀[38;5;233m▀[39m [48;5;235;38;5;24m▄[48;5;235;38;5;25m▄[38;5;25m▄[38;5;25m▄[38;5;24m▄[48;5;235;38;5;25m▄[49;39m  [00m
     [48;5;235;38;5;235m▄[48;5;24;38;5;25m▄[48;5;25;38;5;25m▄[48;5;24;38;5;25m▄[48;5;235;38;5;25m▄[48;5;25;38;5;25m▄[48;5;235;38;5;235m▄[49;39m         [48;5;67;38;5;67m█[48;5;25;38;5;25m█[48;5;25;38;5;25m█[48;5;25;38;5;25m█[48;5;235;38;5;235m█[48;5;25;38;5;25m█[49;39m  [00m
   [38;5;233m▄[38;5;235m▄[48;5;235;38;5;24m▄[48;5;235;38;5;25m▄[49;38;5;235m▄[39m             [48;5;67;38;5;25m▄[48;5;25;38;5;25m▄[48;5;25;38;5;25m▄[48;5;25;38;5;25m▄[48;5;235;38;5;25m▄[48;5;25;38;5;25m▄[49;39m  [00m
   [38;5;235m▀[48;5;25;38;5;24m▄[48;5;24;38;5;25m▄[48;5;25;38;5;68m▄[48;5;24;38;5;25m▄[49;38;5;25m▄[39m      [38;5;235m▄[38;5;235m▄[38;5;17m▄[39m       [38;5;25m▄[38;5;25m▄[38;5;235m▄[39m [00m
    [38;5;23m▀[48;5;110;38;5;60m▄[48;5;110;38;5;254m▄[48;5;25;38;5;25m▄[48;5;25;38;5;25m▄[48;5;233;38;5;25m▄[49;38;5;235m▄[38;5;24m▄[38;5;25m▄[48;5;60;38;5;25m▄[48;5;67;38;5;25m▄[48;5;25;38;5;25m▄[48;5;25;38;5;110m▄[48;5;25;38;5;110m▄[48;5;25;38;5;25m▄[48;5;233;38;5;25m▄[49;39m   [38;5;233m▄[48;5;17;38;5;25m▄[48;5;25;38;5;25m▄[48;5;24;38;5;25m▄[48;5;25;38;5;24m▄[49;38;5;233m▀[39m[00m
      [38;5;60m▀[48;5;153;38;5;24m▄[48;5;68;38;5;110m▄[48;5;25;38;5;67m▄[48;5;25;38;5;25m▄[48;5;110;38;5;25m▄[48;5;67;38;5;255m▄[48;5;32;38;5;110m▄[48;5;68;38;5;110m▄[48;5;68;38;5;67m▄[48;5;25;38;5;110m▄[48;5;25;38;5;110m▄[38;5;110m▄[48;5;25;38;5;67m▄[48;5;24;38;5;67m▄[48;5;233;38;5;25m▄[49;38;5;25m▄[48;5;24;38;5;25m▄[48;5;25;38;5;25m█[38;5;25m▄[48;5;25;38;5;24m▄[49;38;5;17m▀[39m [00m
        [38;5;233m▀[38;5;24m▀[48;5;25;38;5;235m▄[48;5;25;38;5;25m█[48;5;153;38;5;110m▄[48;5;67;38;5;110m▄[48;5;252;38;5;255m▄[48;5;254;38;5;231m▄[48;5;254m▄[48;5;253;38;5;224m▄[48;5;252;38;5;255m▄[48;5;110;38;5;231m▄[48;5;110;38;5;231m▄[48;5;61;38;5;110m▄[48;5;25;38;5;25m▄[38;5;24m▄[48;5;25;38;5;233m▄[49;38;5;24m▀[39m   [00m
          [48;5;235;38;5;235m▄[48;5;25;38;5;25m█[48;5;67;38;5;67m▄[48;5;110;38;5;110m▄[48;5;255;38;5;255m▄[48;5;231;38;5;231m█[48;5;255;38;5;216m▄[48;5;223;38;5;209m▄[48;5;223;38;5;223m▄[48;5;231;38;5;231m█[48;5;231;38;5;231m▄[48;5;110;38;5;110m▄[48;5;235;38;5;235m▄[49;39m      [00m
          [48;5;235;38;5;235m▄[48;5;25;38;5;25m█[48;5;32;38;5;25m▄[48;5;67;38;5;25m▄[48;5;255;38;5;254m▄[48;5;231;38;5;255m▄[48;5;209;38;5;180m▄[48;5;209;38;5;223m▄[48;5;224;38;5;173m▄[48;5;231;38;5;255m▄[48;5;231;38;5;255m▄[48;5;110;38;5;67m▄[48;5;235;38;5;235m▄[49;39m      [00m
           [48;5;25;38;5;235m▄[48;5;25;38;5;25m▄[38;5;25m█[48;5;32m▄[48;5;110;38;5;25m▄[48;5;110;38;5;25m▄[48;5;110m▄[48;5;110m▄[48;5;110m▄[48;5;67m▄[48;5;25;38;5;25m▄[49;39m       [00m
            [48;5;25;38;5;25m▄[48;5;25;38;5;25m▄[38;5;25m▄[48;5;25;38;5;25m▄[49;38;5;235m▀[38;5;235m▀[48;5;25;38;5;25m▄[48;5;25;38;5;25m█[48;5;25;38;5;25m▄[49;39m        [00m
         [38;5;188m▄[48;5;242;38;5;188m▄[48;5;242;38;5;188m▄[48;5;25;38;5;250m▄[48;5;25;38;5;67m▄[48;5;67;38;5;67m▄[48;5;25;38;5;68m▄[48;5;250;38;5;25m▄[48;5;188;38;5;188m▄[48;5;25;38;5;110m▄[48;5;68;38;5;32m▄[48;5;25;38;5;67m▄[48;5;250;38;5;68m▄[48;5;188;38;5;251m▄[48;5;247;38;5;237m▄[49;39m     [00m
         [38;5;237m▀[38;5;242m▀[38;5;242m▀[38;5;247m▀[38;5;188m▀[38;5;251m▀[38;5;188m▀[38;5;188m▀[38;5;188m▀[38;5;188m▀[38;5;188m▀[38;5;188m▀[38;5;247m▀[38;5;237m▀[39m      [00m`
