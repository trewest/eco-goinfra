package hive

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	ibiv1alpha1 "github.com/openshift-kni/eco-goinfra/pkg/hive/extensions/v1alpha1"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	goclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ImageClusterInstallBuilder provides struct for the imageclusterinstall object containing connection to
// the cluster and the imageclusterinstall definitions.
type ImageClusterInstallBuilder struct {
	Definition *ibiv1alpha1.ImageClusterInstall
	Object     *ibiv1alpha1.ImageClusterInstall
	errorMsg   string
	apiClient  goclient.Client
}

// NewImageClusterInstallBuilder creates a new instance of ImageClusterInstallBuilder.
func NewImageClusterInstallBuilder(
	apiClient *clients.Settings, name, nsname, imageset string) *ImageClusterInstallBuilder {
	glog.V(100).Infof(
		"Initializing new imageclusterinstall structure with the following params: "+
			"name: %s, namespace: %s, imageset: %s",
		name, nsname, imageset)

	if apiClient == nil {
		return nil
	}

	builder := &ImageClusterInstallBuilder{
		apiClient: apiClient.Client,
		Definition: &ibiv1alpha1.ImageClusterInstall{
			Spec: ibiv1alpha1.ImageClusterInstallSpec{
				ImageSetRef: hivev1.ClusterImageSetReference{
					Name: imageset,
				},
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the imageclusterinstall is empty")

		builder.errorMsg = "imageclusterinstall 'name' cannot be empty"
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the imageclusterinstall is empty")

		builder.errorMsg = "imageclusterinstall 'nsname' cannot be empty"
	}

	if imageset == "" {
		glog.V(100).Infof("The imageset of the imageclusterinstall is empty")

		builder.errorMsg = "imageclusterinstall 'imageset' cannot be empty"
	}

	return builder
}

// PullImageClusterInstall retrieves an existing imageclusterinstall from the cluster.
func PullImageClusterInstall(apiClient *clients.Settings, name, nsname string) (*ImageClusterInstallBuilder, error) {
	glog.V(100).Infof(
		"Pulling existing imageclusterinstall with name %s from namespace %s", name, nsname)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient is nil")

		return nil, fmt.Errorf("apiClient cannot be nil")
	}

	builder := &ImageClusterInstallBuilder{
		apiClient: apiClient.Client,
		Definition: &ibiv1alpha1.ImageClusterInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the imageclusterinstall is empty")

		return nil, fmt.Errorf("imageclusterinstall 'name' cannot be empty")
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the imageclusterinstall is empty")

		return nil, fmt.Errorf("imageclusterinstall 'nsname' cannot be empty")
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("imageclusterinstall object %s does not exist in namespace %s", name, nsname)
	}

	builder.Definition = builder.Object

	return builder, nil
}

// WithHostname sets hostname of installed node.
func (builder *ImageClusterInstallBuilder) WithHostname(hostname string) *ImageClusterInstallBuilder {
	if valid, _ := builder.validate(); !valid {
		return nil
	}

	builder.Definition.Spec.Hostname = hostname

	return builder
}

// WithClusterDeployment links imageclusterinstall to an existing cluster deployment.
func (builder *ImageClusterInstallBuilder) WithClusterDeployment(name string) *ImageClusterInstallBuilder {
	if valid, _ := builder.validate(); !valid {
		return nil
	}

	if builder.Definition.Spec.ClusterDeploymentRef == nil {
		builder.Definition.Spec.ClusterDeploymentRef = &corev1.LocalObjectReference{}
	}

	builder.Definition.Spec.ClusterDeploymentRef.Name = name

	return builder
}

// WithExtraManifests includes manifests via configmap name.
func (builder *ImageClusterInstallBuilder) WithExtraManifests(name string) *ImageClusterInstallBuilder {
	if valid, _ := builder.validate(); !valid {
		return nil
	}

	builder.Definition.Spec.ExtraManifestsRefs =
		append(builder.Definition.Spec.ExtraManifestsRefs, corev1.LocalObjectReference{
			Name: name,
		})

	return builder
}

// WithMachineNetwork specifies the machine network where nodes will be installed.
func (builder *ImageClusterInstallBuilder) WithMachineNetwork(network string) *ImageClusterInstallBuilder {
	if valid, _ := builder.validate(); !valid {
		return nil
	}

	builder.Definition.Spec.MachineNetwork = network

	return builder
}

// WithSSHKey adds specified ssh key to authorized_keys of installed nodes.
func (builder *ImageClusterInstallBuilder) WithSSHKey(sshKey string) *ImageClusterInstallBuilder {
	if valid, _ := builder.validate(); !valid {
		return nil
	}

	builder.Definition.Spec.SSHKey = sshKey

	return builder
}

// WithVersion sets the seedimage version for the imageclusterinstall.
func (builder *ImageClusterInstallBuilder) WithVersion(version string) *ImageClusterInstallBuilder {
	if valid, _ := builder.validate(); !valid {
		return nil
	}

	builder.Definition.Spec.Version = version

	return builder
}

// GetCompletedCondition returns Completed condition from imageclusterinstall.
func (builder *ImageClusterInstallBuilder) GetCompletedCondition() (*hivev1.ClusterInstallCondition, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof("Getting Completed condition from imageclusterinstall %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		return nil, fmt.Errorf("cannot get condition from non-existent imageclusterinstall")
	}

	for _, condition := range builder.Object.Status.Conditions {
		if condition.Type == hivev1.ClusterInstallCompleted {
			return &condition, nil
		}
	}

	return nil, fmt.Errorf("cannot find completed condition in imageclusterinstall status")
}

// GetFailedCondition returns Failed condition from imageclusterinstall.
func (builder *ImageClusterInstallBuilder) GetFailedCondition() (*hivev1.ClusterInstallCondition, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof("Getting Failed condition from imageclusterinstall %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		return nil, fmt.Errorf("cannot get condition from non-existent imageclusterinstall")
	}

	for _, condition := range builder.Object.Status.Conditions {
		if condition.Type == hivev1.ClusterInstallFailed {
			return &condition, nil
		}
	}

	return nil, fmt.Errorf("cannot find failed condition in imageclusterinstall status")
}

// GetRequirementsMetCondition returns RequirementsMet condition from imageclusterinstall.
func (builder *ImageClusterInstallBuilder) GetRequirementsMetCondition() (*hivev1.ClusterInstallCondition, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof("Getting RequirementsMet condition from imageclusterinstall %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		return nil, fmt.Errorf("cannot get condition from non-existent imageclusterinstall")
	}

	for _, condition := range builder.Object.Status.Conditions {
		if condition.Type == hivev1.ClusterInstallRequirementsMet {
			return &condition, nil
		}
	}

	return nil, fmt.Errorf("cannot find requirements met condition in imageclusterinstall status")
}

// GetStoppedCondition returns Stopped condition from imageclusterinstall.
func (builder *ImageClusterInstallBuilder) GetStoppedCondition() (*hivev1.ClusterInstallCondition, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof("Getting Stopped condition from imageclusterinstall %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		return nil, fmt.Errorf("cannot get condition from non-existent imageclusterinstall")
	}

	for _, condition := range builder.Object.Status.Conditions {
		if condition.Type == hivev1.ClusterInstallStopped {
			return &condition, nil
		}
	}

	return nil, fmt.Errorf("cannot find stopped condition in imageclusterinstall status")
}

// Get fetches the defined imageclusterinstall from the cluster.
func (builder *ImageClusterInstallBuilder) Get() (*ibiv1alpha1.ImageClusterInstall, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof("Getting imageclusterinstall %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	imageClusterInstall := &ibiv1alpha1.ImageClusterInstall{}
	err := builder.apiClient.Get(context.TODO(), goclient.ObjectKey{
		Name:      builder.Definition.Name,
		Namespace: builder.Definition.Namespace,
	}, imageClusterInstall)

	if err != nil {
		return nil, err
	}

	return imageClusterInstall, err
}

// Create generates a imageclusterinstall on the cluster.
func (builder *ImageClusterInstallBuilder) Create() (*ImageClusterInstallBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Creating the imageclusterinstall %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	if !builder.Exists() {
		err = builder.apiClient.Create(context.TODO(), builder.Definition)
		if err == nil {
			builder.Object = builder.Definition
		}
	}

	return builder, err
}

// Update modifies an existing imageclusterinstall on the cluster.
func (builder *ImageClusterInstallBuilder) Update(force bool) (*ImageClusterInstallBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Updating imageclusterinstall %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	err := builder.apiClient.Update(context.TODO(), builder.Definition)

	if err != nil {
		if force {
			glog.V(100).Infof(
				msg.FailToUpdateNotification("imageclusterinstall", builder.Definition.Name, builder.Definition.Namespace))

			err := builder.Delete()

			if err != nil {
				glog.V(100).Infof(
					msg.FailToUpdateError("imageclusterinstall", builder.Definition.Name, builder.Definition.Namespace))

				return nil, err
			}

			return builder.Create()
		}
	}

	if err == nil {
		builder.Object = builder.Definition
	}

	return builder, err
}

// Delete removes an imageclusterinstall from the cluster.
func (builder *ImageClusterInstallBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting the imageclusterinstall %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		return fmt.Errorf("imageclusterinstall cannot be deleted because it does not exist")
	}

	err := builder.apiClient.Delete(context.TODO(), builder.Definition)

	if err != nil {
		return fmt.Errorf("cannot delete imageclusterinstall: %w", err)
	}

	builder.Object = nil

	return nil
}

// Exists checks if the defined imageclusterinstall has already been created.
func (builder *ImageClusterInstallBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof("Checking if imageclusterinstall %s exists in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	builder.Object, err = builder.Get()

	return err == nil || !k8serrors.IsNotFound(err)
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *ImageClusterInstallBuilder) validate() (bool, error) {
	resourceCRD := "ImageClusterInstall"

	if builder == nil {
		glog.V(100).Infof("The %s builder is uninitialized", resourceCRD)

		return false, fmt.Errorf("error: received nil %s builder", resourceCRD)
	}

	if builder.Definition == nil {
		glog.V(100).Infof("The %s is undefined", resourceCRD)

		builder.errorMsg = msg.UndefinedCrdObjectErrString(resourceCRD)
	}

	if builder.apiClient == nil {
		glog.V(100).Infof("The %s builder apiclient is nil", resourceCRD)

		builder.errorMsg = fmt.Sprintf("%s builder cannot have nil apiClient", resourceCRD)
	}

	if builder.errorMsg != "" {
		glog.V(100).Infof("The %s builder has error message: %s", resourceCRD, builder.errorMsg)

		return false, fmt.Errorf(builder.errorMsg)
	}

	return true, nil
}
