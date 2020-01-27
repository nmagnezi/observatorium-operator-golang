
# Image URL to use all building/pushing image targets
IMG ?= observatorium-operator:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

FIRST_GOPATH:=$(firstword $(subst :, ,$(shell go env GOPATH)))
GOBINDATA_BIN=$(FIRST_GOPATH)/bin/go-bindata

# Copy the logic to get kustomize from performance-addon-operators
CACHE_DIR="_cache"
TOOLS_DIR="$(CACHE_DIR)/tools"
KUSTOMIZE_VERSION="v3.5.3"
KUSTOMIZE_PLATFORM ?= "linux_amd64"
KUSTOMIZE_BIN="kustomize"
KUSTOMIZE_TAR="$(KUSTOMIZE_BIN)_$(KUSTOMIZE_VERSION)_$(KUSTOMIZE_PLATFORM).tar.gz"
KUSTOMIZE="$(TOOLS_DIR)/$(KUSTOMIZE_BIN)"

CONTAINER_CLIENT=podman

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests kustomize
	cd config/manager && ../../$(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.4 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

manifests/bindata.go: $(GOBINDATA_BIN)
	# Using "-modtime 1" to make generate target deterministic. It sets all file time stamps to unix timestamp 1
	go-bindata -mode 420 -modtime 1 -pkg manifests -o $@ assets/...

kustomize:
	@if [ ! -x "$(KUSTOMIZE)" ]; then\
		echo "Downloading kustomize $(KUSTOMIZE_VERSION)";\
		mkdir -p $(TOOLS_DIR);\
		curl -JL https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/$(KUSTOMIZE_VERSION)/$(KUSTOMIZE_TAR) -o $(TOOLS_DIR)/$(KUSTOMIZE_TAR);\
		tar -xvf $(TOOLS_DIR)/$(KUSTOMIZE_TAR) -C $(TOOLS_DIR);\
		rm -rf $(TOOLS_DIR)/$(KUSTOMIZE_TAR);\
		chmod +x $(KUSTOMIZE);\
	else\
		echo "Using kustomize cached at $(KUSTOMIZE)";\
	fi

container-build-operator-courier:
	$(CONTAINER_CLIENT) build -f tools/operator-courier/Dockerfile -t courier-build-container .

bundle-push: container-build-operator-courier
	@QUAY_USERNAME=$(QUAY_USERNAME) QUAY_PASSWORD=$(QUAY_PASSWORD) ./tools/operator-courier/push.sh


.PHONY: bundle-push
