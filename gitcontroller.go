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
package main

import (
	commands "github.com/fabric8io/gitcontroller/cmds"
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func main() {
	cmds := &cobra.Command{
		Use:   "gitcontroller",
		Short: "gitcontroller performs rolling updates on applications using git to store configurations",
		Long: `gitcontroller performs rolling updates on applications using git to store configurations.\n
								Find more information at http://fabric8.io.`,
		Run: runHelp,
	}

	cmds.PersistentFlags().BoolP("yes", "y", false, "assume yes")
	cmds.PersistentFlags().String(commands.Namespace, "", "namespace to query")
	cmds.PersistentFlags().String(commands.Selector, "", "label selector to query")

	f := cmdutil.NewFactory(nil)
	f.BindFlags(cmds.PersistentFlags())

	cmds.AddCommand(commands.NewCmdRun(f))
	cmds.AddCommand(commands.NewCmdCheck(f))
	cmds.AddCommand(commands.NewCmdVersion())

	cmds.Execute()
}
