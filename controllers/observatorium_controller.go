/*

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/nmagnezi/observatorium-operator/client"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrl_client "sigs.k8s.io/controller-runtime/pkg/client"

	observatoriumv1alpha1 "github.com/nmagnezi/observatorium-operator/api/v1alpha1"
	"github.com/nmagnezi/observatorium-operator/manifests"
	"github.com/nmagnezi/observatorium-operator/tasks"
)

// "sigs.k8s.io/controller-runtime/pkg/client"

// ObservatoriumReconciler reconciles a Observatorium object
type ObservatoriumReconciler struct {
	Client    *client.Client
	CrdClient ctrl_client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
}

// +kubebuilder:rbac:groups=observatorium.observatorium,resources=observatoria,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=observatorium.observatorium,resources=observatoria/status,verbs=get;update;patch

func (r *ObservatoriumReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("observatorium", req.NamespacedName)
	log.Info("Start Reconcile")

	var observatorium observatoriumv1alpha1.Observatorium
	if err := r.CrdClient.Get(ctx, req.NamespacedName, &observatorium); err != nil {
		log.Error(err, "unable to fetch Observatorium")
		return ctrl.Result{}, ctrl_client.IgnoreNotFound(err)
	}

	factory := manifests.NewFactory(req.NamespacedName.Namespace, "r.namespaceUserWorkload123", observatorium)
	tl := tasks.NewTaskRunner(
		r.Client,
		[]*tasks.TaskSpec{
			tasks.NewTaskSpec("Updating Thanos Querier", tasks.NewThanosQuerierTask(r.Client, factory)),
			tasks.NewTaskSpec("Updating Thanos Querier Cache", tasks.NewThanosQuerierCacheTask(r.Client, factory)),
			tasks.NewTaskSpec("Updating Thanos Compactor", tasks.NewThanosCompactorTask(r.Client, factory)),
			tasks.NewTaskSpec("Updating Thanos Store", tasks.NewThanosStoreTask(r.Client, factory)),
			tasks.NewTaskSpec("Updating Thanos Ruler", tasks.NewThanosRulerTask(r.Client, factory)),
			tasks.NewTaskSpec("Updating Thanos Receive Controller", tasks.NewThanosReceiveControllerTask(r.Client, factory)),
		},
	)

	taskName, err := tl.RunAll()
	if err != nil {
		// klog.Infof("Updating ClusterOperator status to failed. Err: %v", err)
		// failedTaskReason := strings.Join(strings.Fields(taskName+"Failed"), "")
		// reportErr := r.client.StatusReporter().SetFailed(err, failedTaskReason)
		// if reportErr != nil {
		// 	klog.Errorf("error occurred while setting status to failed: %v", reportErr)
		// }
		log.Info("Finish Reconcile Failed")
		r.Log.Error(err, taskName)
		return ctrl.Result{}, err
	}
	log.Info("Finish Reconcile Success")
	return ctrl.Result{}, nil
}

func (r *ObservatoriumReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&observatoriumv1alpha1.Observatorium{}).
		Complete(r)
}
