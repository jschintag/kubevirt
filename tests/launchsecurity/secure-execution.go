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
	"kubevirt.io/kubevirt/tests/console"
	"kubevirt.io/kubevirt/tests/decorators"
	"kubevirt.io/kubevirt/tests/framework/kubevirt"
	"kubevirt.io/kubevirt/tests/libnode"
	"kubevirt.io/kubevirt/tests/libpod"
	"kubevirt.io/kubevirt/tests/libvmifact"
	"kubevirt.io/kubevirt/tests/libvmops"

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
	Context("Ensure cluster can run Secure Execution VMs", decorators.SecureExecution, decorators.SigCompute, func() {
		var vmi *kubevirtv1.VirtualMachineInstance
		BeforeEach(func() {
			vmi = libvmifact.NewFedora()
			// Enabling launchsecurity here won't prevent starting non-SE VMs
			vmi.Spec.Domain.LaunchSecurity = &kubevirtv1.LaunchSecurity{}

			// TODO: Mount the hostkey here via ConfigMap

			By("Launching a normal VM to convert it to Secure Execution")
			vmi = libvmops.RunVMIAndExpectLaunch(vmi, 240)

			// TODO: Convert VM to SE here

			By("Expecting the VirtualMachineInstance console")
			Expect(console.LoginToFedora(vmi)).To(Succeed())
		})

		It("Should schedule the VM on Secure Execution enabled nodes", func() {
			pod, err := libpod.GetPodByVirtualMachineInstance(vmi, vmi.Namespace)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Spec.NodeSelector).To(HaveKeyWithValue(kubevirtv1.SecureExecutionLabel, "true"))
		})

		It("Should launche a Secure Execution VM", func() {
			output, err := console.RunCommandAndStoreOutput(vmi, "cat /sys/firmware/uv/prot_virt_guest", 30)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal("1\n"))
		})
	})
})
