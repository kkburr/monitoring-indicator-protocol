package grafana

import (
	//"fmt"
	//"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapEditor interface {
	Create(*v1.ConfigMap) (*v1.ConfigMap, error)
	Update(*v1.ConfigMap) (*v1.ConfigMap, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1.ConfigMap, error)
}

type Controller struct {
	cmEditor ConfigMapEditor
}

func NewController(configMap ConfigMapEditor) *Controller {
	return &Controller{
		cmEditor: configMap,
	}
}

func (c *Controller) OnAdd(obj interface{}) {
	// convert to config map
	// check for existance
	// update or create depending on existance
}

func (c *Controller) OnUpdate(oldObj, newObj interface{}) {
	// convert to config map
	// check for existance
	// update or create depending on existance
}

func (c *Controller) OnDelete(obj interface{}) {
	// convert to config map
	// delete
}
