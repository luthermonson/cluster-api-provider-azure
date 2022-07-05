/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capierrors "sigs.k8s.io/cluster-api/errors"
)

const (
	// LabelAgentPoolMode represents mode of an agent pool. Possible values include: System, User.
	LabelAgentPoolMode = "azuremanagedmachinepool.infrastructure.cluster.x-k8s.io/agentpoolmode"

	// NodePoolModeSystem represents mode system for azuremachinepool.
	NodePoolModeSystem NodePoolMode = "System"

	// NodePoolModeUser represents mode user for azuremachinepool.
	NodePoolModeUser NodePoolMode = "User"
)

// NodePoolMode enumerates the values for agent pool mode.
type NodePoolMode string

// AzureManagedMachinePoolSpec defines the desired state of AzureManagedMachinePool.
type AzureManagedMachinePoolSpec struct {

	// Mode - represents mode of an agent pool. Possible values include: System, User.
	// +kubebuilder:validation:Enum=System;User
	Mode string `json:"mode"`

	// SKU is the size of the VMs in the node pool.
	SKU string `json:"sku"`

	// OSDiskSizeGB is the disk size for every machine in this agent pool.
	// If you specify 0, it will apply the default osDisk size according to the vmSize specified.
	OSDiskSizeGB *int32 `json:"osDiskSizeGB,omitempty"`

	// ProviderIDList is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`

	// Scaling specifies the autoscaling parameters for the node pool.
	// +optional
	Scaling *ManagedMachinePoolScaling `json:"scaling,omitempty"`

	// OSType specifies the virtual machine operating system. Default to Linux. Possible values include: 'Linux', 'Windows'
	// +kubebuilder:validation:Enum=Linux;Windows
	// +optional
	OSType *string `json:"osType,omitempty"`
}

// ManagedMachinePoolScaling specifies scaling options.
type ManagedMachinePoolScaling struct {
	// MinSize specifies the minimum number node count.
	MinSize *int32 `json:"minSize"`

	// MazSize specifies the maximum number node count.
	MaxSize *int32 `json:"maxSize"`
}

// AzureManagedMachinePoolStatus defines the observed state of AzureManagedMachinePool.
type AzureManagedMachinePoolStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Replicas is the most recently observed number of replicas.
	// +optional
	Replicas int32 `json:"replicas"`

	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	ErrorReason *capierrors.MachineStatusError `json:"errorReason,omitempty"`

	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=azuremanagedmachinepools,scope=Namespaced,categories=cluster-api,shortName=ammp
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// AzureManagedMachinePool is the Schema for the azuremanagedmachinepools API.
type AzureManagedMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureManagedMachinePoolSpec   `json:"spec,omitempty"`
	Status AzureManagedMachinePoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AzureManagedMachinePoolList contains a list of AzureManagedMachinePools.
type AzureManagedMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureManagedMachinePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureManagedMachinePool{}, &AzureManagedMachinePoolList{})
}
