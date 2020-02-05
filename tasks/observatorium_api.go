package tasks

import (
	"github.com/nmagnezi/observatorium-operator/client"
	"github.com/nmagnezi/observatorium-operator/manifests"
	"github.com/pkg/errors"
)

type ObservatoriumApiTask struct {
	client  *client.Client
	factory *manifests.Factory
}

func NewThanosObservatoriumApiTask(client *client.Client, factory *manifests.Factory) *ObservatoriumApiTask {
	return &ObservatoriumApiTask{
		client:  client,
		factory: factory,
	}
}

func (t *ObservatoriumApiTask) Run() error {

	svc, err := t.factory.ObservatoriumApiService()
	if err != nil {
		return errors.Wrap(err, "initializing Observatorium API Service failed")
	}

	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling Observatorium API Service failed")
	}

	dep, err := t.factory.ObservatoriumApiDeployment()
	if err != nil {
		return errors.Wrap(err, "initializing Observatorium API Deployment failed")
	}

	err = t.client.CreateOrUpdateDeployment(dep)
	if err != nil {
		return errors.Wrap(err, "reconciling Observatorium API Deployment failed")
	}
	return nil
}
