package tasks

import (
	"github.com/nmagnezi/observatorium-operator/client"
	"github.com/nmagnezi/observatorium-operator/manifests"
	"github.com/pkg/errors"
)

type ThanosStoreTask struct {
	client  *client.Client
	factory *manifests.Factory
}

func NewThanosStoreTask(client *client.Client, factory *manifests.Factory) *ThanosStoreTask {
	return &ThanosStoreTask{
		client:  client,
		factory: factory,
	}
}

func (t *ThanosStoreTask) Run() error {

	svc, err := t.factory.ThanosStoreService()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Store Service failed")
	}

	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Store Service failed")
	}

	s, err := t.factory.ThanosStoreStatefulSet()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Store Deployment failed")
	}

	err = t.client.CreateOrUpdateStatefulSet(s)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Store Deployment failed")
	}
	return nil
}
