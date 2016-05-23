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
	k8sclient "k8s.io/kubernetes/pkg/client/unversioned"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

const (
	consoleMetadataUrl           = "io/fabric8/apps/console/maven-metadata.xml"
	baseConsoleUrl               = "io/fabric8/apps/console/%[1]s/console-%[1]s-kubernetes.json"
	consoleKubernetesMetadataUrl = "io/fabric8/apps/console-kubernetes/maven-metadata.xml"
	baseConsoleKubernetesUrl     = "io/fabric8/apps/console-kubernetes/%[1]s/console-kubernetes-%[1]s-kubernetes.json"

	devopsTemplatesDistroUrl = "io/fabric8/forge/distro/distro/%[1]s/distro-%[1]s-templates.zip"
	devOpsMetadataUrl        = "io/fabric8/forge/distro/distro/maven-metadata.xml"

	kubeflixTemplatesDistroUrl = "io/fabric8/kubeflix/distro/distro/%[1]s/distro-%[1]s-templates.zip"
	kubeflixMetadataUrl        = "io/fabric8/kubeflix/distro/distro/maven-metadata.xml"

	zipkinTemplatesDistroUrl = "io/fabric8/zipkin/packages/distro/%[1]s/distro-%[1]s-templates.zip"
	zipkinMetadataUrl        = "io/fabric8/zipkin/packages/distro/maven-metadata.xml"

	iPaaSTemplatesDistroUrl = "io/fabric8/ipaas/distro/distro/%[1]s/distro-%[1]s-templates.zip"
	iPaaSMetadataUrl        = "io/fabric8/ipaas/distro/distro/maven-metadata.xml"

	Fabric8SCC    = "fabric8"
	PrivilegedSCC = "privileged"
	RestrictedSCC = "restricted"

	runFlag             = "app"
	versioniPaaSFlag    = "version-ipaas"
	versionDevOpsFlag   = "version-devops"
	versionKubeflixFlag = "version-kubeflix"
	versionZipkinFlag   = "version-zipkin"
	mavenRepoFlag       = "maven-repo"
	dockerRegistryFlag  = "docker-registry"

	typeLabel          = "type"
	teamTypeLabelValue = "team"
)

type createFunc func(c *k8sclient.Client, f *cmdutil.Factory, name string) (Result, error)

func NewCmdRun(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "watches the Deployments and ReplicationControllers for changes to their git repositories and to perform rolling upgrades when they change",
		Long:  `watches the Deployments and ReplicationControllers for changes to their git repositories and to perform rolling upgrades when they change`,
		PreRun: func(cmd *cobra.Command, args []string) {
			showBanner()
		},
		Run: func(cmd *cobra.Command, args []string) {
			c, cfg := client.NewClient(f)
			ns, _, _ := f.DefaultNamespace()
			util.Info("Running gitcontroller on the ")
			util.Success(string(util.TypeOfMaster(c)))
			util.Info(" installation at ")
			util.Success(cfg.Host)
			util.Info(" in namespace ")
			util.Successf("%s\n\n", ns)

		},
	}
	cmd.PersistentFlags().StringP("domain", "d", defaultDomain(), "The domain name to append to the service name to access web applications")
	/*
		cmd.PersistentFlags().String("api-server", "", "overrides the api server url")
		cmd.PersistentFlags().String(runFlag, "", "The name of the fabric8 app to startup. e.g. use `--app=cd-pipeline` to run the main CI/CD pipeline app")
		cmd.PersistentFlags().Bool(templatesFlag, true, "Should the standard Fabric8 templates be installed?")
		cmd.PersistentFlags().Bool(consoleFlag, true, "Should the Fabric8 console be deployed?")
	*/
	return cmd
}
