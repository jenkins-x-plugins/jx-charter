package main

import (
	"os"

	"github.com/jenkins-x-plugins/jx-charter/pkg/apis/chart"
	"github.com/jenkins-x-plugins/jx-charter/pkg/apis/chart/v1alpha1"
	"github.com/jenkins-x/jx-api/v4/pkg/schemagen"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

var (
	resourceKinds = []schemagen.ResourceKind{
		{
			APIVersion: chart.GroupAndVersion,
			Name:       "chart",
			Resource:   &v1alpha1.Chart{},
		},
	}
)

func main() {
	out := "schema"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	err := schemagen.GenerateSchemas(resourceKinds, out)
	if err != nil {
		log.Logger().Errorf("failed: %v", err)
		os.Exit(1)
	}
	log.Logger().Infof("completed the plugin generator")
	os.Exit(0)
}
