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

package externaldns

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/dynamic"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterclientset "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	clusterinformers "sigs.k8s.io/cluster-api/pkg/client/informers_generated/externalversions"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	configclientset "github.com/openshift/client-go/config/clientset/versioned"
	installertypes "github.com/openshift/installer/pkg/types"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"

	etcdv1 "github.com/openshift/cluster-etcd-operator/pkg/apis/etcd/v1"
	"github.com/openshift/cluster-etcd-operator/pkg/controller/assets"
)

const (
	operatorNamespace        = "cluster-etcd-operator"
	clusterAPINamespace      = "openshift-cluster-api"
	clusterConfigNamespace   = "kube-system"
	clusterConfigName        = "cluster-config-v1"
	installConfigKey         = "install-config"
	workloadFailingCondition = "WorkloadFailing"
)

var (
	defaultResyncInterval = 10 * time.Minute
	configName            = types.NamespacedName{Name: "instance"}
)

// Add creates a new Etcd Operator reconciler and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and start it when the manager is started.
func Add(mgr manager.Manager) error {
	reconciler, err := newReconciler(mgr)
	if err != nil {
		return err
	}
	return add(mgr, reconciler)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) (*ReconcileEtcdConfig, error) {
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
	configClient, err := configclientset.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}
	apiExtensionsClient, err := apiextclient.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	clusterAPIInformerFactory := clusterinformers.NewSharedInformerFactoryWithOptions(
		clusterAPIClient,
		defaultResyncInterval,
		clusterinformers.WithNamespace(clusterAPINamespace))

	eventRecorder := setupEventRecorder(kubeClient)

	machineInformer := clusterAPIInformerFactory.Cluster().V1alpha1().Machines().Informer()

	targetNamespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		return nil, err
	}

	// Retrieve cluster install config
	installConfig, err := getInstallConfig(kubeClient)
	if err != nil {
		log.WithError(err).Error("failed to get installconfig")
		return nil, err
	}

	// Retrieve DNS config
	dnsConfig := &configv1.DNS{}
	// TODO: Use this when client issue is resolved.
	/*
		err = mgr.GetClient().Get(context.TODO(), types.NamespacedName{Name: "cluster"}, dnsConfig)
		if err != nil {
			log.WithError(err).Error("failed to get dns config")
			return nil, err
		}
	*/

	// Retrieve the typed cluster version config.
	clusterVersionConfig, err := configClient.ConfigV1().ClusterVersions().Get("version", metav1.GetOptions{})
	if err != nil {
		log.WithError(err).Error("failed to get clusterversion 'version'")
		return nil, err
	}

	// Build a provider based on existing config
	provider, err := getProvider(kubeClient, targetNamespace, installConfig, dnsConfig, clusterVersionConfig)
	if err != nil {
		log.WithError(err).Error("failed to setup provider")
		return nil, err
	}

	return &ReconcileEtcdConfig{
		Client:              mgr.GetClient(),
		scheme:              mgr.GetScheme(),
		externalDNSImage:    externalDNSImage,
		kubeClient:          kubeClient,
		clusterAPIClient:    clusterAPIClient,
		apiExtensionsClient: apiExtensionsClient,
		eventRecorder:       eventRecorder,
		machineInformer:     machineInformer,
		targetNamespace:     targetNamespace,
		provider:            provider,
	}, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r *ReconcileEtcdConfig) error {
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

	handlerFunc := handler.ToRequestsFunc(func(mapObject handler.MapObject) []reconcile.Request {
		return []reconcile.Request{{NamespacedName: configName}}
	})

	// Watch for changes to machines in cluster API namespace
	err = c.Watch(&source.Informer{Informer: r.machineInformer}, &handler.EnqueueRequestsFromMapFunc{ToRequests: handlerFunc})
	if err != nil {
		return err
	}

	// Watch for Deployments
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: handlerFunc})
	if err != nil {
		return err
	}

	// Watch for Service Accounts
	err = c.Watch(&source.Kind{Type: &corev1.ServiceAccount{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: handlerFunc})
	if err != nil {
		return err
	}
	return nil
}

func ensureOperatorConfig(mgr manager.Manager) error {
	dynamicClient, err := dynamic.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}
	v1helpers.EnsureOperatorConfigExists(
		dynamicClient,
		assets.MustAsset("config/operator/config.yaml"),
		schema.GroupVersionResource{Group: etcdv1.SchemeGroupVersion.Group, Version: "v1", Resource: "etcdconfigs"},
	)
	return nil
}

var _ reconcile.Reconciler = &ReconcileEtcdConfig{}

// ReconcileEtcdConfig reconciles a EtcdConfig object
type ReconcileEtcdConfig struct {
	client.Client
	scheme              *runtime.Scheme
	externalDNSImage    string
	kubeClient          kubeclient.Interface
	clusterAPIClient    clusterclientset.Interface
	apiExtensionsClient apiextclient.Interface
	eventRecorder       events.Recorder
	machineInformer     cache.SharedIndexInformer
	targetNamespace     string
	provider            provider
}

// Reconcile reads that state of the cluster for a EtcdConfig object and makes changes based on the state read
// and what is in the EtcdConfig.Spec
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

	switch instance.Spec.ManagementState {
	case operatorv1.Unmanaged:
		return reconcile.Result{}, nil

	case operatorv1.Removed:
		err = r.removeEtcdDNS()
		return reconcile.Result{}, err
	}

	if err = r.syncEtcdDNS(instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileEtcdConfig) removeEtcdDNS() error {
	if err := r.kubeClient.CoreV1().Namespaces().Delete(r.targetNamespace, nil); err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (r *ReconcileEtcdConfig) syncEtcdDNS(config *etcdv1.EtcdConfig) error {

	originalConfig := config.DeepCopy()

	// Ensure that the CRD is defined
	crd := resourceread.ReadCustomResourceDefinitionV1Beta1OrDie(assets.MustAsset("config/crds/etcd_v1_dnsendpoint.yaml"))
	_, _, err := resourceapply.ApplyCustomResourceDefinition(r.apiExtensionsClient.ApiextensionsV1beta1(), r.eventRecorder, crd)
	if err != nil {
		return err
	}

	// Apply the rest of the external-dns application artifacts
	results := resourceapply.ApplyDirectly(r.kubeClient, r.eventRecorder, assets.Asset,
		"config/external-dns/namespace.yaml",
		"config/external-dns/role.yaml",
		"config/external-dns/binding.yaml",
		"config/external-dns/sa.yaml",
	)
	resourcesThatForceRedeployment := sets.NewString("config/external-dns/sa.yaml")
	forceDeployment := config.Generation != config.Status.ObservedGeneration

	errs := []error{}

	for _, currResult := range results {
		if currResult.Error != nil {
			log.WithError(currResult.Error).WithField("file", currResult.File).Error("Apply error")
			errs = append(errs, fmt.Errorf("%q (%T): %v", currResult.File, currResult.Type, currResult.Error))
			continue
		}

		if currResult.Changed && resourcesThatForceRedeployment.Has(currResult.File) {
			forceDeployment = true
		}
	}

	// Ensure that the AWS credentials secret exists
	// TODO: Remove when cred minter provides the secret
	creds, err := r.provider.Creds()
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to obtain creds: %v", err))
	} else {
		_, changedSecret, err := resourceapply.ApplySecret(r.kubeClient.CoreV1(), r.eventRecorder, creds)
		if err != nil {
			errs = append(errs, fmt.Errorf("%q: %v", "config/external-dns/secret.yaml", err))
		}
		if changedSecret {
			forceDeployment = true
		}
	}

	actualDeployment, _, err := r.syncExternalDNSDeployment(config, forceDeployment)
	if err != nil {
		log.WithError(err).WithField("file", "config/external-dns/deployment.yaml").Error("Apply error")
		errs = append(errs, fmt.Errorf("%q: %v", "config/external-dns/deployment.yaml", err))
	}

	if actualDeployment.Status.ReadyReplicas > 0 {
		v1helpers.SetOperatorCondition(&config.Status.Conditions, operatorv1.OperatorCondition{
			Type:               operatorv1.OperatorStatusTypeAvailable,
			Status:             operatorv1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
		})
	} else {
		v1helpers.SetOperatorCondition(&config.Status.Conditions, operatorv1.OperatorCondition{
			Type:               operatorv1.OperatorStatusTypeAvailable,
			Status:             operatorv1.ConditionFalse,
			Reason:             "NoPodsAvailable",
			Message:            "no deployment pods available on any node.",
			LastTransitionTime: metav1.Now(),
		})
	}

	config.Status.ObservedGeneration = config.ObjectMeta.Generation
	config.Status.ReadyReplicas = actualDeployment.Status.ReadyReplicas
	config.Status.Version = "1.0.0"
	resourcemerge.SetDeploymentGeneration(&config.Status.Generations, actualDeployment)
	if len(errs) > 0 {
		message := ""
		for _, err := range errs {
			message = message + err.Error() + "\n"
		}
		v1helpers.SetOperatorCondition(&config.Status.Conditions, operatorv1.OperatorCondition{
			Type:               workloadFailingCondition,
			Status:             operatorv1.ConditionTrue,
			Message:            message,
			Reason:             "SyncError",
			LastTransitionTime: metav1.Now(),
		})
	} else {
		v1helpers.SetOperatorCondition(&config.Status.Conditions, operatorv1.OperatorCondition{
			Type:               workloadFailingCondition,
			Status:             operatorv1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
		})
	}
	if !equality.Semantic.DeepEqual(config.Status, originalConfig.Status) {
		if err := r.Status().Update(context.TODO(), config); err != nil {
			log.WithError(err).Error("Status update error")
			return err
		}
	}

	if len(errs) > 0 {
		log.WithField("Errors", errs).Error("errors occurred, requeuing")
		return utilerrors.NewAggregate(errs)
	}
	return nil
}

func (r *ReconcileEtcdConfig) syncExternalDNSDeployment(config *etcdv1.EtcdConfig, forceDeployment bool) (*appsv1.Deployment, bool, error) {
	deployment := resourceread.ReadDeploymentV1OrDie(assets.MustAsset("config/external-dns/deployment.yaml"))
	container := &deployment.Spec.Template.Spec.Containers[0]
	container.Image = r.externalDNSImage
	container.ImagePullPolicy = corev1.PullAlways
	for _, arg := range r.provider.Args() {
		container.Args = append(container.Args, arg)
	}
	if len(config.Spec.DNSConfig.LogLevel) > 0 {
		container.Args = append(container.Args, fmt.Sprintf("--log-level=%s", config.Spec.DNSConfig.LogLevel))
	}
	container.Env = append(container.Env, r.provider.Env()...)
	return resourceapply.ApplyDeployment(
		r.kubeClient.AppsV1(),
		r.eventRecorder,
		deployment,
		resourcemerge.ExpectedDeploymentGeneration(deployment, config.Status.Generations),
		forceDeployment)
}

func getInstallConfig(client kubeclient.Interface) (*installertypes.InstallConfig, error) {
	cm, err := client.CoreV1().ConfigMaps(clusterConfigNamespace).Get(clusterConfigName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed getting clusterconfig %s/%s: %v", clusterConfigNamespace, clusterConfigName, err)
	}
	installConfigData, ok := cm.Data[installConfigKey]
	if !ok {
		return nil, fmt.Errorf("missing %q in configmap", installConfigKey)
	}
	installConfig := &installertypes.InstallConfig{}
	if err := yaml.Unmarshal([]byte(installConfigData), installConfig); err != nil {
		return nil, fmt.Errorf("invalid InstallConfig: %v yaml: %s", err, installConfigData)
	}
	return installConfig, nil
}

func getProvider(client kubeclient.Interface, namespace string, installConfig *installertypes.InstallConfig, dnsConfig *configv1.DNS, clusterVersionConfig *configv1.ClusterVersion) (provider, error) {
	switch {
	case installConfig.Platform.AWS != nil:
		return newAWSProvider(client, namespace, installConfig, dnsConfig, clusterVersionConfig), nil
	}
	return nil, fmt.Errorf("Platform is unsupported")
}

type provider interface {
	Env() []corev1.EnvVar
	Args() []string
	Creds() (*corev1.Secret, error)
}
