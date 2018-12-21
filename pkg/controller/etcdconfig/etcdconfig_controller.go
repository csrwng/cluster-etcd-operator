/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package etcdconfig

import (
	"context"
	"fmt"
	"os"
	"time"

	etcdv1 "github.com/openshift/cluster-etcd-operator/pkg/apis/etcd/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	kubeinformers "k8s.io/client-go/informers"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterclientset "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	clusterinformers "sigs.k8s.io/cluster-api/pkg/client/informers_generated/externalversions"

	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/cluster-etcd-operator/pkg/controller/assets"
)

const (
	operatorNamespace = "cluster-etcd-operator"
)

var (
	log                   = logf.Log.WithName("cluster-etcd-operator")
	defaultResyncInterval = 10 * time.Minute
)

// Add creates a new Etcd Operator and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	reconciler, informers, err := newReconciler(mgr)
	if err != nil {
		return err
	}
	return add(mgr, reconciler)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) (reconcile.Reconciler, error) {
	externalDNSImage := os.Getenv("IMAGE")
	if len(externalDNSImage) == 0 {
		return nil, fmt.Errorf("no image specified, specify one via the IMAGE env var")
	}
	kubeClient, err := kubeclient.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}
	clusterAPIClient, err := clusterclientset.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}
	clusterAPIInformerFactory := clusterinformers.NewSharedInformerFactoryWithOptions(
		clusterAPIClient,
		defaultResyncInterval,
		clusterinformers.WithNamespace(clusterAPINamespace))
	machineInformer := clusterAPIInformerFactory.Cluster().V1alpha1().Machines().Informer()

	eventRecorder := setupEventRecorder(kubeClient)

	return &ReconcileEtcdConfig{
		Client:           mgr.GetClient(),
		scheme:           mgr.GetScheme(),
		externalDNSImage: externalDNSImage,
		kubeClient:       kubeClient,
		clusterAPIClient: clusterAPIClient,
		eventRecorder:    eventRecorder,
		machineInformer:  machineInformer,
	}, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("etcdconfig-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	if err = ensureOperatorConfig(mgr); err != nil {
		return err
	}
	// Watch for changes to EtcdConfig
	err = c.Watch(&source.Kind{Type: &etcdv1.EtcdConfig{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

func mustAsset(name string) []byte {
	b, err := assets.Asset(name)
	if err != nil {
		panic(fmt.Sprintf("Asset not found: %s. Error: %v", name, err))
	}
	return b
}

func ensureOperatorConfig(mgr manager.Manager) error {
	dynamicClient, err := dynamic.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}
	v1helpers.EnsureOperatorConfigExists(
		dynamicClient,
		mustAsset("config/operator/config.yaml"),
		schema.GroupVersionResource{Group: etcdv1.SchemeGroupVersion.Group, Version: "v1", Resource: "etcdconfigs"},
	)
	return nil
}

var _ reconcile.Reconciler = &ReconcileEtcdConfig{}

// ReconcileEtcdConfig reconciles a EtcdConfig object
type ReconcileEtcdConfig struct {
	client.Client
	scheme           *runtime.Scheme
	externalDNSImage string
	kubeClient       kubeclient.Interface
	clusterAPIClient clusterclientset.Interface
	eventRecorder    events.Recorder
}

// Reconcile reads that state of the cluster for a EtcdConfig object and makes changes based on the state read
// and what is in the EtcdConfig.Spec
// +kubebuilder:rbac:groups=etcd.operator.openshift.io,resources=etcdconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=etcd.operator.openshift.io,resources=etcdconfigs/status,verbs=get;update;patch
func (r *ReconcileEtcdConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the EtcdConfig instance
	instance := &etcdv1.EtcdConfig{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
