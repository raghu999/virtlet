/*
Copyright 2017 Mirantis

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

package manager

import (
	"fmt"

	"github.com/Mirantis/virtlet/pkg/metadata"
	"github.com/Mirantis/virtlet/pkg/tapmanager"
)

// recoverNetworkNamespaces recovers all the active VM network namespaces
// from previous Virtlet run by scanning the metadata store and starting
// dhcp server for each namespace that's still active
func recoverNetworkNamespaces(metadataStore metadata.Store, fdManager tapmanager.FDManager) (allErrors []error) {
	sandboxes, err := metadataStore.ListPodSandboxes(nil)
	if err != nil {
		allErrors = append(allErrors, err)
		return
	}

	for _, s := range sandboxes {
		psi, err := s.Retrieve()
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("can not retrieve PodSandboxInfo for sandbox id %q: %v", s.GetID(), err))
			continue
		}
		if psi == nil {
			allErrors = append(allErrors, fmt.Errorf("inconsistent database. Found pod %q sandbox but can not retrive its metadata", s.GetID()))
			continue
		}

		if err := fdManager.Recover(
			s.GetID(),
			tapmanager.GetFDPayload{
				ContainerSideNetwork: psi.ContainerSideNetwork,
				Description: &tapmanager.PodNetworkDesc{
					PodID:   s.GetID(),
					PodNs:   psi.Metadata.GetNamespace(),
					PodName: psi.Metadata.GetName(),
				},
			},
		); err != nil {
			allErrors = append(allErrors, fmt.Errorf("error recovering netns for %q pod: %v", s.GetID(), err))
		}
	}
	return
}
