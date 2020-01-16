module github.com/nmagnezi/observatorium-operator

go 1.13

require (
	github.com/coreos/prometheus-operator v0.35.0
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/openshift/api v3.9.1-0.20190809235250-af7bae2945fe+incompatible
	github.com/openshift/client-go v0.0.0-20191022152013-2823239d2298
	github.com/pkg/errors v0.8.1
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	k8s.io/api v0.17.1
	k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783
	k8s.io/apimachinery v0.17.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.17.1
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.8.2-0.20190819201610-48b2c9c8eae2 // v1.8.2 is misleading as Prometheus does not have v2 module. This is pointing to one commit after 2.12.0.
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
)
