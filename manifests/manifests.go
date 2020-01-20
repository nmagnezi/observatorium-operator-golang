package manifests

import (
	"bytes"

	// #nosec
	"fmt"
	"hash/fnv"
	"io"
	"strconv"

	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	observatoriumv1alpha1 "github.com/nmagnezi/observatorium-operator/api/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	securityv1 "github.com/openshift/api/security/v1"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
)

var (
	ThanosQuerierDeployment      = "assets/thanos-querier-deployment.yaml"
	ThanosQuerierService         = "assets/thanos-querier-service.yaml"
	ThanosQuerierCacheDeployment = "assets/thanos-querier-cache-deployment.yaml"
	ThanosQuerierCacheService    = "assets/thanos-querier-cache-service.yaml"
	ThanosQuerierCacheConfigMap  = "assets/thanos-querier-cache-configmap.yaml"
	ThanosCompactorStatefulSet   = "assets/thanos-compactor-statefulSet.yaml"
	ThanosCompactorService       = "assets/thanos-compactor-service.yaml"
	ThanosStoreStatefulSet       = "assets/thanos-store-statefulSet.yaml"
	ThanosStoreService           = "assets/thanos-store-service.yaml"
)

func MustAssetReader(asset string) io.Reader {
	return bytes.NewReader(MustAsset(asset))
}

type Factory struct {
	namespace, namespaceUserWorkload string
	crd                              observatoriumv1alpha1.Observatorium
}

func NewFactory(namespace, namespaceUserWorkload string, c observatoriumv1alpha1.Observatorium) *Factory {
	return &Factory{
		namespace:             namespace,
		namespaceUserWorkload: namespaceUserWorkload,
		crd:                   c,
	}
}

func (f *Factory) NewDeployment(manifest io.Reader) (*appsv1.Deployment, error) {
	d, err := NewDeployment(manifest)
	if err != nil {
		return nil, err
	}

	if d.GetNamespace() == "" {
		d.SetNamespace(f.namespace)
	}

	return d, nil
}

func (f *Factory) NewService(manifest io.Reader) (*v1.Service, error) {
	s, err := NewService(manifest)
	if err != nil {
		return nil, err
	}

	if s.GetNamespace() == "" {
		s.SetNamespace(f.namespace)
	}

	return s, nil
}

func (f *Factory) NewStatefulSet(manifest io.Reader) (*appsv1.StatefulSet, error) {
	s, err := NewStatefulSet(manifest)
	if err != nil {
		return nil, err
	}

	if s.GetNamespace() == "" {
		s.SetNamespace(f.namespace)
	}

	return s, nil
}

func (f *Factory) NewConfigMap(manifest io.Reader) (*v1.ConfigMap, error) {
	cm, err := NewConfigMap(manifest)
	if err != nil {
		return nil, err
	}

	if cm.GetNamespace() == "" {
		cm.SetNamespace(f.namespace)
	}

	return cm, nil
}

func (f *Factory) ThanosQuerierDeployment(grpcTLS *v1.Secret, enableUserWorkloadMonitoring bool, trustedCA *v1.ConfigMap) (*appsv1.Deployment, error) {
	d, err := f.NewDeployment(MustAssetReader(ThanosQuerierDeployment))
	if err != nil {
		return nil, err
	}

	d.Namespace = f.namespace
	d.Spec.Replicas = f.crd.Spec.Thanos.Querier.Replicas

	return d, nil
}

func (f *Factory) ThanosQuerierCacheDeployment() (*appsv1.Deployment, error) {
	d, err := f.NewDeployment(MustAssetReader(ThanosQuerierCacheDeployment))
	if err != nil {
		return nil, err
	}

	d.Namespace = f.namespace
	d.Spec.Replicas = f.crd.Spec.Thanos.QuerierCache.Replicas
	return d, nil
}

func (f *Factory) ThanosQuerierCacheService() (*v1.Service, error) {
	s, err := f.NewService(MustAssetReader(ThanosQuerierCacheService))
	if err != nil {
		return nil, err
	}

	s.Namespace = f.namespace

	return s, nil
}

func (f *Factory) ThanosQuerierCacheConfigMap() (*v1.ConfigMap, error) {
	c, err := f.NewConfigMap(MustAssetReader(ThanosQuerierCacheConfigMap))
	if err != nil {
		return nil, err
	}

	c.Namespace = f.namespace

	return c, nil
}

func (f *Factory) ThanosQuerierService() (*v1.Service, error) {
	s, err := f.NewService(MustAssetReader(ThanosQuerierService))
	if err != nil {
		return nil, err
	}

	s.Namespace = f.namespace

	return s, nil
}

func (f *Factory) ThanosCompactorService() (*v1.Service, error) {
	s, err := f.NewService(MustAssetReader(ThanosCompactorService))
	if err != nil {
		return nil, err
	}

	s.Namespace = f.namespace

	return s, nil
}

func (f *Factory) ThanosCompactorStatefulSet() (*appsv1.StatefulSet, error) {
	d, err := f.NewStatefulSet(MustAssetReader(ThanosCompactorStatefulSet))
	if err != nil {
		return nil, err
	}

	d.Namespace = f.namespace
	d.Spec.Replicas = f.crd.Spec.Thanos.Compactor.Replicas
	d.Spec.Template.Spec.Containers[0].Resources.Limits[v1.ResourceCPU] = f.crd.Spec.Thanos.Compactor.Resources.Limits[v1.ResourceCPU]
	d.Spec.Template.Spec.Containers[0].Resources.Limits[v1.ResourceMemory] = f.crd.Spec.Thanos.Compactor.Resources.Limits[v1.ResourceMemory]
	d.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceCPU] = f.crd.Spec.Thanos.Compactor.Resources.Requests[v1.ResourceCPU]
	d.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceMemory] = f.crd.Spec.Thanos.Compactor.Resources.Requests[v1.ResourceMemory]

	return d, nil
}

func (f *Factory) ThanosStoreService() (*v1.Service, error) {
	s, err := f.NewService(MustAssetReader(ThanosStoreService))
	if err != nil {
		return nil, err
	}

	s.Namespace = f.namespace

	return s, nil
}

func (f *Factory) ThanosStoreStatefulSet() (*appsv1.StatefulSet, error) {
	d, err := f.NewStatefulSet(MustAssetReader(ThanosStoreStatefulSet))
	if err != nil {
		return nil, err
	}

	d.Namespace = f.namespace
	d.Spec.Replicas = f.crd.Spec.Thanos.Store.Replicas
	d.Spec.Template.Spec.Containers[0].Resources.Limits[v1.ResourceCPU] = f.crd.Spec.Thanos.Store.Resources.Limits[v1.ResourceCPU]
	d.Spec.Template.Spec.Containers[0].Resources.Limits[v1.ResourceMemory] = f.crd.Spec.Thanos.Store.Resources.Limits[v1.ResourceMemory]
	d.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceCPU] = f.crd.Spec.Thanos.Store.Resources.Requests[v1.ResourceCPU]
	d.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceMemory] = f.crd.Spec.Thanos.Store.Resources.Requests[v1.ResourceMemory]
	d.Spec.VolumeClaimTemplates[0].Spec.StorageClassName = f.crd.Spec.Thanos.Store.StorageClass
	d.Spec.VolumeClaimTemplates[0].Spec.Resources.Requests[v1.ResourceStorage] = resource.MustParse(*f.crd.Spec.Thanos.Store.PVCSize)
	return d, nil
}

func NewDaemonSet(manifest io.Reader) (*appsv1.DaemonSet, error) {
	ds := appsv1.DaemonSet{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&ds)
	if err != nil {
		return nil, err
	}

	return &ds, nil
}

func NewService(manifest io.Reader) (*v1.Service, error) {
	s := v1.Service{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func NewEndpoints(manifest io.Reader) (*v1.Endpoints, error) {
	e := v1.Endpoints{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&e)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func NewRoute(manifest io.Reader) (*routev1.Route, error) {
	r := routev1.Route{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func NewSecret(manifest io.Reader) (*v1.Secret, error) {
	s := v1.Secret{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func NewClusterRoleBinding(manifest io.Reader) (*rbacv1.ClusterRoleBinding, error) {
	crb := rbacv1.ClusterRoleBinding{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&crb)
	if err != nil {
		return nil, err
	}

	return &crb, nil
}

func NewClusterRole(manifest io.Reader) (*rbacv1.ClusterRole, error) {
	cr := rbacv1.ClusterRole{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&cr)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}

func NewRoleBinding(manifest io.Reader) (*rbacv1.RoleBinding, error) {
	rb := rbacv1.RoleBinding{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&rb)
	if err != nil {
		return nil, err
	}

	return &rb, nil
}

func NewRole(manifest io.Reader) (*rbacv1.Role, error) {
	r := rbacv1.Role{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func NewRoleBindingList(manifest io.Reader) (*rbacv1.RoleBindingList, error) {
	rbl := rbacv1.RoleBindingList{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&rbl)
	if err != nil {
		return nil, err
	}

	return &rbl, nil
}

func NewRoleList(manifest io.Reader) (*rbacv1.RoleList, error) {
	rl := rbacv1.RoleList{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&rl)
	if err != nil {
		return nil, err
	}

	return &rl, nil
}

func NewConfigMap(manifest io.Reader) (*v1.ConfigMap, error) {
	cm := v1.ConfigMap{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&cm)
	if err != nil {
		return nil, err
	}

	return &cm, nil
}

func NewConfigMapList(manifest io.Reader) (*v1.ConfigMapList, error) {
	cml := v1.ConfigMapList{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&cml)
	if err != nil {
		return nil, err
	}

	return &cml, nil
}

func NewServiceAccount(manifest io.Reader) (*v1.ServiceAccount, error) {
	sa := v1.ServiceAccount{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&sa)
	if err != nil {
		return nil, err
	}

	return &sa, nil
}

func NewPrometheus(manifest io.Reader) (*monv1.Prometheus, error) {
	p := monv1.Prometheus{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func NewPrometheusRule(manifest io.Reader) (*monv1.PrometheusRule, error) {
	p := monv1.PrometheusRule{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func NewAlertmanager(manifest io.Reader) (*monv1.Alertmanager, error) {
	a := monv1.Alertmanager{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func NewServiceMonitor(manifest io.Reader) (*monv1.ServiceMonitor, error) {
	sm := monv1.ServiceMonitor{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&sm)
	if err != nil {
		return nil, err
	}

	return &sm, nil
}

func NewDeployment(manifest io.Reader) (*appsv1.Deployment, error) {
	d := appsv1.Deployment{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func NewStatefulSet(manifest io.Reader) (*appsv1.StatefulSet, error) {
	s := appsv1.StatefulSet{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func NewIngress(manifest io.Reader) (*v1beta1.Ingress, error) {
	i := v1beta1.Ingress{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func NewAPIService(manifest io.Reader) (*apiregistrationv1beta1.APIService, error) {
	s := apiregistrationv1beta1.APIService{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func NewSecurityContextConstraints(manifest io.Reader) (*securityv1.SecurityContextConstraints, error) {
	s := securityv1.SecurityContextConstraints{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// HashTrustedCA synthesizes a configmap just by copying "ca-bundle.crt" from the given configmap
// and naming it by hashing the contents of "ca-bundle.crt".
// It adds "monitoring.openshift.io/name" and "monitoring.openshift.io/hash" labels.
// Any other labels from the given configmap are discarded.
//
// It returns nil, if the given configmap does not contain the "ca-bundle.crt" the data key
// or data is empty string.
func (f *Factory) HashTrustedCA(caBundleCM *v1.ConfigMap, prefix string) *v1.ConfigMap {
	caBundle, ok := caBundleCM.Data["ca-bundle.crt"]
	if !ok || caBundle == "" {
		// We return here instead of erroring out as we need
		// "ca-bundle.crt" to be there. This can mean that
		// the CA was not propagated yet. In that case we
		// will catch this on next sync loop.
		return nil
	}

	h := fnv.New64()
	h.Write([]byte(caBundle))
	hash := strconv.FormatUint(h.Sum64(), 32)

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-monitoring",
			Name:      fmt.Sprintf("%s-trusted-ca-bundle-%s", prefix, hash),
			Labels: map[string]string{
				"monitoring.openshift.io/name": prefix,
				"monitoring.openshift.io/hash": hash,
			},
		},
		Data: map[string]string{
			"ca-bundle.crt": caBundle,
		},
	}
}

// HashSecret synthesizes a secret by setting the given data
// and naming it by hashing the values of the given data.
//
// For simplicity, data is expected to be given in a key-value format,
// i.e. HashSecret(someSecret, value1, key1, value2, key2, ...).
//
// It adds "monitoring.openshift.io/name" and "monitoring.openshift.io/hash" labels.
// Any other labels from the given secret are discarded.
//
// It still returns a secret if the given secret does not contain any data.
func (f *Factory) HashSecret(secret *v1.Secret, data ...string) (*v1.Secret, error) {
	h := fnv.New64()
	m := make(map[string][]byte)

	var err error
	for i := 0; i < len(data)/2; i++ {
		k := data[i*2]
		v := []byte(data[i*2+1])
		_, err = h.Write(v)
		m[k] = v
	}
	if err != nil {
		return nil, errors.Wrap(err, "error hashing tls data")
	}
	hash := strconv.FormatUint(h.Sum64(), 32)

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: secret.GetNamespace(),
			Name:      fmt.Sprintf("%s-%s", secret.GetName(), hash),
			Labels: map[string]string{
				"monitoring.openshift.io/name": secret.GetName(),
				"monitoring.openshift.io/hash": hash,
			},
		},
		Data: m,
	}, nil
}

func trustedCABundleVolumeMount(name string) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      name,
		ReadOnly:  true,
		MountPath: "/etc/pki/ca-trust/extracted/pem/",
	}
}

func trustedCABundleVolume(configMapName, volumeName string) v1.Volume {
	yes := true

	return v1.Volume{
		Name: volumeName,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: configMapName,
				},
				Optional: &yes,
			},
		},
	}
}
