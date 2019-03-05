package prometheus

import (
	"encoding/json"
	"fmt"
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"reflect"
)

type ConfigMapPatcher interface {
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ConfigMap, err error)
}

type ConfigRenderer interface  {
	Upsert(i *v1alpha1.IndicatorDocument)
	Delete(i *v1alpha1.IndicatorDocument)

	fmt.Stringer
}

type Controller struct {
	renderer ConfigRenderer
	patcher  ConfigMapPatcher
}

func NewController(configMap ConfigMapPatcher, cr ConfigRenderer) *Controller {
	return &Controller{
		patcher: configMap,
		renderer: cr,
	}
}

type patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func (c *Controller) OnAdd(obj interface{}) {
	i, ok := obj.(*v1alpha1.IndicatorDocument)
	if !ok {
		return
	}
	c.renderer.Upsert(i)
	c.patch()
}

func (c *Controller) OnUpdate(oldObj, newObj interface{}) {
	if !reflect.DeepEqual(oldObj, newObj){
		c.OnAdd(newObj)
	}
}

func (c *Controller) OnDelete(obj interface{}) {
	i, ok := obj.(*v1alpha1.IndicatorDocument)
	if !ok {
		return
	}
	c.renderer.Delete(i)
	c.patch()
}

func (c *Controller) patch() {
	patches := []patch{
		{
			Op:    "replace",
			Path:  "/data/alerts",
			Value: c.renderer.String(),
		},
	}

	data, err := json.Marshal(patches)
	if err != nil {
		log.Printf("unable to marshal JSON: %s", err)
		return
	}

	_, err = c.patcher.Patch("prometheus-server", types.JSONPatchType, data)
	if err != nil {
		log.Printf("unable to patch config map: %s", err)
	}
}
