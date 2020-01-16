package tasks

import (
	"github.com/nmagnezi/observatorium-operator/client"
	"github.com/nmagnezi/observatorium-operator/manifests"
	"github.com/pkg/errors"
)

type ThanosCompactorTask struct {
	client  *client.Client
	factory *manifests.Factory
}

func NewThanosCompactorTask(client *client.Client, factory *manifests.Factory) *ThanosCompactorTask {
	return &ThanosCompactorTask{
		client:  client,
		factory: factory,
	}
}

func (t *ThanosCompactorTask) Run() error {

	svc, err := t.factory.ThanosCompactorService()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Compactor Service failed")
	}

	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Compactor Service failed")
	}

	s, err := t.factory.ThanosCompactorStatefulSet()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Compactor Deployment failed")
	}

	err = t.client.CreateOrUpdateStatefulSet(s)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Compactor Deployment failed")
	}
	return nil
}
