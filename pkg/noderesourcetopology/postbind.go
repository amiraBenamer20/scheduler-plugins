/*
Copyright 2023 The Kubernetes Authors.

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

package noderesourcetopology

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	// "sigs.k8s.io/scheduler-plugins/pkg/noderesourcetopology/logging"
	
	"github.com/amiraBenamer20/scheduler-plugins/pkg/noderesourcetopology/logging"
)

func (tm *TopologyMatch) PostBind(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) {
	lh := klog.FromContext(ctx).WithValues(logging.KeyPod, klog.KObj(pod), logging.KeyPodUID, logging.PodUID(pod), logging.KeyNode, nodeName)
	lh.V(4).Info(logging.FlowBegin)
	defer lh.V(4).Info(logging.FlowEnd)

	tm.nrtCache.PostBind(nodeName, pod)
}
