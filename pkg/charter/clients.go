package charter

import (
	"fmt"

	"github.com/jenkins-x-plugins/jx-charter/pkg/client/clientset/versioned"
	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"
)

// LazyCreateChartClient lazy creates the jx client if its not defined
func LazyCreateChartClient(client versioned.Interface) (versioned.Interface, error) {
	if client != nil {
		return client, nil
	}
	f := kubeclient.NewFactory()
	cfg, err := f.CreateKubeConfig()
	if err != nil {
		return client, fmt.Errorf("failed to get kubernetes config: %w", err)
	}
	client, err = versioned.NewForConfig(cfg)
	if err != nil {
		return client, fmt.Errorf("error building jx clientset: %w", err)
	}
	return client, nil
}
