package run

import (
	"context"
	"flag"
	"fmt"
	"github.com/jenkins-x-plugins/jx-charter/pkg/handlers"
	"github.com/jenkins-x-plugins/jx-charter/pkg/rootcmd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	coreInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
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

	Port       string
	Namespace  string
	KubeClient kubernetes.Interface

	CoreInformerFactory coreInformers.SharedInformerFactory
	HelmInformer        cache.SharedIndexInformer
	IsReady             *atomic.Value
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

	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", "", "The kubernetes Namespace to watch for PipelineRun and PipelineActivity resources. Defaults to the current namespace")
	flag.StringVar(&o.Port, "port", "8080", "port the health endpoint should listen on")

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
 	fmt.Printf("got Secret %s/%s\n", r.Namespace, r.Namespace)
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
	isReady := o.IsReady
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
