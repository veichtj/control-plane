package internal

import (
	"context"

	reconcilerApi "github.com/kyma-incubator/reconciler/pkg/keb"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/ptr"
	"github.com/kyma-project/control-plane/components/provisioner/pkg/gqlschema"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	SCMigrationComponentName          = "sc-migration"
	BTPOperatorComponentName          = "btp-operator"
	HelmBrokerComponentName           = "helm-broker"
	ServiceCatalogComponentName       = "service-catalog"
	ServiceCatalogAddonsComponentName = "service-catalog-addons"
	ServiceManagerComponentName       = "service-manager-proxy"

	// BTP Operator overrides keys
	BTPOperatorClientID     = "manager.secret.clientid"
	BTPOperatorClientSecret = "manager.secret.clientsecret"
	BTPOperatorURL          = "manager.secret.url"    // deprecated, for btp-operator v0.2.0
	BTPOperatorSMURL        = "manager.secret.sm_url" // for btp-operator v0.2.3
	BTPOperatorTokenURL     = "manager.secret.tokenurl"
)

type ClusterIDGetter func() (string, error)

func DisableServiceManagementComponents(r ProvisionerInputCreator) {
	r.DisableOptionalComponent(SCMigrationComponentName)
	r.DisableOptionalComponent(HelmBrokerComponentName)
	r.DisableOptionalComponent(ServiceCatalogComponentName)
	r.DisableOptionalComponent(ServiceCatalogAddonsComponentName)
	r.DisableOptionalComponent(ServiceManagerComponentName)
	r.DisableOptionalComponent(BTPOperatorComponentName)
}

func getBTPOperatorProvisioningOverrides(creds *ServiceManagerOperatorCredentials) []*gqlschema.ConfigEntryInput {
	return []*gqlschema.ConfigEntryInput{
		{
			Key:    BTPOperatorClientID,
			Value:  creds.ClientID,
			Secret: ptr.Bool(true),
		},
		{
			Key:    BTPOperatorClientSecret,
			Value:  creds.ClientSecret,
			Secret: ptr.Bool(true),
		},
		{
			Key:   BTPOperatorURL,
			Value: creds.ServiceManagerURL,
		},
		{
			Key:   BTPOperatorSMURL,
			Value: creds.ServiceManagerURL,
		},
		{
			Key:   BTPOperatorTokenURL,
			Value: creds.URL,
		},
	}
}

func getBTPOperatorUpdateOverrides(creds *ServiceManagerOperatorCredentials, clusterId string) []*gqlschema.ConfigEntryInput {
	if clusterId == "" {
		return []*gqlschema.ConfigEntryInput{}
	}
	return []*gqlschema.ConfigEntryInput{
		{
			Key:   "cluster.id",
			Value: clusterId,
		},
	}
}

func GetBTPOperatorReconcilerOverrides(creds *ServiceManagerOperatorCredentials, clusterIdGetter ClusterIDGetter) ([]reconcilerApi.Configuration, error) {
	id, err := clusterIdGetter()
	if err != nil {
		return nil, err
	}
	provisioning := getBTPOperatorProvisioningOverrides(creds)
	update := getBTPOperatorUpdateOverrides(creds, id)
	all := append(provisioning, update...)
	var config []reconcilerApi.Configuration
	for _, c := range all {
		secret := false
		if c.Secret != nil {
			secret = *c.Secret
		}
		rc := reconcilerApi.Configuration{Key: c.Key, Value: c.Value, Secret: secret}
		config = append(config, rc)
	}
	return config, nil
}

func CreateBTPOperatorProvisionInput(r ProvisionerInputCreator, creds *ServiceManagerOperatorCredentials) {
	overrides := getBTPOperatorProvisioningOverrides(creds)
	r.AppendOverrides(BTPOperatorComponentName, overrides)
}

func GetClusterIDWithKubeconfig(kubeconfig string) ClusterIDGetter {
	return func() (string, error) {
		cfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
		if err != nil {
			return "", err
		}
		cs, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return "", err
		}
		cm, err := cs.CoreV1().ConfigMaps("kyma-system").Get(context.Background(), "cluster-info", metav1.GetOptions{})
		if k8serrors.IsNotFound(err) {
			return "", nil
		}
		if err != nil {
			return "", err
		}
		return cm.Data["id"], nil
	}
}
