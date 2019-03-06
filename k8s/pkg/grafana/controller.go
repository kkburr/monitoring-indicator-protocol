package grafana

import (
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"log"

	//"fmt"
	//"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapEditor interface {
	Create(*v1.ConfigMap) (*v1.ConfigMap, error)
	Update(*v1.ConfigMap) (*v1.ConfigMap, error)
	Delete(name string, options *metav1.DeleteOptions) error
}

type Controller struct {
	cmEditor ConfigMapEditor
}

func NewController(configMap ConfigMapEditor) *Controller {
	return &Controller{
		cmEditor: configMap,
	}
}

// TODO: what happens when you have two config maps with the same filename (ie: dashboard.json)

// TODO: evaluate edge case where object might already exist
func (c *Controller) OnAdd(obj interface{}) {
	d, ok := obj.(*v1alpha1.IndicatorDocument)
	if !ok {
		return
	}
	configMap, err := ConfigMap(d, nil)
	if err != nil {
		log.Printf("Failed to generate ConfigMap: %s", err)
	}
	_, err = c.cmEditor.Create(configMap)
	if err != nil {
		log.Printf("Failed to create ConfigMap: %s", err)
	}
}

// TODO: evaluate edge case where object might not exist
func (c *Controller) OnUpdate(oldObj, newObj interface{}) {
	// convert to config map
	// call update
}

// TODO: evaluate edge case where object might not exist
func (c *Controller) OnDelete(obj interface{}) {
	// convert to config map
	// call delete
}
