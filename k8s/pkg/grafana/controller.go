package grafana

import (
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"log"
	"reflect"

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
	doc, ok := obj.(*v1alpha1.IndicatorDocument)
	if !ok {
		return
	}
	configMap, err := ConfigMap(doc, nil)
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
	newDoc, ok := newObj.(*v1alpha1.IndicatorDocument)
	if !ok {
		return
	}
	if oldObj != nil {
		oldDoc, ok := oldObj.(*v1alpha1.IndicatorDocument)
		if !ok {
			return
		}
		if reflect.DeepEqual(newDoc, oldDoc) {
			return
		}
	}
	configMap, err := ConfigMap(newDoc, nil)
	if err != nil {
		log.Printf("Failed to generate ConfigMap: %s", err)
	}
	_, err = c.cmEditor.Update(configMap)
	if err != nil {
		log.Printf("Failed to update ConfigMap: %s", err)
	}
}

// TODO: evaluate edge case where object might not exist
func (c *Controller) OnDelete(obj interface{}) {
	doc, _ := obj.(*v1alpha1.IndicatorDocument)
	configMap, err := ConfigMap(doc, nil)
	if err != nil {
		log.Printf("Failed to generate ConfigMap: %s", err)
	}
	err = c.cmEditor.Delete(configMap.Name, &metav1.DeleteOptions{
		Preconditions: &metav1.Preconditions{
			UID: &configMap.ObjectMeta.UID,
		},
	})
	if err != nil {
		log.Printf("Failed to delete ConfigMap: %s", err)
	}
}
