# Observatorium Operator

:warning: This project is a work in progress.

## How to deploy (Kubernetes)

### Prerequisites

#### S3 storage endpoint. For testing purposes you may use minio as follows

```
kubectl create namespace minio
kubectl create -f https://raw.githubusercontent.com/nmagnezi/observatorium-operator/master/hack/kubernetes/minio.yaml
```

#### Observatorium namespace thanos-objectstorage secret 

```
kubectl create -f https://raw.githubusercontent.com/nmagnezi/observatorium-operator/master/hack/kubernetes/observatorium_namespace_and_secret.yaml
```

#### RBAC configuration

```
kubectl create -f https://raw.githubusercontent.com/nmagnezi/observatorium-operator/master/hack/kubernetes/observatorium-operator-cluster-role.yaml
kubectl create -f https://raw.githubusercontent.com/nmagnezi/observatorium-operator/master/hack/kubernetes/observatorium-operator-cluster-role_binding.yaml
```

### Deploy the operator
* Deployment via image taken from [quay.io](https://quay.io/repository/nmagnezi/observatorium-operator?tab=tags)

#### Install CRDs
```
kubectl -n observatorium create -f https://raw.githubusercontent.com/nmagnezi/observatorium-operator/master/config/crd/bases/observatorium.observatorium_observatoria.yaml
kubectl apply -f https://raw.githubusercontent.com/coreos/kube-prometheus/master/manifests/setup/prometheus-operator-0servicemonitorCustomResourceDefinition.yaml
```
#### Install Operator Manager
```
kubectl -n observatorium create -f https://raw.githubusercontent.com/nmagnezi/observatorium-operator/master/hack/kubernetes/operator.yaml
```
### Deploy CR
```
kubectl -n observatorium create -f https://raw.githubusercontent.com/nmagnezi/observatorium-operator/master/hack/kubernetes/observatorium_v1alpha1_observatorium.yaml
```

### Expected Outcome
```
DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "observatorium", "request": "observatorium/observatorium-sample"}

$ kubectl -n observatorium get all
NAME                                              READY   STATUS    RESTARTS   AGE
pod/controller-manager-749b46864d-zqldk           1/1     Running   0          46s
pod/observatorium-querier-cache-dbc7b78f9-6bdx8   1/1     Running   0          14s
pod/observatorium-querier-cache-dbc7b78f9-hzkgz   1/1     Running   0          14s
pod/observatorium-querier-cache-dbc7b78f9-vmjk6   1/1     Running   0          14s
pod/thanos-compactor-0                            2/2     Running   0          15s
pod/thanos-querier-564cc64746-9fhzr               2/2     Running   0          15s
pod/thanos-querier-564cc64746-rzbbq               2/2     Running   0          15s
pod/thanos-querier-564cc64746-wwzcd               2/2     Running   0          15s
pod/thanos-receive-controller-6f784f6548-lfhm4    1/1     Running   0          14s
pod/thanos-receive-default-0                      1/2     Running   0          6s
pod/thanos-ruler-0                                1/1     Running   0          15s
pod/thanos-store-0                                2/2     Running   0          15s

NAME                                TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                         AGE
service/observatorium-cache         ClusterIP   10.96.3.255     <none>        9090/TCP                        15s
service/thanos-compactor            ClusterIP   10.96.9.18      <none>        10902/TCP                       15s
service/thanos-querier              ClusterIP   10.96.153.150   <none>        10901/TCP,9090/TCP              15s
service/thanos-receive-controller   ClusterIP   10.96.248.60    <none>        8080/TCP                        15s
service/thanos-receive-default      ClusterIP   None            <none>        10901/TCP,10902/TCP,19291/TCP   6s
service/thanos-ruler                ClusterIP   None            <none>        10901/TCP,10902/TCP             15s
service/thanos-store                ClusterIP   None            <none>        10901/TCP,10902/TCP             15s

NAME                                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/controller-manager            1/1     1            1           46s
deployment.apps/observatorium-querier-cache   3/3     3            3           14s
deployment.apps/thanos-querier                3/3     3            3           15s
deployment.apps/thanos-receive-controller     1/1     1            1           14s

NAME                                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/controller-manager-749b46864d           1         1         1       46s
replicaset.apps/observatorium-querier-cache-dbc7b78f9   3         3         3       14s
replicaset.apps/thanos-querier-564cc64746               3         3         3       15s
replicaset.apps/thanos-receive-controller-6f784f6548    1         1         1       14s

NAME                                      READY   AGE
statefulset.apps/thanos-compactor         1/1     15s
statefulset.apps/thanos-receive-default   0/3     6s
statefulset.apps/thanos-ruler             1/1     15s
statefulset.apps/thanos-store             1/1     15s
```

## How to deploy (OpenShift)
TBD

## Known limitations
* See [open issues](https://github.com/nmagnezi/observatorium-operator/issues)