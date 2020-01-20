package tasks

import (
	"github.com/nmagnezi/observatorium-operator/client"
	"github.com/nmagnezi/observatorium-operator/manifests"
	"github.com/pkg/errors"
)

type ThanosQuerierCacheTask struct {
	client  *client.Client
	factory *manifests.Factory
}

func NewThanosQuerierCacheTask(client *client.Client, factory *manifests.Factory) *ThanosQuerierCacheTask {
	return &ThanosQuerierCacheTask{
		client:  client,
		factory: factory,
	}
}

func (t *ThanosQuerierCacheTask) Run() error {

	svc, err := t.factory.ThanosQuerierCacheService()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Querier Cache Service failed")
	}

	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Querier Cache Service failed")
	}

	c, err := t.factory.ThanosQuerierCacheConfigMap()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Querier Cache ConfigMap failed")
	}

	err = t.client.CreateOrUpdateConfigMap(c)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Querier Cache ConfigMap failed")
	}

	dep, err := t.factory.ThanosQuerierCacheDeployment()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Querier Cache Deployment failed")
	}

	err = t.client.CreateOrUpdateDeployment(dep)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Querier Cache Deployment failed")
	}
	return nil

}
