package schematests

import (
	"github.com/jenkins-x-plugins/jx-charter/pkg/apis/chart/v1alpha1"
	"github.com/jenkins-x/jx-api/v4/pkg/util"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestValidateExternalSecret(t *testing.T) {
	t.Parallel()

	path := filepath.Join("test_data", "chart.yaml")
	data, err := ioutil.ReadFile(path)
	require.NoError(t, err, "failed to load %s", path)

	deploy := &v1alpha1.Chart{}
	err = yaml.Unmarshal(data, deploy)
	require.NoError(t, err, "failed to unmarshal %s", path)

	results, err := util.ValidateYaml(deploy, data)
	t.Logf("got results %#v\n", results)

	require.NoError(t, err, "should not have failed to validate yaml file %s", path)

	require.Empty(t, results, "should not have validation errors for file %s", path)
}
