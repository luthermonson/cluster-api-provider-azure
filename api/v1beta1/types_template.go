/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta1

import "github.com/pkg/errors"

type AzureClusterTemplateResourceSpec struct {
	AzureClusterClassSpec `json:",inline"`

	// NetworkSpec encapsulates all things related to Azure network.
	// +optional
	NetworkSpec NetworkTemplateSpec `json:"networkSpec,omitempty"`

	// BastionSpec encapsulates all things related to the Bastions in the cluster.
	// +optional
	BastionSpec BastionTemplateSpec `json:"bastionSpec,omitempty"`
}

type NetworkTemplateSpec struct {
	NetworkClassSpec `json:",inline"`

	// Vnet is the configuration for the Azure virtual network.
	// +optional
	Vnet VnetTemplateSpec `json:"vnet,omitempty"`

	// Subnets is the configuration for the control-plane subnet and the node subnet.
	// +optional
	Subnets SubnetTemplatesSpec `json:"subnets,omitempty"`

	// APIServerLB is the configuration for the control-plane load balancer.
	// +optional
	APIServerLB LoadBalancerClassSpec `json:"apiServerLB,omitempty"`

	// NodeOutboundLB is the configuration for the node outbound load balancer.
	// +optional
	NodeOutboundLB *LoadBalancerClassSpec `json:"nodeOutboundLB,omitempty"`

	// ControlPlaneOutboundLB is the configuration for the control-plane outbound load balancer.
	// This is different from APIServerLB, and is used only in private clusters (optionally) for enabling outbound traffic.
	// +optional
	ControlPlaneOutboundLB *LoadBalancerClassSpec `json:"controlPlaneOutboundLB,omitempty"`
}

// GetControlPlaneSubnetTemplate returns the cluster control plane subnet template.
func (n *NetworkTemplateSpec) GetControlPlaneSubnetTemplate() (SubnetTemplateSpec, error) {
	for _, sn := range n.Subnets {
		if sn.Role == SubnetControlPlane || sn.Role == SubnetAll {
			return sn, nil
		}
	}
	return SubnetTemplateSpec{}, errors.Errorf("no subnet template found with role %s or %s", SubnetControlPlane, SubnetAll)
}

// UpdateControlPlaneSubnet updates the cluster control plane subnet.
func (n *NetworkTemplateSpec) UpdateControlPlaneSubnetTemplate(subnet SubnetTemplateSpec) {
	for i, sn := range n.Subnets {
		if sn.Role == SubnetControlPlane || sn.Role == SubnetAll {
			n.Subnets[i] = subnet
		}
	}
}

type VnetTemplateSpec struct {
	VnetClassSpec `json:",inline"`

	// Peerings defines a list of peerings of the newly created virtual network with existing virtual networks.
	// +optional
	Peerings VnetPeeringsTemplateSpec `json:"peerings,omitempty"`
}

type VnetPeeringsTemplateSpec []VnetPeeringClassSpec

type SubnetTemplateSpec struct {
	SubnetClassSpec `json:",inline"`

	// SecurityGroup defines the NSG (network security group) that should be attached to this subnet.
	// +optional
	SecurityGroup SecurityGroupClass `json:"securityGroup,omitempty"`

	// NatGateway associated with this subnet.
	// +optional
	NatGateway NatGatewayClassSpec `json:"natGateway,omitempty"`
}

func (s SubnetTemplateSpec) IsNatGatewayEnabled() bool {
	return s.NatGateway.Name != ""
}

type SubnetTemplatesSpec []SubnetTemplateSpec

type BastionTemplateSpec struct {
	// +optional
	AzureBastion *AzureBastionTemplateSpec `json:"azureBastion,omitempty"`
}

type AzureBastionTemplateSpec struct {
	// +optional
	Subnet SubnetTemplateSpec `json:"subnet,omitempty"`
}
