
# Image URL to use all building/pushing image targets
IMG ?= controller:latest

all: fmt vet lint build

# Run tests
test-unit:
	go test ./pkg/... ./cmd/... -coverprofile cover.out

verify: verify-bindata verify-gofmt vet

# Build operator binary
build: generate
	go build -o bin/operator github.com/openshift/cluster-etcd-operator/cmd/operator

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crds
	kustomize build config/default | kubectl apply -f -

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Install manifests on cluster
install:
	kubectl apply -f ./config/crds

# Generate code
generate: generate-bindata generate-crds
	go generate ./pkg/... ./cmd/...

verify-gofmt:
	hack/verify-gofmt.sh

verify-bindata:
	hack/verify-bindata.sh

generate-bindata:
	hack/update-bindata.sh

generate-crds:
	go run ./vendor/sigs.k8s.io/controller-tools/cmd/controller-gen crd

# Build the docker image
docker-build:
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

# Utility target to run operator
run:
	WATCH_NAMESPACE=openshift-cluster-etcd-dns IMAGE=quay.io/csrwng/external-dns:latest ./bin/operator
