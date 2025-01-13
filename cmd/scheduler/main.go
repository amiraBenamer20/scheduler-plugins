/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"

	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // for rest client metric registration
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/kubernetes/cmd/kube-scheduler/app"

	"github.com/amiraBenamer20/scheduler-plugins/pkg/capacityscheduling"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/coscheduling"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/networkaware/networkoverhead"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/networkaware/topologicalsort"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/network-cost-aware/networkcost"//Amira
	"github.com/amiraBenamer20/scheduler-plugins/pkg/network-cost-aware/topologicalcnsort"//Amira
	"github.com/amiraBenamer20/scheduler-plugins/pkg/noderesources"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/noderesourcetopology"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/podstate"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/preemptiontoleration"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/qos"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/sysched"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/trimaran/loadvariationriskbalancing"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/trimaran/lowriskovercommitment"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/trimaran/peaks"
	"github.com/amiraBenamer20/scheduler-plugins/pkg/trimaran/targetloadpacking"

	// Ensure scheme package is initialized.
	_ "github.com/amiraBenamer20/scheduler-plugins/apis/config/scheme"
)

func main() {
	// Register custom plugins to the scheduler framework.
	// Later they can consist of scheduler profile(s) and hence
	// used by various kinds of workloads.
	command := app.NewSchedulerCommand(
		app.WithPlugin(capacityscheduling.Name, capacityscheduling.New),
		app.WithPlugin(coscheduling.Name, coscheduling.New),
		app.WithPlugin(loadvariationriskbalancing.Name, loadvariationriskbalancing.New),
		app.WithPlugin(networkoverhead.Name, networkoverhead.New),
		app.WithPlugin(topologicalsort.Name, topologicalsort.New),
		app.WithPlugin(networkcost.Name, networkcost.New),//Amira
		app.WithPlugin(topologicalcnsort.Name, topologicalcnsort.New),//Amira
		app.WithPlugin(noderesources.AllocatableName, noderesources.NewAllocatable),
		app.WithPlugin(noderesourcetopology.Name, noderesourcetopology.New),
		app.WithPlugin(preemptiontoleration.Name, preemptiontoleration.New),
		app.WithPlugin(targetloadpacking.Name, targetloadpacking.New),
		app.WithPlugin(lowriskovercommitment.Name, lowriskovercommitment.New),
		app.WithPlugin(sysched.Name, sysched.New),
		app.WithPlugin(peaks.Name, peaks.New),
		// Sample plugins below.
		// app.WithPlugin(crossnodepreemption.Name, crossnodepreemption.New),
		app.WithPlugin(podstate.Name, podstate.New),
		app.WithPlugin(qos.Name, qos.New),
	)

	code := cli.Run(command)
	os.Exit(code)
}
