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
	"github.com/fabric8io/gitcontroller/git"
	"github.com/fabric8io/gitcontroller/util"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/api"

	"fmt"
	tapi "github.com/openshift/origin/pkg/template/api"
	tapiv1 "github.com/openshift/origin/pkg/template/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions"
	k8sclient "k8s.io/kubernetes/pkg/client/unversioned"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"os"
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

	/*
		rclist, err := c.ReplicationControllers(ns).List(listOpts)
		if err != nil {
			return err
		}
		for _, rc := range rclist.Items {
			err = checkRC(c, &rc)
			if err != nil {
				return err
			}
		}
	*/
	return nil
}

func checkRC(c *k8sclient.Client, rc *api.ReplicationController) error {
	template := rc.Spec.Template
	if template != nil {
		result, err := checkPodSpec(c, rc.Kind, &rc.ObjectMeta, &template.Spec)
		if err != nil {
			return err
		}
		if result {
			return fmt.Errorf("TODO update RC")
		}
	}
	return nil
}

func checkDeployment(c *k8sclient.Client, dep *extensions.Deployment, ns string) error {
	template := dep.Spec.Template
	result, err := checkPodSpec(c, dep.Kind, &dep.ObjectMeta, &template.Spec)
	if err != nil {
		return err
	}
	if result {
		_, err = c.Extensions().Deployments(ns).Update(dep)
		return err

	}
	return nil
}

func checkPodSpec(c *k8sclient.Client, kind string, metadata *api.ObjectMeta, podSpec *api.PodSpec) (bool, error) {
	result := false
	if podSpec != nil {
		for _, volume := range podSpec.Volumes {
			source := volume.VolumeSource
			gitRepo := source.GitRepo
			if gitRepo != nil {
				util.Infof("Found git repo for volume: %v\n", volume.Name)
				repo := gitRepo.Repository
				revision := gitRepo.Revision

				newrevision, err := checkIfGitUpdated(repo, revision, metadata, kind, volume.Name)
				if err != nil {
					return false, err
				}
				if newrevision != revision {
					util.Infof("Revision updated from %s to %s for volume: %v namespace: %s name: %s\n", revision, newrevision, volume.Name, metadata.Namespace, metadata.Name)
					gitRepo.Revision = newrevision
					result = true
				} else {
					util.Infof("Revision still at %s for volume: %v namespace: %s name: %s\n", newrevision, volume.Name, metadata.Namespace, metadata.Name)
				}
			}
		}
	}
	return result, nil
}

func checkIfGitUpdated(repo string, revision string, metadata *api.ObjectMeta, kind string, volumeName string) (string, error) {
	path := DataDir + "/" + metadata.Namespace + "/" + kind + "/" + metadata.Name + "/" + volumeName

	err := os.MkdirAll(path, 0777)
	if err != nil {
		return "", err
	}
	hasGit, err := exists(path + "/.git")
	if err != nil {
		return "", err
	}
	if !hasGit {
		err = git.GitClone(repo, path)
	} else {
		err = git.GitPull(path)
	}
	if err != nil {
		return "", err
	}
	return git.GitLatestCommitSince(path, revision)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
