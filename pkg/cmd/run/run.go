package run

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jenkins-x-plugins/jx-charter/pkg/apis/chart/v1alpha1"
	"github.com/jenkins-x-plugins/jx-charter/pkg/handlers"
	"github.com/jenkins-x-plugins/jx-charter/pkg/helmdecoder"
	"github.com/jenkins-x-plugins/jx-charter/pkg/rootcmd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	coreInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/yaml"
)

var (
	cmdLong = templates.LongDesc(`
		Runs the charter controller which watches helm Secrets and creates helm Chart CRDs
`)

	cmdExample = templates.Examples(`
		# watch for helm Secret resources and create/update the associated Chart CRDs
		%s run
	`)
)

// Options the options for this command
type Options struct {
	options.BaseOptions

	Port           string
	ResyncInterval time.Duration
	WatchNamespace string
	Namespace      string
	KubeClient     kubernetes.Interface

	CoreInformerFactory coreInformers.SharedInformerFactory
	HelmInformer        cache.SharedIndexInformer
	IsReady             atomic.Value
	Stop                chan struct{}
}

// NewCmdRun creates a command object for the command
func NewCmdRun() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Runs the charter controller which watches helm Secrets and creates helm Chart CRDs",
		Long:    cmdLong,
		Example: fmt.Sprintf(cmdExample, rootcmd.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}

	cmd.Flags().StringVarP(&o.WatchNamespace, "namespace", "n", "", "The kubernetes namespace to watch for helm Secrets")
	cmd.Flags().DurationVar(&o.ResyncInterval, "resync-interval", 1*time.Minute, "resync interval between full re-list operations")
	cmd.Flags().StringVar(&o.Port, "port", "8080", "port the health endpoint should listen on")

	o.BaseOptions.AddBaseFlags(cmd)
	return cmd, o
}

// Validate verifies things are setup correctly
func (o *Options) Validate() error {
	var err error

	o.KubeClient, o.Namespace, err = kube.LazyCreateKubeClientAndNamespace(o.KubeClient, o.Namespace)
	if err != nil {
		return errors.Wrapf(err, "failed to create kube client")
	}

	if o.CoreInformerFactory == nil {
		o.CoreInformerFactory = coreInformers.NewSharedInformerFactoryWithOptions(
			o.KubeClient,
			o.ResyncInterval,
			coreInformers.WithNamespace(o.WatchNamespace),
		)
	}

	return nil
}

func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate options")
	}

	o.HelmInformer = o.CoreInformerFactory.Core().V1().Secrets().Informer()

	o.HelmInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			r := obj.(*v1.Secret)
			upsertSecret(r)
			log.Logger().Debugf("handled Add event for helm secret %s", r.Name)
		},

		UpdateFunc: func(old, new interface{}) {
			r := new.(*v1.Secret)
			upsertSecret(r)
			log.Logger().Debugf("handled Update event for deployment %s", r.Name)
		},
	})

	log.Logger().Info("starting charter controller")

	// start the informers outside of the health endpoints
	go func() {
		o.Start()
	}()

	// health endpoint is used by kubernetes and changes to ready once informer caches are syncd
	o.startHealthEndpoint()
	return nil
}

func upsertSecret(r *v1.Secret) {
	if r == nil {
		return
	}
	release, err := helmdecoder.ConvertSecretToHelmRelease(r)
	if err != nil {
		log.Logger().Warnf("failed to decode Secret %s/%s due to %v\n", r.Namespace, r.Namespace, err.Error())
		return
	}
	if release == nil {
		return
	}

	ch := &v1alpha1.Chart{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.APIVersion,
			Kind:       v1alpha1.KindChart,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        release.Name,
			Namespace:   release.Namespace,
			Annotations: r.Annotations,
			Labels:      r.Labels,
		},
	}

	if release.Chart != nil && release.Chart.Metadata != nil {
		ch.Spec.Metadata = *release.Chart.Metadata
	}
	if release.Info != nil {
		ch.Status = v1alpha1.ToChartStatus(release.Info)
	}

	data, err := yaml.Marshal(ch)
	if err != nil {
		log.Logger().Warnf("failed to marshal Release for secret %s/%s due to %v\n", r.Namespace, r.Namespace, err.Error())
		return
	}

	fmt.Printf("\n%s\n", string(data))
}

func (o *Options) Start() {
	o.Stop = make(chan struct{})

	// Kubernetes serves an utility to handle API crashes
	defer runtime.HandleCrash()

	if o.CoreInformerFactory != nil {
		o.CoreInformerFactory.Start(o.Stop)
		if !cache.WaitForCacheSync(o.Stop, o.HelmInformer.HasSynced) {
			runtime.HandleError(fmt.Errorf("timed out waiting for deployment caches to sync"))
			return
		}
	}

	o.IsReady.Store(true)
	<-o.Stop
}

func (o *Options) startHealthEndpoint() {
	isReady := &o.IsReady
	r := handlers.Router(isReady)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", o.Port),
		Handler: r,
	}
	go func() {
		log.Logger().Fatal(srv.ListenAndServe())
	}()

	log.Logger().Infof("The service is ready to listen and serve.")

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Logger().Infof("Got SIGINT...")
	case syscall.SIGTERM:
		log.Logger().Infof("Got SIGTERM...")
	}

	log.Logger().Infof("The service is shutting down...")
	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Logger().Fatalf("failed to shutdown cleanly %v", err)
	}
	log.Logger().Infof("Done")
}
