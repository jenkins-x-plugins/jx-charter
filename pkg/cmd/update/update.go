package update

import (
	"context"
	"fmt"

	"github.com/jenkins-x-plugins/jx-charter/pkg/charter"
	"github.com/jenkins-x-plugins/jx-charter/pkg/client/clientset/versioned"
	"github.com/jenkins-x-plugins/jx-charter/pkg/rootcmd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"

	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	cmdLong = templates.LongDesc(`
		Creates or Updates helm Chart CRDs from the helm Secrets
`)

	cmdExample = templates.Examples(`
		# creates or updates any missing helm Chart resources from the helm Secrets
		%s update
	`)
)

// Options the options for this command
type Options struct {
	options.BaseOptions

	WatchNamespace string
	Namespace      string
	KubeClient     kubernetes.Interface
	ChartClient    versioned.Interface
}

// NewCmdUpdate creates a command object for the command
func NewCmdUpdate() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Creates or Updates helm Chart CRDs from the helm Secrets",
		Long:    cmdLong,
		Example: fmt.Sprintf(cmdExample, rootcmd.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}

	cmd.Flags().StringVarP(&o.WatchNamespace, "namespace", "n", "", "The kubernetes namespace to look for helm Secrets")

	o.BaseOptions.AddBaseFlags(cmd)
	return cmd, o
}

// Validate verifies things are setup correctly
func (o *Options) Validate() error {
	var err error

	o.KubeClient, o.Namespace, err = kube.LazyCreateKubeClientAndNamespace(o.KubeClient, o.Namespace)
	if err != nil {
		return fmt.Errorf("failed to create kube client: %w", err)
	}

	o.ChartClient, err = charter.LazyCreateChartClient(o.ChartClient)
	if err != nil {
		return fmt.Errorf("failed to create chart client: %w", err)
	}

	return nil
}

func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate options: %w", err)
	}

	ctx := context.TODO()

	listOptions := metav1.ListOptions{}
	list, err := o.KubeClient.CoreV1().Secrets(o.WatchNamespace).List(ctx, listOptions)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("could not list Secrets in namespace %s: %w", o.WatchNamespace, err)
	}
	if list == nil {
		return nil
	}

	for i := range list.Items {
		r := &list.Items[i]
		err := charter.UpsertChartFromSecret(ctx, o.ChartClient, r)
		if err != nil {
			return fmt.Errorf("failed to process secret %s: %w", r.Name, err)
		}
	}
	return nil
}
