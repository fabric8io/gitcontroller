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
	"github.com/fabric8io/gitcontroller/util"
	tapi "github.com/openshift/origin/pkg/template/api"
	tapiv1 "github.com/openshift/origin/pkg/template/api/v1"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/api"
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
			util.Info("Checking git repositories")
		},
	}
	return cmd
}
