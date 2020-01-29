package tasks

import (
	"github.com/nmagnezi/observatorium-operator/client"
	"github.com/nmagnezi/observatorium-operator/manifests"
	"github.com/pkg/errors"
)

type ThanosReceiveControllerTask struct {
	client  *client.Client
	factory *manifests.Factory
}

func NewThanosReceiveControllerTask(client *client.Client, factory *manifests.Factory) *ThanosReceiveControllerTask {
	return &ThanosReceiveControllerTask{
		client:  client,
		factory: factory,
	}
}

func (t *ThanosReceiveControllerTask) Run() error {

	cm, err := t.factory.ThanosReceiveControllerConfigMap()
	if err != nil {
		return errors.Wrap(err, "initializingThanos Receive Controller ConfigMap failed")
	}

	err = t.client.CreateOrUpdateConfigMap(cm)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Controller ConfigMap failed")
	}

	s, err := t.factory.ThanosReceiveControllerService()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Receive Controller Service failed")
	}

	err = t.client.CreateOrUpdateService(s)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Controller Service failed")
	}

	sa, err := t.factory.ThanosReceiveControllerServiceAccount()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Receive Controller ServiceAccount failed")
	}

	err = t.client.CreateOrUpdateServiceAccount(sa)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Controller ServiceAccount failed")
	}

	dep, err := t.factory.ThanosReceiveControllerDeployment()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Receive Controller Deployment failed")
	}

	err = t.client.CreateOrUpdateDeployment(dep)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Controller Deployment failed")
	}

	rc, err := t.factory.ThanosReceiveControllerRoleConfig()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Receive Controller Role config failed")
	}

	err = t.client.CreateOrUpdateRole(rc)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Controller Role config failed")
	}

	rb, err := t.factory.ThanosReceiveControllerRoleBinding()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Receive Controller RoleBinding failed")
	}

	err = t.client.CreateOrUpdateRoleBinding(rb)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Controller RoleBinding failed")
	}

	sm, err := t.factory.ThanosReceiveControllerServiceMonitor()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Receive Controller ServiceMonitor failed")
	}

	err = t.client.CreateOrUpdateServiceMonitor(sm)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Controller ServiceMonitor failed")
	}

	s, err = t.factory.ThanosReceiveDefaultService()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Receive Default Service failed")
	}

	err = t.client.CreateOrUpdateService(s)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Default Service failed")
	}

	ds, err := t.factory.ThanosReceiveDefaultStatefulSet()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Receive Default StatefulSet failed")
	}

	err = t.client.CreateOrUpdateStatefulSet(ds)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Default StatefulSet failed")
	}

	s, err = t.factory.ThanosReceiveService()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Receive Service failed")
	}

	err = t.client.CreateOrUpdateService(s)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Receive Service failed")
	}
	return nil
}
