package grafana_test

import (
	. "github.com/onsi/gomega"
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/grafana"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestController(t *testing.T) {
	t.Run("it adds", func(t *testing.T) {
		g := NewGomegaWithT(t)

		spyConfigMapEditor := &spyConfigMapEditor{g: g}

		controller := grafana.NewController(spyConfigMapEditor)

		i := &v1alpha1.IndicatorDocument{
			ObjectMeta: metav1.ObjectMeta{
				Name: "rabbit-mq-resource-name",
			},
			Spec: v1alpha1.IndicatorDocumentSpec{
				Product: v1alpha1.Product{
					Name:    "rabbit-mq-product-name",
					Version: "v1.0",
				},
				Indicators: []v1alpha1.Indicator{
					{
						Name:   "qps",
						Promql: "rate(qps)",
					},
				},
				Layout: v1alpha1.Layout{
					Title: "rabbit-mq-layout-title",
				},
			},
		}

		controller.OnAdd(i)

		spyConfigMapEditor.expectCreated([]*v1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rabbit-mq-resource-name-620771403",
					Labels: map[string]string{
						"grafana_dashboard": "true",
					},
				},
				Data: map[string]string{
					"dashboard.json": `{
					  "title": "rabbit-mq-layout-title",
					  "rows": [
					    {
					      "title": "qps",
					      "panels": [
					        {
					          "title": "qps",
					          "type": "graph",
					          "targets": [
					            {
					              "expr": "rate(qps)"
					            }
					          ],
					          "thresholds": null
					        }
					      ]
					    }
					  ]
					}`,
				},
			},
		})
	})

	t.Run("fails to add a non-indicator", func(t *testing.T) {
		g := NewGomegaWithT(t)
		spyConfigMapEditor := &spyConfigMapEditor{g: g}
		controller := grafana.NewController(spyConfigMapEditor)

		controller.OnAdd(666)

		spyConfigMapEditor.expectCreated(nil)
	})

	//t.Run("it updates existing indicators", func(t *testing.T) {
	//	g := NewGomegaWithT(t)
	//
	//	spyConfigMapEditor := &spyConfigMapEditor{g: g}
	//	spyConfigRenderer := &spyConfigRenderer{
	//		g: g,
	//		config: "new-config",
	//	}
	//	p := prometheus.NewController(spyConfigMapEditor, spyConfigRenderer)
	//
	//	i1 := &v1alpha1.IndicatorDocument{
	//		ObjectMeta: metav1.ObjectMeta{
	//			Name:            "rabbit-mq-monitoring-1",
	//		},
	//	}
	//	i2 := &v1alpha1.IndicatorDocument{
	//		ObjectMeta: metav1.ObjectMeta{
	//			Name:            "rabbit-mq-monitoring-2",
	//		},
	//	}
	//
	//	p.OnAdd(i1)
	//	p.OnUpdate(i1, i2)
	//
	//	spyConfigRenderer.assertUpsert(0, i1)
	//	spyConfigRenderer.assertUpsert(1, i2)
	//	spyConfigMapEditor.expectPatches([]string{
	//		"new-config",
	//		"new-config",
	//	})
	//})
	//
	//t.Run("it does not upsert if the indicator is unchanged", func(t *testing.T) {
	//	g := NewGomegaWithT(t)
	//
	//	spyConfigMapEditor := &spyConfigMapEditor{g: g}
	//	spyConfigRenderer := &spyConfigRenderer{
	//		g: g,
	//		config: "new-config",
	//	}
	//	p := prometheus.NewController(spyConfigMapEditor, spyConfigRenderer)
	//
	//	i1 := &v1alpha1.IndicatorDocument{
	//		ObjectMeta: metav1.ObjectMeta{
	//			Name:            "rabbit-mq-monitoring-1",
	//		},
	//	}
	//	i2 := &v1alpha1.IndicatorDocument{
	//		ObjectMeta: metav1.ObjectMeta{
	//			Name:            "rabbit-mq-monitoring-1",
	//		},
	//	}
	//
	//	p.OnAdd(i1)
	//	p.OnUpdate(i1, i2)
	//
	//	spyConfigRenderer.assertUpsert(0, i1)
	//	spyConfigRenderer.assertUpsertLen(1)
	//	spyConfigMapEditor.expectPatches([]string{
	//		"new-config",
	//	})
	//})
	//
	//t.Run("it deletes existing indicators", func(t *testing.T) {
	//	g := NewGomegaWithT(t)
	//
	//	spyConfigMapEditor := &spyConfigMapEditor{g: g}
	//	spyConfigRenderer := &spyConfigRenderer{
	//		g: g,
	//		config: "new-config",
	//	}
	//	p := prometheus.NewController(spyConfigMapEditor, spyConfigRenderer)
	//
	//	i := &v1alpha1.IndicatorDocument{
	//		ObjectMeta: metav1.ObjectMeta{
	//			Name:            "rabbit-mq-monitoring-1",
	//		},
	//	}
	//
	//	p.OnDelete(i)
	//
	//	spyConfigRenderer.assertDelete(0, i)
	//	spyConfigRenderer.assertDeleteLen(1)
	//	spyConfigMapEditor.expectPatches([]string{
	//		"new-config",
	//	})
	//})
	//
	//t.Run("it does nothing when given non-indicators", func(t *testing.T) {
	//	g := NewGomegaWithT(t)
	//
	//	spyConfigMapEditor := &spyConfigMapEditor{g: g}
	//	spyConfigRenderer := &spyConfigRenderer{g: g}
	//	p := prometheus.NewController(spyConfigMapEditor, spyConfigRenderer)
	//
	//	p.OnAdd(nil)
	//	p.OnAdd("nothing")
	//	p.OnAdd(42)
	//
	//	p.OnUpdate(nil,nil)
	//	p.OnUpdate("nothing", "something")
	//	p.OnUpdate(42, 23)
	//
	//	p.OnDelete(nil)
	//	p.OnDelete("nothing")
	//	p.OnDelete(42)
	//
	//	spyConfigRenderer.assertUpsertLen(0)
	//	g.Expect(spyConfigMapEditor.patchCalled).To(BeFalse())
	//})
}

type spyConfigMapEditor struct {
	g           *GomegaWithT
	createCalls []*v1.ConfigMap
	updateCalls []*v1.ConfigMap
	deleteCalls []string
}

func (s *spyConfigMapEditor) Create(cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	s.createCalls = append(s.createCalls, cm)
	return nil, nil
}

func (s *spyConfigMapEditor) Update(cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	s.updateCalls = append(s.updateCalls, cm)
	return nil, nil
}

func (s *spyConfigMapEditor) Delete(name string, _ *metav1.DeleteOptions) error {
	s.deleteCalls = append(s.deleteCalls, name)
	return nil
}

func (s *spyConfigMapEditor) expectCreated(cms []*v1.ConfigMap) {
	s.g.Expect(s.createCalls).To(HaveLen(len(cms)))
	for i, cm := range cms {
		s.g.Expect(s.createCalls[i].Name).To(Equal(cm.Name))
		s.g.Expect(s.createCalls[i].Labels).To(Equal(cm.Labels))
		s.g.Expect(s.createCalls[i].Data["dashboard.json"]).To(MatchJSON(cm.Data["dashboard.json"]))
	}
}

func (s *spyConfigMapEditor) expectUpdated(cms []*v1.ConfigMap) {
	s.g.Expect(s.updateCalls).To(Equal(cms))
}

func (s *spyConfigMapEditor) expectDeleted(names []string) {
	s.g.Expect(s.deleteCalls).To(Equal(names))
}
