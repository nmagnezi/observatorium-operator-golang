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
pod/controller-manager-749b46864d-ftg7l           1/1     Running   0          63s
pod/observatorium-api-77448bc78d-wvvdh            1/1     Running   0          48s
pod/observatorium-querier-cache-dbc7b78f9-247c6   1/1     Running   0          47s
pod/observatorium-querier-cache-dbc7b78f9-jjwr2   1/1     Running   0          47s
pod/observatorium-querier-cache-dbc7b78f9-q8b6l   1/1     Running   0          47s
pod/thanos-compactor-0                            2/2     Running   0          47s
pod/thanos-querier-564cc64746-6gsch               2/2     Running   0          48s
pod/thanos-querier-564cc64746-tt9vf               2/2     Running   0          48s
pod/thanos-querier-564cc64746-wtlgf               2/2     Running   0          48s
pod/thanos-receive-controller-6f784f6548-h2fdt    1/1     Running   0          47s
pod/thanos-receive-default-0                      2/2     Running   0          38s
pod/thanos-receive-default-1                      2/2     Running   0          33s
pod/thanos-receive-default-2                      2/2     Running   0          29s
pod/thanos-ruler-0                                1/1     Running   0          47s
pod/thanos-store-0                                2/2     Running   0          47s

NAME                                TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)                         AGE
service/observatorium-api           ClusterIP   10.96.197.3    <none>        8080/TCP                        50s
service/observatorium-cache         ClusterIP   10.96.94.202   <none>        9090/TCP                        50s
service/thanos-compactor            ClusterIP   10.96.102.18   <none>        10902/TCP                       50s
service/thanos-querier              ClusterIP   10.96.153.12   <none>        10901/TCP,9090/TCP              50s
service/thanos-receive              ClusterIP   10.96.240.17   <none>        10901/TCP,10902/TCP,19291/TCP   23s
service/thanos-receive-controller   ClusterIP   10.96.111.63   <none>        8080/TCP                        50s
service/thanos-receive-default      ClusterIP   None           <none>        10901/TCP,10902/TCP,19291/TCP   39s
service/thanos-ruler                ClusterIP   None           <none>        10901/TCP,10902/TCP             50s
service/thanos-store                ClusterIP   None           <none>        10901/TCP,10902/TCP             50s

NAME                                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/controller-manager            1/1     1            1           64s
deployment.apps/observatorium-api             1/1     1            1           49s
deployment.apps/observatorium-querier-cache   3/3     3            3           49s
deployment.apps/thanos-querier                3/3     3            3           49s
deployment.apps/thanos-receive-controller     1/1     1            1           48s

NAME                                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/controller-manager-749b46864d           1         1         1       64s
replicaset.apps/observatorium-api-77448bc78d            1         1         1       49s
replicaset.apps/observatorium-querier-cache-dbc7b78f9   3         3         3       49s
replicaset.apps/thanos-querier-564cc64746               3         3         3       49s
replicaset.apps/thanos-receive-controller-6f784f6548    1         1         1       48s

NAME                                      READY   AGE
statefulset.apps/thanos-compactor         1/1     49s
statefulset.apps/thanos-receive-default   3/3     39s
statefulset.apps/thanos-ruler             1/1     49s
statefulset.apps/thanos-store             1/1     49s
```

## How to deploy (OpenShift)
TBD

## Known limitations
* See [open issues](https://github.com/nmagnezi/observatorium-operator/issues)