package charter

import (
	"github.com/jenkins-x-plugins/jx-charter/pkg/client/clientset/versioned"
	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"
	"github.com/pkg/errors"
)

// LazyCreateChartClient lazy creates the jx client if its not defined
func LazyCreateChartClient(client versioned.Interface) (versioned.Interface, error) {
	if client != nil {
		return client, nil
	}
	f := kubeclient.NewFactory()
	cfg, err := f.CreateKubeConfig()
	if err != nil {
		return client, errors.Wrap(err, "failed to get kubernetes config")
	}
	client, err = versioned.NewForConfig(cfg)
	if err != nil {
		return client, errors.Wrap(err, "error building jx clientset")
	}
	return client, nil
}
