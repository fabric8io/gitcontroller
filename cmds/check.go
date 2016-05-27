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
	"github.com/fabric8io/gitcontroller/client"
	"github.com/fabric8io/gitcontroller/util"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/api"

	tapi "github.com/openshift/origin/pkg/template/api"
	tapiv1 "github.com/openshift/origin/pkg/template/api/v1"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func NewCmdCheck(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Checks the resources to see if they need to be updated with a new git version",
		Long:  `Checks the resources to see if they need to be updated with a new git version`,
		PreRun: func(cmd *cobra.Command, args []string) {
			tapi.AddToScheme(api.Scheme)
			tapiv1.AddToScheme(api.Scheme)
			showBanner()
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := checkResources(f, cmd)
			if err != nil {
				printError(err)
			}
		},
	}
	return cmd
}

func checkResources(f *cmdutil.Factory, cmd *cobra.Command) error {
	c, _ := client.NewClient(f)
	ns := cmd.Flags().Lookup(Namespace).Value.String()
	if len(ns) <= 0 {
		ns, _, _ = f.DefaultNamespace()
	}
	selector := cmd.Flags().Lookup(Selector).Value.String()

	util.Info("Checking git repositories of resources in namespace ")
	util.Success(ns)
	util.Info(" with selector")
	util.Success(selector)
	util.Info("\n")

	listOpts, err := createListOpts(selector)
	if err != nil {
		return err
	}
	deplist, err := c.Extensions().Deployments(ns).List(*listOpts)
	if err != nil {
		return err
	}
	for _, dep := range deplist.Items {
		err = checkDeployment(c, &dep, ns)
		if err != nil {
			return err
		}
	}
	return nil
}
