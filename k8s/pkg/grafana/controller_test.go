package grafana_test

import (
	. "github.com/onsi/gomega"
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/grafana"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestController(t *testing.T) {
	t.Run("it adds", func(t *testing.T) {
		g := NewGomegaWithT(t)
		spyConfigMapEditor := &spyConfigMapEditor{g: g}
		controller := grafana.NewController(spyConfigMapEditor)

		controller.OnAdd(indicatorDocument())

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

		spyConfigMapEditor.expectThatNothingWasCreated()
	})

	t.Run("it updates", func(t *testing.T) {
		g := NewGomegaWithT(t)

		spyConfigMapEditor := &spyConfigMapEditor{g: g}
		p := grafana.NewController(spyConfigMapEditor)

		p.OnUpdate(nil, indicatorDocument())

		spyConfigMapEditor.expectUpdated([]*v1.ConfigMap{
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

	t.Run("fails to update a non-indicators", func(t *testing.T) {
		g := NewGomegaWithT(t)
		spyConfigMapEditor := &spyConfigMapEditor{g: g}
		controller := grafana.NewController(spyConfigMapEditor)

		controller.OnUpdate(indicatorDocument(), 616)
		spyConfigMapEditor.expectThatNothingWasUpdated()

		controller.OnUpdate("asdf", indicatorDocument())
		spyConfigMapEditor.expectThatNothingWasUpdated()
	})

	t.Run("does not update when new and old objects are the same", func(t *testing.T) {
		g := NewGomegaWithT(t)
		spyConfigMapEditor := &spyConfigMapEditor{g: g}
		controller := grafana.NewController(spyConfigMapEditor)

		controller.OnUpdate(nil, nil)
		spyConfigMapEditor.expectThatNothingWasUpdated()
		controller.OnUpdate(indicatorDocument(), indicatorDocument())
		spyConfigMapEditor.expectThatNothingWasUpdated()
	})

	t.Run("it deletes", func(t *testing.T) {
		g := NewGomegaWithT(t)

		spyConfigMapEditor := &spyConfigMapEditor{g: g}
		p := grafana.NewController(spyConfigMapEditor)

		p.OnDelete(indicatorDocument())

		uid := types.UID("some-uid")
		spyConfigMapEditor.expectDeleted([]deleteCall{
			{
				name: "rabbit-mq-resource-name-620771403",
				do: &metav1.DeleteOptions{
					Preconditions: &metav1.Preconditions{
						UID: &uid,
					},
				},
			},
		})
	})

	t.Run("fails to delete a non-indicators", func(t *testing.T) {
		g := NewGomegaWithT(t)
		spyConfigMapEditor := &spyConfigMapEditor{g: g}
		controller := grafana.NewController(spyConfigMapEditor)

		controller.OnDelete("non-indicator")

		spyConfigMapEditor.expectThatNothingWasDeleted()
	})

	// TODO: test that a namespace provided to the controller is set in the cm objects
}

type deleteCall struct {
	name string
	do   *metav1.DeleteOptions
}

type spyConfigMapEditor struct {
	g           *GomegaWithT
	createCalls []*v1.ConfigMap
	updateCalls []*v1.ConfigMap
	deleteCalls []deleteCall
}

func (s *spyConfigMapEditor) Create(cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	s.createCalls = append(s.createCalls, cm)
	return nil, nil
}

func (s *spyConfigMapEditor) Update(cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	s.updateCalls = append(s.updateCalls, cm)
	return nil, nil
}

func (s *spyConfigMapEditor) Delete(name string, do *metav1.DeleteOptions) error {
	s.deleteCalls = append(s.deleteCalls, deleteCall{name: name, do: do})
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
	s.g.Expect(s.updateCalls).To(HaveLen(len(cms)))
	for i, cm := range cms {
		s.g.Expect(s.updateCalls[i].Name).To(Equal(cm.Name))
		s.g.Expect(s.updateCalls[i].Labels).To(Equal(cm.Labels))
		s.g.Expect(s.updateCalls[i].Data["dashboard.json"]).To(MatchJSON(cm.Data["dashboard.json"]))
	}
}

func (s *spyConfigMapEditor) expectDeleted(dcs []deleteCall) {
	s.g.Expect(s.deleteCalls).To(Equal(dcs))
}

func (s *spyConfigMapEditor) expectThatNothingWasCreated() {
	s.g.Expect(s.createCalls).To(BeNil())
}
func (s *spyConfigMapEditor) expectThatNothingWasUpdated() {
	s.g.Expect(s.updateCalls).To(BeNil())
}
func (s *spyConfigMapEditor) expectThatNothingWasDeleted() {
	s.g.Expect(s.deleteCalls).To(BeNil())
}

func indicatorDocument() *v1alpha1.IndicatorDocument {
	return &v1alpha1.IndicatorDocument{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rabbit-mq-resource-name",
			UID:  types.UID("some-uid"),
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
}