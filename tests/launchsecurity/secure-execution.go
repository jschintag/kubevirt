/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright The KubeVirt Authors
 *
 */

package launchsecurity

import (
	"kubevirt.io/kubevirt/tests/decorators"
	"kubevirt.io/kubevirt/tests/framework/kubevirt"
	"kubevirt.io/kubevirt/tests/libnode"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kubevirtv1 "kubevirt.io/api/core/v1"
)

var _ = Describe("[sig-compute]IBM Secure Execution", decorators.SecureExecution, decorators.SigCompute, func() {
	Context("Node Labels", func() {
		It("Should have nodes with Secure Execution Label", func() {
			virtclient := kubevirt.Client()
			nodes := libnode.GetAllSchedulableNodes(virtclient)
			hasNodeWithSELabel := false
			for _, node := range nodes.Items {
				if _, ok := node.Labels[kubevirtv1.SecureExecutionLabel]; ok {
					hasNodeWithSELabel = true
					break
				}
			}
			Expect(hasNodeWithSELabel).To(BeTrue())
		})
	})
})
