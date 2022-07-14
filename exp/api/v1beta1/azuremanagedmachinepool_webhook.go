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

package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/Azure/go-autorest/autorest/to"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/maps"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedmachinepools,verbs=create;update,versions=v1beta1,name=default.azuremanagedmachinepools.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (m *AzureManagedMachinePool) Default(client client.Client) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}
	m.Labels[LabelAgentPoolMode] = m.Spec.Mode

	if m.Spec.Name == nil || *m.Spec.Name == "" {
		m.Spec.Name = &m.Name
	}
}

//+kubebuilder:webhook:verbs=update;delete,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedmachinepools,versions=v1beta1,name=validation.azuremanagedmachinepools.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (m *AzureManagedMachinePool) ValidateCreate(client client.Client) error {
	validators := []func() error{
		m.validateMaxPods,
		m.validateOSType,
		m.validateName,
	}

	var errs []error
	for _, validator := range validators {
		if err := validator(); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (m *AzureManagedMachinePool) ValidateUpdate(oldRaw runtime.Object, client client.Client) error {
	old := oldRaw.(*AzureManagedMachinePool)
	var allErrs field.ErrorList

	if m.Spec.SKU != old.Spec.SKU {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "SKU"),
				m.Spec.SKU,
				"field is immutable"))
	}

	if old.Spec.OSType != nil {
		// Prevent OSType modification if it was already set to some value
		if m.Spec.OSType == nil {
			// unsetting the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "OSType"),
					m.Spec.OSType,
					"field is immutable, unsetting is not allowed"))
		} else if *m.Spec.OSType != *old.Spec.OSType {
			// changing the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "OSType"),
					*m.Spec.OSType,
					"field is immutable"))
		}
	}

	if old.Spec.OSDiskSizeGB != nil {
		// Prevent OSDiskSizeGB modification if it was already set to some value
		if m.Spec.OSDiskSizeGB == nil {
			// unsetting the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "OSDiskSizeGB"),
					m.Spec.OSDiskSizeGB,
					"field is immutable, unsetting is not allowed"))
		} else if *m.Spec.OSDiskSizeGB != *old.Spec.OSDiskSizeGB {
			// changing the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "OSDiskSizeGB"),
					*m.Spec.OSDiskSizeGB,
					"field is immutable"))
		}
	}

	if !reflect.DeepEqual(m.Spec.Taints, old.Spec.Taints) {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "Taints"),
				m.Spec.Taints,
				"field is immutable"))
	}

	// custom headers are immutable
	oldCustomHeaders := maps.FilterByKeyPrefix(old.ObjectMeta.Annotations, azure.CustomHeaderPrefix)
	newCustomHeaders := maps.FilterByKeyPrefix(m.ObjectMeta.Annotations, azure.CustomHeaderPrefix)
	if !reflect.DeepEqual(oldCustomHeaders, newCustomHeaders) {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("metadata", "annotations"),
				m.ObjectMeta.Annotations,
				fmt.Sprintf("annotations with '%s' prefix are immutable", azure.CustomHeaderPrefix)))
	}

	if !ensureStringSlicesAreEqual(m.Spec.AvailabilityZones, old.Spec.AvailabilityZones) {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "AvailabilityZones"),
				m.Spec.AvailabilityZones,
				"field is immutable"))
	}

	if old.Spec.MaxPods != nil {
		// Prevent MaxPods modification if it was already set to some value
		if m.Spec.MaxPods == nil {
			// unsetting the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "MaxPods"),
					m.Spec.MaxPods,
					"field is immutable, unsetting is not allowed"))
		} else if *m.Spec.MaxPods != *old.Spec.MaxPods {
			// changing the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "MaxPods"),
					*m.Spec.MaxPods,
					"field is immutable"))
		}
	}

	if old.Spec.OsDiskType != nil {
		// Prevent OSDiskType modification if it was already set to some value
		if m.Spec.OsDiskType == nil || to.String(m.Spec.OsDiskType) == "" {
			// unsetting the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "OsDiskType"),
					m.Spec.OsDiskType,
					"field is immutable, unsetting is not allowed"))
		} else if *m.Spec.OsDiskType != *old.Spec.OsDiskType {
			// changing the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "OsDiskType"),
					m.Spec.OsDiskType,
					"field is immutable"))
		}
	}

	if old.Spec.EnableUltraSSD != nil {
		// Prevent EnabledUltraSSD modification if it was already set to some value
		if m.Spec.EnableUltraSSD == nil {
			// unsetting the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "EnableUltraSSD"),
					m.Spec.EnableUltraSSD,
					"field is immutable, unsetting is not allowed"))
		} else if *m.Spec.EnableUltraSSD != *old.Spec.EnableUltraSSD {
			// changing the field is not allowed
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "EnableUltraSSD"),
					m.Spec.EnableUltraSSD,
					"field is immutable"))
		}
	} else {
		if m.Spec.EnableUltraSSD != nil {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "EnableUltraSSD"),
					m.Spec.EnableUltraSSD,
					"field is immutable, unsetting is not allowed"))
		}
	}
	if len(allErrs) != 0 {
		return apierrors.NewInvalid(GroupVersion.WithKind("AzureManagedMachinePool").GroupKind(), m.Name, allErrs)
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (m *AzureManagedMachinePool) ValidateDelete(client client.Client) error {
	return nil
}

func (m *AzureManagedMachinePool) validateMaxPods() error {
	if m.Spec.MaxPods != nil {
		if to.Int32(m.Spec.MaxPods) < 10 || to.Int32(m.Spec.MaxPods) > 250 {
			return field.Invalid(
				field.NewPath("Spec", "MaxPods"),
				m.Spec.MaxPods,
				"MaxPods must be between 10 and 250")
		}
	}

	return nil
}

func (r *AzureManagedMachinePool) validateOSType() error {
	if r.Spec.Mode == string(NodePoolModeSystem) {
		if r.Spec.OSType != nil && *r.Spec.OSType != azure.LinuxOS {
			return field.Forbidden(
				field.NewPath("Spec", "OSType"),
				"System node pooll must have OSType 'Linux'")
		}
	}

	return nil
}

func (r *AzureManagedMachinePool) validateName() error {
	if r.Spec.OSType != nil && *r.Spec.OSType == azure.WindowsOS {
		if len(r.Name) > 6 {
			return field.Invalid(
				field.NewPath("Name"),
				r.Name,
				"Windows agent pool name can not be longer than 6 characters.")
		}
	}

	return nil
}

func ensureStringSlicesAreEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	m := map[string]bool{}
	for _, v := range a {
		m[v] = true
	}

	for _, v := range b {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}
