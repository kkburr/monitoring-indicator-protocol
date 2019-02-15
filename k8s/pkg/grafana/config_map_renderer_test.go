package grafana_test

import (
	"encoding/base64"
	"github.com/pivotal/monitoring-indicator-protocol/pkg/indicator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"errors"

	. "github.com/onsi/gomega"

	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/grafana"
)

func TestEmptyIndicatorDocument(t *testing.T) {
	g := NewGomegaWithT(t)

	var doc *v1alpha1.IndicatorDocument
	_, err := grafana.ConfigMap(doc, nil)
	g.Expect(err).To(HaveOccurred())
}

func TestNoLayoutGeneratesDefaultDashboard(t *testing.T) {
	g := NewGomegaWithT(t)

	doc := &v1alpha1.IndicatorDocument{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-name",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.IndicatorDocumentSpec{
			Product: v1alpha1.Product{
				Name:    "my_app",
				Version: "1.0.1",
			},
			Indicators: []v1alpha1.Indicator{
				{
					Name:   "latency",
					Promql: "histogram_quantile(0.9, latency)",
					Thresholds: []v1alpha1.Threshold{
						{
							Level: "critical",
							Gte:   floatVar(100.2),
						},
					},
					Documentation: map[string]string{
						"title": "90th Percentile Latency",
					},
				},
			},
		},
	}

	cm, err := grafana.ConfigMap(doc, func (indicator.Document) (string, error) {
		return "the-expected-json", nil
	})

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(cm.Name).To(Equal("test-name"))
	g.Expect(cm.Namespace).To(Equal("test-namespace"))
	b64Data := cm.Data["dashboard.json"]
	cmJSON, err := base64.StdEncoding.DecodeString(b64Data)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(string(cmJSON)).To(Equal("the-expected-json"))
}

func TestDashboardMapperError(t *testing.T) {
	g := NewGomegaWithT(t)

	doc := &v1alpha1.IndicatorDocument{}

	_, err := grafana.ConfigMap(doc, func (indicator.Document) (string, error) {
		return "", errors.New("some-error")
	})

	g.Expect(err).To(HaveOccurred())
}

func floatVar(f float64) *float64 {
	return &f
}
