package assets

import (
	"fmt"
	"strings"
)

var _config_crds_etcd_v1_dnsendpoint_yaml = []byte(`apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: dnsendpoints.etcd.operator.openshift.io
spec:
  group: etcd.operator.openshift.io
  names:
    kind: DNSEndpoint
    plural: dnsendpoints
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            endpoints:
              description: Endpoints is the list of DNS records to create/update
              items:
                properties:
                  dnsName:
                    description: The hostname of the DNS record
                    type: string
                  labels:
                    description: Labels stores labels defined for the Endpoint
                    type: object
                  providerSpecific:
                    description: ProviderSpecific stores provider specific config
                    type: object
                  recordTTL:
                    description: TTL for the record
                    format: int64
                    type: integer
                  recordType:
                    description: RecordType type of record, e.g. CNAME, A, SRV, TXT
                      etc
                    type: string
                  targets:
                    description: The targets the DNS record points to
                    items:
                      type: string
                    type: array
                type: object
              type: array
          type: object
        status:
          properties:
            observedGeneration:
              description: ObservedGeneration is the generation observed by the external-dns
                controller.
              format: int64
              type: integer
          type: object
  version: v1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func config_crds_etcd_v1_dnsendpoint_yaml() ([]byte, error) {
	return _config_crds_etcd_v1_dnsendpoint_yaml, nil
}

var _config_crds_etcd_v1_etcdconfig_yaml = []byte(`apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: etcdconfigs.etcd.operator.openshift.io
spec:
  group: etcd.operator.openshift.io
  names:
    kind: EtcdConfig
    plural: etcdconfigs
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            dnsConfig:
              description: dnsConfig specifies configuration for etcd DNS entry management
              properties:
                automaticUpdates:
                  description: automaticUpdates indicates that the DNS entries should
                    be automatically updated based on the IP address of master Machine
                    resources
                  type: boolean
                logLevel:
                  description: 'logLevel is the level of logging for the external-dns
                    controller Valid values: debug, info, warning, error, fatal'
                  type: string
              required:
              - automaticUpdates
              type: object
            managementState:
              description: managementState indicates whether and how the operator
                should manage the component
              type: string
          required:
          - managementState
          - dnsConfig
          type: object
        status:
          properties:
            conditions:
              description: conditions is a list of conditions and their status
              items:
                properties:
                  lastTransitionTime:
                    format: date-time
                    type: string
                  message:
                    type: string
                  reason:
                    type: string
                  status:
                    type: string
                  type:
                    type: string
                required:
                - type
                - status
                type: object
              type: array
            generations:
              description: generations are used to determine when an item needs to
                be reconciled or has changed in a way that needs a reaction.
              items:
                properties:
                  group:
                    description: group is the group of the thing you're tracking
                    type: string
                  hash:
                    description: hash is an optional field set for resources without
                      generation that are content sensitive like secrets and configmaps
                    type: string
                  lastGeneration:
                    description: lastGeneration is the last generation of the workload
                      controller involved
                    format: int64
                    type: integer
                  name:
                    description: name is the name of the thing you're tracking
                    type: string
                  namespace:
                    description: namespace is where the thing you're tracking is
                    type: string
                  resource:
                    description: resource is the resource type of the thing you're
                      tracking
                    type: string
                required:
                - group
                - resource
                - namespace
                - name
                - lastGeneration
                - hash
                type: object
              type: array
            observedGeneration:
              description: observedGeneration is the last generation change you've
                dealt with
              format: int64
              type: integer
            readyReplicas:
              description: readyReplicas indicates how many replicas are ready and
                at the desired state
              format: int32
              type: integer
            version:
              description: version is the level this availability applies to
              type: string
          required:
          - version
          - readyReplicas
          - generations
          type: object
  version: v1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func config_crds_etcd_v1_etcdconfig_yaml() ([]byte, error) {
	return _config_crds_etcd_v1_etcdconfig_yaml, nil
}

var _config_default_kustomization_yaml = []byte(`# Adds namespace to all resources.
namespace: cluster-etcd-operator-system

# Value of this field is prepended to the
# names of all resources, e.g. a deployment named
# "wordpress" becomes "alices-wordpress".
# Note that it should also match with the prefix (text before '-') of the namespace
# field above.
namePrefix: cluster-etcd-operator-

# Labels to add to all resources and selectors.
#commonLabels:
#  someName: someValue

# Each entry in this list must resolve to an existing
# resource definition in YAML.  These are the resource
# files that kustomize reads, modifies and emits as a
# YAML string, with resources separated by document
# markers ("---").
resources:
- ../rbac/rbac_role.yaml
- ../rbac/rbac_role_binding.yaml
- ../manager/manager.yaml
  # Comment the following 3 lines if you want to disable
  # the auth proxy (https://github.com/brancz/kube-rbac-proxy)
  # which protects your /metrics endpoint.
- ../rbac/auth_proxy_service.yaml
- ../rbac/auth_proxy_role.yaml
- ../rbac/auth_proxy_role_binding.yaml

patches:
- manager_image_patch.yaml
  # Protect the /metrics endpoint by putting it behind auth.
  # Only one of manager_auth_proxy_patch.yaml and
  # manager_prometheus_metrics_patch.yaml should be enabled.
- manager_auth_proxy_patch.yaml
  # If you want your controller-manager to expose the /metrics
  # endpoint w/o any authn/z, uncomment the following line and
  # comment manager_auth_proxy_patch.yaml.
  # Only one of manager_auth_proxy_patch.yaml and
  # manager_prometheus_metrics_patch.yaml should be enabled.
#- manager_prometheus_metrics_patch.yaml

vars:
- name: WEBHOOK_SECRET_NAME
  objref:
    kind: Secret
    name: webhook-server-secret
    apiVersion: v1
`)

func config_default_kustomization_yaml() ([]byte, error) {
	return _config_default_kustomization_yaml, nil
}

var _config_default_manager_auth_proxy_patch_yaml = []byte(`# This patch inject a sidecar container which is a HTTP proxy for the controller manager,
# it performs RBAC authorization against the Kubernetes API using SubjectAccessReviews.
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: controller-manager
  namespace: system
spec:
  template:
    spec:
      containers:
      - name: kube-rbac-proxy
        image: gcr.io/kubebuilder/kube-rbac-proxy:v0.4.0
        args:
        - "--secure-listen-address=0.0.0.0:8443"
        - "--upstream=http://127.0.0.1:8080/"
        - "--logtostderr=true"
        - "--v=10"
        ports:
        - containerPort: 8443
          name: https
      - name: manager
        args:
        - "--metrics-addr=127.0.0.1:8080"
`)

func config_default_manager_auth_proxy_patch_yaml() ([]byte, error) {
	return _config_default_manager_auth_proxy_patch_yaml, nil
}

var _config_default_manager_image_patch_yaml = []byte(`apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: controller-manager
  namespace: system
spec:
  template:
    spec:
      containers:
      # Change the value of image field below to your controller image URL
      - image: IMAGE_URL
        name: manager
`)

func config_default_manager_image_patch_yaml() ([]byte, error) {
	return _config_default_manager_image_patch_yaml, nil
}

var _config_default_manager_prometheus_metrics_patch_yaml = []byte(`# This patch enables Prometheus scraping for the manager pod.
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: controller-manager
  namespace: system
spec:
  template:
    metadata:
      annotations:
        prometheus.io/scrape: 'true'
    spec:
      containers:
      # Expose the prometheus metrics on default port
      - name: manager
        ports:
        - containerPort: 8080
          name: metrics
          protocol: TCP
`)

func config_default_manager_prometheus_metrics_patch_yaml() ([]byte, error) {
	return _config_default_manager_prometheus_metrics_patch_yaml, nil
}

var _config_external_dns_binding_yaml = []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  creationTimestamp: null
  name: external-dns
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: external-dns
subjects:
- kind: ServiceAccount
  name: external-dns
  namespace: openshift-etcd-dns
`)

func config_external_dns_binding_yaml() ([]byte, error) {
	return _config_external_dns_binding_yaml, nil
}

var _config_external_dns_deployment_yaml = []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: external-dns
  namespace: openshift-etcd-dns
spec:
  replicas: 1
  selector:
    matchLabels:
      name: external-dns
  template:
    metadata:
      labels:
        name: external-dns
    spec:
      strategy:
        type: Recreate
      serviceAccountName: external-dns
      priorityClassName: system-cluster-critical
      containers:
      - name: external-dns
        image: quay.io/openshift/external-dns:latest
        args:
        - --source=crd
        - --crd-source-apiversion=etcd.operator.openshift.io/v1
        - --crd-source-kind=DNSEndpoint
        - --registry=noop
        - --policy=upsert-only
        - --namespace=openshift-etcd-dns

`)

func config_external_dns_deployment_yaml() ([]byte, error) {
	return _config_external_dns_deployment_yaml, nil
}

var _config_external_dns_role_yaml = []byte(``)

func config_external_dns_role_yaml() ([]byte, error) {
	return _config_external_dns_role_yaml, nil
}

var _config_external_dns_sa_yaml = []byte(``)

func config_external_dns_sa_yaml() ([]byte, error) {
	return _config_external_dns_sa_yaml, nil
}

var _config_manager_manager_yaml = []byte(`apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  name: system
---
apiVersion: v1
kind: Service
metadata:
  name: controller-manager-service
  namespace: system
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
spec:
  selector:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  ports:
  - port: 443
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: controller-manager
  namespace: system
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
spec:
  selector:
    matchLabels:
      control-plane: controller-manager
      controller-tools.k8s.io: "1.0"
  serviceName: controller-manager-service
  template:
    metadata:
      labels:
        control-plane: controller-manager
        controller-tools.k8s.io: "1.0"
    spec:
      containers:
      - command:
        - /manager
        image: controller:latest
        imagePullPolicy: Always
        name: manager
        env:
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: SECRET_NAME
            value: $(WEBHOOK_SECRET_NAME)
        resources:
          limits:
            cpu: 100m
            memory: 30Mi
          requests:
            cpu: 100m
            memory: 20Mi
        ports:
        - containerPort: 9876
          name: webhook-server
          protocol: TCP
        volumeMounts:
        - mountPath: /tmp/cert
          name: cert
          readOnly: true
      terminationGracePeriodSeconds: 10
      volumes:
      - name: cert
        secret:
          defaultMode: 420
          secretName: webhook-server-secret
---
apiVersion: v1
kind: Secret
metadata:
  name: webhook-server-secret
  namespace: system
`)

func config_manager_manager_yaml() ([]byte, error) {
	return _config_manager_manager_yaml, nil
}

var _config_operator_config_yaml = []byte(`apiVersion: etcd.operator.openshift.io/v1alpha1
kind: EtcdConfig
metadata:
  name: instance
spec:
  managementState: Managed
  dnsConfig:
    automaticUpdates: true
    logLevel: info
`)

func config_operator_config_yaml() ([]byte, error) {
	return _config_operator_config_yaml, nil
}

var _config_rbac_auth_proxy_role_yaml = []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: proxy-role
rules:
- apiGroups: ["authentication.k8s.io"]
  resources:
  - tokenreviews
  verbs: ["create"]
- apiGroups: ["authorization.k8s.io"]
  resources:
  - subjectaccessreviews
  verbs: ["create"]
`)

func config_rbac_auth_proxy_role_yaml() ([]byte, error) {
	return _config_rbac_auth_proxy_role_yaml, nil
}

var _config_rbac_auth_proxy_role_binding_yaml = []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: proxy-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: proxy-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: system
`)

func config_rbac_auth_proxy_role_binding_yaml() ([]byte, error) {
	return _config_rbac_auth_proxy_role_binding_yaml, nil
}

var _config_rbac_auth_proxy_service_yaml = []byte(`apiVersion: v1
kind: Service
metadata:
  annotations:
    prometheus.io/port: "8443"
    prometheus.io/scheme: https
    prometheus.io/scrape: "true"
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  name: controller-manager-metrics-service
  namespace: system
spec:
  ports:
  - name: https
    port: 8443
    targetPort: https
  selector:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
`)

func config_rbac_auth_proxy_service_yaml() ([]byte, error) {
	return _config_rbac_auth_proxy_service_yaml, nil
}

var _config_rbac_rbac_role_yaml = []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: manager-role
rules:
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - apps
  resources:
  - deployments/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - etcd.operator.openshift.io
  resources:
  - etcdconfigs
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - etcd.operator.openshift.io
  resources:
  - etcdconfigs/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - admissionregistration.k8s.io
  resources:
  - mutatingwebhookconfigurations
  - validatingwebhookconfigurations
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
`)

func config_rbac_rbac_role_yaml() ([]byte, error) {
	return _config_rbac_rbac_role_yaml, nil
}

var _config_rbac_rbac_role_binding_yaml = []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  creationTimestamp: null
  name: manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: manager-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: system
`)

func config_rbac_rbac_role_binding_yaml() ([]byte, error) {
	return _config_rbac_rbac_role_binding_yaml, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"config/crds/etcd_v1_dnsendpoint.yaml":                 config_crds_etcd_v1_dnsendpoint_yaml,
	"config/crds/etcd_v1_etcdconfig.yaml":                  config_crds_etcd_v1_etcdconfig_yaml,
	"config/default/kustomization.yaml":                    config_default_kustomization_yaml,
	"config/default/manager_auth_proxy_patch.yaml":         config_default_manager_auth_proxy_patch_yaml,
	"config/default/manager_image_patch.yaml":              config_default_manager_image_patch_yaml,
	"config/default/manager_prometheus_metrics_patch.yaml": config_default_manager_prometheus_metrics_patch_yaml,
	"config/external-dns/binding.yaml":                     config_external_dns_binding_yaml,
	"config/external-dns/deployment.yaml":                  config_external_dns_deployment_yaml,
	"config/external-dns/role.yaml":                        config_external_dns_role_yaml,
	"config/external-dns/sa.yaml":                          config_external_dns_sa_yaml,
	"config/manager/manager.yaml":                          config_manager_manager_yaml,
	"config/operator/config.yaml":                          config_operator_config_yaml,
	"config/rbac/auth_proxy_role.yaml":                     config_rbac_auth_proxy_role_yaml,
	"config/rbac/auth_proxy_role_binding.yaml":             config_rbac_auth_proxy_role_binding_yaml,
	"config/rbac/auth_proxy_service.yaml":                  config_rbac_auth_proxy_service_yaml,
	"config/rbac/rbac_role.yaml":                           config_rbac_rbac_role_yaml,
	"config/rbac/rbac_role_binding.yaml":                   config_rbac_rbac_role_binding_yaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"config": {nil, map[string]*_bintree_t{
		"crds": {nil, map[string]*_bintree_t{
			"etcd_v1_dnsendpoint.yaml": {config_crds_etcd_v1_dnsendpoint_yaml, map[string]*_bintree_t{}},
			"etcd_v1_etcdconfig.yaml":  {config_crds_etcd_v1_etcdconfig_yaml, map[string]*_bintree_t{}},
		}},
		"default": {nil, map[string]*_bintree_t{
			"kustomization.yaml":                    {config_default_kustomization_yaml, map[string]*_bintree_t{}},
			"manager_auth_proxy_patch.yaml":         {config_default_manager_auth_proxy_patch_yaml, map[string]*_bintree_t{}},
			"manager_image_patch.yaml":              {config_default_manager_image_patch_yaml, map[string]*_bintree_t{}},
			"manager_prometheus_metrics_patch.yaml": {config_default_manager_prometheus_metrics_patch_yaml, map[string]*_bintree_t{}},
		}},
		"external-dns": {nil, map[string]*_bintree_t{
			"binding.yaml":    {config_external_dns_binding_yaml, map[string]*_bintree_t{}},
			"deployment.yaml": {config_external_dns_deployment_yaml, map[string]*_bintree_t{}},
			"role.yaml":       {config_external_dns_role_yaml, map[string]*_bintree_t{}},
			"sa.yaml":         {config_external_dns_sa_yaml, map[string]*_bintree_t{}},
		}},
		"manager": {nil, map[string]*_bintree_t{
			"manager.yaml": {config_manager_manager_yaml, map[string]*_bintree_t{}},
		}},
		"operator": {nil, map[string]*_bintree_t{
			"config.yaml": {config_operator_config_yaml, map[string]*_bintree_t{}},
		}},
		"rbac": {nil, map[string]*_bintree_t{
			"auth_proxy_role.yaml":         {config_rbac_auth_proxy_role_yaml, map[string]*_bintree_t{}},
			"auth_proxy_role_binding.yaml": {config_rbac_auth_proxy_role_binding_yaml, map[string]*_bintree_t{}},
			"auth_proxy_service.yaml":      {config_rbac_auth_proxy_service_yaml, map[string]*_bintree_t{}},
			"rbac_role.yaml":               {config_rbac_rbac_role_yaml, map[string]*_bintree_t{}},
			"rbac_role_binding.yaml":       {config_rbac_rbac_role_binding_yaml, map[string]*_bintree_t{}},
		}},
	}},
}}
