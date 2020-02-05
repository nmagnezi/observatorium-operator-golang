package tasks

import (
	"github.com/nmagnezi/observatorium-operator/client"
	"github.com/nmagnezi/observatorium-operator/manifests"
	"github.com/pkg/errors"
)

type ThanosRulerTask struct {
	client  *client.Client
	factory *manifests.Factory
}

func NewThanosRulerTask(client *client.Client, factory *manifests.Factory) *ThanosRulerTask {
	return &ThanosRulerTask{
		client:  client,
		factory: factory,
	}
}

func (t *ThanosRulerTask) Run() error {

	svc, err := t.factory.ThanosRulerService()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Ruler Service failed")
	}

	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Ruler Service failed")
	}

	s, err := t.factory.ThanosRulerStatefulSet()
	if err != nil {
		return errors.Wrap(err, "initializing Thanos Ruler StatefulSet failed")
	}

	err = t.client.CreateOrUpdateStatefulSet(s)
	if err != nil {
		return errors.Wrap(err, "reconciling Thanos Ruler StatefulSet failed")
	}
	return nil
}
