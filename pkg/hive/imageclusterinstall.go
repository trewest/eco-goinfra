package hive

import (
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	ibiv1alpha1 "github.com/openshift/image-based-install-operator/api/v1alpha1"
)

type ImageClusterInstallBuilder struct {
	Definition *ibiv1alpha1.ImageClusterInstall
	Object     *ibiv1alpha1.ImageClusterInstall
	errorMsg   string
	apiClient  *clients.Settings
}

func NewImageClusterInstallBuilder(apiClient *clients.Settings, name, nsname, imageSet string) {

}
