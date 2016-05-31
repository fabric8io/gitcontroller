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
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/watch"

	"fmt"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	k8sclient "k8s.io/kubernetes/pkg/client/unversioned"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"os"
	"strconv"
	"time"
)

type createFunc func(c *k8sclient.Client, f *cmdutil.Factory, name string) (Result, error)

type GitWatcher struct {
	ListOpts    *api.ListOptions
	Namespace   string
	KubeClient  *k8sclient.Client
	Deployments map[string]*extensions.Deployment

	// channels
	CheckC chan *extensions.Deployment
}

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
			ns := os.Getenv(NamespaceEnvVar)
			if len(ns) <= 0 {
				ns, _, _ = f.DefaultNamespace()
			}
			util.Info("Running gitcontroller on the ")
			util.Success(string(util.TypeOfMaster(c)))
			util.Info(" installation at ")
			util.Success(cfg.Host)
			util.Info(" in namespace ")
			util.Successf("%s\n\n", ns)

			selector := cmd.Flags().Lookup(Selector).Value.String()
			pollTime := cmd.Flags().Lookup(PollTime).Value.String()
			pollSeconds, err := strconv.Atoi(pollTime)
			if err != nil {
				printError(err)
			}

			listOpts, err := createListOpts(selector)
			if err != nil {
				printError(err)
			}

			watcher := GitWatcher{
				ListOpts:    listOpts,
				Namespace:   ns,
				KubeClient:  c,
				Deployments: make(map[string]*extensions.Deployment),
				CheckC:      make(chan *extensions.Deployment),
			}

			err = watcher.loadDeployments()
			if err != nil {
				printError(err)
			}

			go watcher.processLoop(pollSeconds)

			fmt.Println("Starting k8s watch loop")
			watcher.watchLoop()

		},
	}
	cmd.PersistentFlags().Int32(PollTime, 60, "Number of seconds between polls of git repositories")
	return cmd
}

func (watcher *GitWatcher) processLoop(seconds int) {
	util.Infof("Starting to process ticker every %d second(s) and watching kubernetes events\n", seconds)
	tickChan := time.NewTicker(time.Second * time.Duration(seconds)).C

	for {
		select {
		case dep := <-watcher.CheckC:
			watcher.Deployments[toKey(dep)] = dep
			watcher.checkDependency(dep)

		case <-tickChan:
			for _, dep := range watcher.Deployments {
				watcher.checkDependency(dep)
			}
		}
	}
}

func (watcher *GitWatcher) checkDependency(dep *extensions.Deployment) {
	checkDeployment(watcher.KubeClient, dep, watcher.Namespace)
}

// lets load the initial dependencies
func (watcher *GitWatcher) loadDeployments() error {
	c := watcher.KubeClient
	ns := watcher.Namespace
	deplist, err := c.Extensions().Deployments(ns).List(*watcher.ListOpts)
	if err != nil {
		return err
	}
	for _, dep := range deplist.Items {
		watcher.Deployments[toKey(&dep)] = &dep
	}
	return nil
}

// watches dependencies; when they change we send them to the channel
func (watcher *GitWatcher) watchLoop() {
	c := watcher.KubeClient
	ns := watcher.Namespace
	w, err := c.Extensions().Deployments(ns).Watch(*watcher.ListOpts)
	if err != nil {
		printError(err)
	}
	kubectl.WatchLoop(w, func(e watch.Event) error {
		o := e.Object
		switch o := o.(type) {
		case *extensions.Deployment:
			watcher.CheckC <- o
		default:
			util.Warnf("Unknown watch object type %v\n", o)
		}
		return nil
	})
}
