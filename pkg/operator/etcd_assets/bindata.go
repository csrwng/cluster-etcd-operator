// Code generated by go-bindata.
// sources:
// bindata/etcd/cm.yaml
// bindata/etcd/defaultconfig.yaml
// bindata/etcd/ns.yaml
// bindata/etcd/pod-cm.yaml
// bindata/etcd/pod.yaml
// bindata/etcd/sa.yaml
// bindata/etcd/svc.yaml
// DO NOT EDIT!

package etcd_assets

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _etcdCmYaml = []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  namespace: openshift-etcd
  name: config
data:
  config.yaml:
`)

func etcdCmYamlBytes() ([]byte, error) {
	return _etcdCmYaml, nil
}

func etcdCmYaml() (*asset, error) {
	bytes, err := etcdCmYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "etcd/cm.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _etcdDefaultconfigYaml = []byte(`apiVersion: kubecontrolplane.config.openshift.io/v1
kind: EtcdConfig
`)

func etcdDefaultconfigYamlBytes() ([]byte, error) {
	return _etcdDefaultconfigYaml, nil
}

func etcdDefaultconfigYaml() (*asset, error) {
	bytes, err := etcdDefaultconfigYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "etcd/defaultconfig.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _etcdNsYaml = []byte(`apiVersion: v1
kind: Namespace
metadata:
  annotations:
    openshift.io/node-selector: ""
  name: openshift-etcd
  labels:
    openshift.io/run-level: "0"
`)

func etcdNsYamlBytes() ([]byte, error) {
	return _etcdNsYaml, nil
}

func etcdNsYaml() (*asset, error) {
	bytes, err := etcdNsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "etcd/ns.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _etcdPodCmYaml = []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  namespace: openshift-etcd
  name: etcd-pod
data:
  pod.yaml:
  forceRedeploymentReason:
  version:
`)

func etcdPodCmYamlBytes() ([]byte, error) {
	return _etcdPodCmYaml, nil
}

func etcdPodCmYaml() (*asset, error) {
	bytes, err := etcdPodCmYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "etcd/pod-cm.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _etcdPodYaml = []byte(`apiVersion: v1
kind: Pod
metadata:
  name: etcd
  namespace: openshift-etcd
  labels:
    app: etcd
    k8s-app: etcd
    etcd: "true"
    revision: "REVISION"
spec:
  initContainers:
    - name: etc-quorum-guard-copy
      image: ${IMAGE}
      imagePullPolicy: IfNotPresent
      terminationMessagePolicy: FallbackToLogsOnError
      command:
        - /bin/sh
        - -c
        - |
          #!/bin/sh
          set -euo pipefail

          cp /etc/kubernetes/static-pod-resources/secrets/etcd-all-peer/etcd-peer-NODE_NAME.crt /etc/kubernetes/etcd-backup-dir/system:etcd-peer-NODE_NAME.crt
          cp /etc/kubernetes/static-pod-resources/secrets/etcd-all-peer/etcd-peer-NODE_NAME.key /etc/kubernetes/etcd-backup-dir/system:etcd-peer-NODE_NAME.key
      resources:
        requests:
          memory: 60Mi
          cpu: 30m
      securityContext:
        privileged: true
      volumeMounts:
        - mountPath: /etc/kubernetes/etcd-backup-dir
          name: etcd-backup-dir
        - mountPath: /etc/kubernetes/static-pod-resources
          name: resource-dir
        - mountPath: /etc/kubernetes/static-pod-certs
          name: cert-dir
  containers:
  - name: etcd
    image: ${IMAGE}
    imagePullPolicy: IfNotPresent
    terminationMessagePolicy: FallbackToLogsOnError
    command:
      - /bin/sh
      - -c
      - |
        #!/bin/sh
        set -euo pipefail

        ETCDCTL="etcdctl --cacert=/etc/kubernetes/static-pod-resources/configmaps/etcd-serving-ca/ca-bundle.crt \
                           --cert=/etc/kubernetes/static-pod-resources/secrets/etcd-all-peer/etcd-peer-NODE_NAME.crt \
                           --key=/etc/kubernetes/static-pod-resources/secrets/etcd-all-peer/etcd-peer-NODE_NAME.key \
                           --endpoints=${ALL_ETCD_ENDPOINTS}"
        ${ETCDCTL} member list

        echo "waiting for member ${NODE_NODE_ENVVAR_NAME_ETCD_DNS_NAME}..."
        COUNT=30
        while [ $COUNT -gt 0 ]; do
          echo "current member list is..."
          ${ETCDCTL} member list
          echo ""
          echo ""

          IS_MEMBER_PRESENT=$(${ETCDCTL} member list | grep -o "${NODE_NODE_ENVVAR_NAME_ETCD_DNS_NAME}.*:2380" || true)
          if [[ -n "${IS_MEMBER_PRESENT:-}" ]]; then
            break
          fi
          sleep 1
          let COUNT=$COUNT-1
        done

        # if the member is not present after 30 seconds
        if [ -z "$IS_MEMBER_PRESENT" ]; then
          echo "member ${NODE_NODE_ENVVAR_NAME_ETCD_DNS_NAME} is not present after 30 seconds"
          exit 1
        fi
        echo "member ${NODE_NODE_ENVVAR_NAME_ETCD_DNS_NAME} is present, continuing"

        initial_cluster=""
        member_output=$( ${ETCDCTL} member list | cut -d',' -f3 )
        for endpoint_key in ${member_output}; do
          endpoint=$(${ETCDCTL} member list | grep $endpoint_key | awk -F'[, ]' '{ print $7 }')
          initial_cluster+="$endpoint_key=$endpoint,"
          echo "adding $endpoint_key=$endpoint,"
        done

        # if the member isn't started, then we need to add exactly what we expect to the initial cluster for this member
        echo "checking for unstarted"
        ${ETCDCTL} member list | grep "${NODE_NODE_ENVVAR_NAME_ETCD_DNS_NAME}" | grep unstarted || true
        IS_MEMBER_UNSTARTED=$(${ETCDCTL} member list | grep "${NODE_NODE_ENVVAR_NAME_ETCD_DNS_NAME}" | grep unstarted || true)
        if [[ -n "${IS_MEMBER_UNSTARTED:-}" ]]; then
          initial_cluster+="NODE_NAME=https://${NODE_NODE_ENVVAR_NAME_ETCD_DNS_NAME}:2380,"
          echo "adding unstarted NODE_NAME=https://${NODE_NODE_ENVVAR_NAME_ETCD_DNS_NAME}:2380,"
          break
        fi

        # trim last comma
        initial_cluster="${initial_cluster::-1}"
        echo $initial_cluster

        # at this point we know this member is added.  To support a transition, we must remove the old etcd pod.
        # move it somewhere safe so we can retrieve it again later if something goes badly.
        mv /etc/kubernetes/manifests/etcd-member.yaml /etc/kubernetes/etcd-backup-dir || true

        export ETCD_INITIAL_CLUSTER=${initial_cluster}
        export ETCD_NAME=${NODE_NODE_ENVVAR_NAME_ETCD_NAME}
        env | grep ETCD | grep -v NODE

        set -x
        exec etcd \
          --initial-advertise-peer-urls=https://${NODE_NODE_ENVVAR_NAME_IP}:2380 \
          --cert-file=/etc/kubernetes/static-pod-resources/secrets/etcd-all-serving/etcd-serving-NODE_NAME.crt \
          --key-file=/etc/kubernetes/static-pod-resources/secrets/etcd-all-serving/etcd-serving-NODE_NAME.key \
          --trusted-ca-file=/etc/kubernetes/static-pod-resources/configmaps/etcd-serving-ca/ca-bundle.crt \
          --client-cert-auth=true \
          --peer-cert-file=/etc/kubernetes/static-pod-resources/secrets/etcd-all-peer/etcd-peer-NODE_NAME.crt \
          --peer-key-file=/etc/kubernetes/static-pod-resources/secrets/etcd-all-peer/etcd-peer-NODE_NAME.key \
          --peer-trusted-ca-file=/etc/kubernetes/static-pod-resources/configmaps/etcd-peer-client-ca/ca-bundle.crt \
          --peer-client-cert-auth=true \
          --advertise-client-urls=https://${NODE_NODE_ENVVAR_NAME_IP}:2379 \
          --listen-client-urls=https://${LISTEN_ON_ALL_IPS}:2379 \
          --listen-peer-urls=https://${LISTEN_ON_ALL_IPS}:2380 \
          --listen-metrics-urls=https://${LISTEN_ON_ALL_IPS}:9978 ||  mv /etc/kubernetes/etcd-backup-dir/etcd-member.yaml /etc/kubernetes/manifests
    env:
${COMPUTED_ENV_VARS}
    resources:
      requests:
        memory: 600Mi
        cpu: 300m
    readinessProbe:
      exec:
        command:
          - /bin/sh
          - -ec
          - "lsof -n -i :2380 | grep LISTEN"
      failureThreshold: 3
      initialDelaySeconds: 3
      periodSeconds: 5
      successThreshold: 1
      timeoutSeconds: 5
    securityContext:
      privileged: true
    volumeMounts:
      - mountPath: /etc/kubernetes/manifests
        name: static-pod-dir
      - mountPath: /etc/kubernetes/etcd-backup-dir
        name: etcd-backup-dir
      - mountPath: /etc/kubernetes/static-pod-resources
        name: resource-dir
      - mountPath: /etc/kubernetes/static-pod-certs
        name: cert-dir
      - mountPath: /var/lib/etcd/
        name: data-dir
  - name: etcd-metrics
    image: ${IMAGE}
    imagePullPolicy: IfNotPresent
    terminationMessagePolicy: FallbackToLogsOnError
    command:
      - /bin/sh
      - -c
      - |
        #!/bin/sh
        set -euo pipefail

        export ETCD_NAME=${NODE_NODE_ENVVAR_NAME_ETCD_NAME}

        exec etcd grpc-proxy start \
          --endpoints https://${NODE_NODE_ENVVAR_NAME_ETCD_DNS_NAME}:9978 \
          --metrics-addr https://${LISTEN_ON_ALL_IPS}:9979 \
          --listen-addr ${LOCALHOST_IP}:9977 \
          --key /etc/kubernetes/static-pod-resources/secrets/etcd-all-peer/etcd-peer-NODE_NAME.key \
          --key-file /etc/kubernetes/static-pod-resources/secrets/etcd-all-serving-metrics/etcd-serving-metrics-NODE_NAME.key \
          --cert /etc/kubernetes/static-pod-resources/secrets/etcd-all-peer/etcd-peer-NODE_NAME.crt \
          --cert-file /etc/kubernetes/static-pod-resources/secrets/etcd-all-serving-metrics/etcd-serving-metrics-NODE_NAME.crt \
          --cacert /etc/kubernetes/static-pod-resources/configmaps/etcd-peer-client-ca/ca-bundle.crt \
          --trusted-ca-file /etc/kubernetes/static-pod-resources/configmaps/etcd-metrics-proxy-serving-ca/ca-bundle.crt
    env:
${COMPUTED_ENV_VARS}
    resources:
      requests:
        memory: 200Mi
        cpu: 100m
    securityContext:
      privileged: true
    volumeMounts:
      - mountPath: /etc/kubernetes/static-pod-resources
        name: resource-dir
      - mountPath: /etc/kubernetes/static-pod-certs
        name: cert-dir
      - mountPath: /var/lib/etcd/
        name: data-dir
  hostNetwork: true
  priorityClassName: system-node-critical
  tolerations:
  - operator: "Exists"
  volumes:
    - hostPath:
        path: /etc/kubernetes/manifests
      name: static-pod-dir
    - hostPath:
        path: /etc/kubernetes/static-pod-resources/etcd-member
      name: etcd-backup-dir
    - hostPath:
        path: /etc/kubernetes/static-pod-resources/etcd-pod-REVISION
      name: resource-dir
    - hostPath:
        path: /etc/kubernetes/static-pod-resources/etcd-certs
      name: cert-dir
    - hostPath:
        path: /var/lib/etcd
        type: ""
      name: data-dir

`)

func etcdPodYamlBytes() ([]byte, error) {
	return _etcdPodYaml, nil
}

func etcdPodYaml() (*asset, error) {
	bytes, err := etcdPodYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "etcd/pod.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _etcdSaYaml = []byte(`apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: openshift-etcd
  name: etcd-sa
`)

func etcdSaYamlBytes() ([]byte, error) {
	return _etcdSaYaml, nil
}

func etcdSaYaml() (*asset, error) {
	bytes, err := etcdSaYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "etcd/sa.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _etcdSvcYaml = []byte(`apiVersion: v1
kind: Service
metadata:
  namespace: openshift-etcd
  name: etcd
  annotations:
    service.alpha.openshift.io/serving-cert-secret-name: serving-cert
    prometheus.io/scrape: "true"
    prometheus.io/scheme: https
spec:
  selector:
    etcd: "true"
  ports:
  - name: https
    port: 443
    targetPort: 10257
`)

func etcdSvcYamlBytes() ([]byte, error) {
	return _etcdSvcYaml, nil
}

func etcdSvcYaml() (*asset, error) {
	bytes, err := etcdSvcYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "etcd/svc.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
var _bindata = map[string]func() (*asset, error){
	"etcd/cm.yaml":            etcdCmYaml,
	"etcd/defaultconfig.yaml": etcdDefaultconfigYaml,
	"etcd/ns.yaml":            etcdNsYaml,
	"etcd/pod-cm.yaml":        etcdPodCmYaml,
	"etcd/pod.yaml":           etcdPodYaml,
	"etcd/sa.yaml":            etcdSaYaml,
	"etcd/svc.yaml":           etcdSvcYaml,
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"etcd": {nil, map[string]*bintree{
		"cm.yaml":            {etcdCmYaml, map[string]*bintree{}},
		"defaultconfig.yaml": {etcdDefaultconfigYaml, map[string]*bintree{}},
		"ns.yaml":            {etcdNsYaml, map[string]*bintree{}},
		"pod-cm.yaml":        {etcdPodCmYaml, map[string]*bintree{}},
		"pod.yaml":           {etcdPodYaml, map[string]*bintree{}},
		"sa.yaml":            {etcdSaYaml, map[string]*bintree{}},
		"svc.yaml":           {etcdSvcYaml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
