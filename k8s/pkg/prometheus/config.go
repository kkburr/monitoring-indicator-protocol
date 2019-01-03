package prometheus

import (
	"github.com/pivotal/indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"github.com/pivotal/indicator-protocol/pkg/indicator"
	"github.com/pivotal/indicator-protocol/pkg/prometheus_alerts"
	"gopkg.in/yaml.v2"
	"log"
	"sync"
)

type Config struct {
	mu sync.Mutex
	indicatorDocuments map[string]*v1alpha1.IndicatorDocument
}

func NewConfig() *Config {
	return &Config{
		indicatorDocuments: map[string]*v1alpha1.IndicatorDocument{},
	}
}

func (c *Config) Upsert(i *v1alpha1.IndicatorDocument) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.indicatorDocuments[key(i)] = i
}

func (c *Config) Delete(i *v1alpha1.IndicatorDocument) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.indicatorDocuments, key(i))
}

func key(i *v1alpha1.IndicatorDocument) string {
	return i.Namespace + "/" + i.Name
}

// String will render out the prometheus config for alert rules.
func (c *Config) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	groups := make([]prometheus_alerts.Group, 0, len(c.indicatorDocuments))
	for k, i := range c.indicatorDocuments {
		doc := mapToDomainObject(i)
		alertDocument := prometheus_alerts.AlertDocumentFrom(doc)
		alertDocument.Groups[0].Name = k
		groups = append(groups, alertDocument.Groups[0])
	}

	out, err := yaml.Marshal(prometheus_alerts.Document{Groups: groups})
	if err != nil {
		log.Printf("Could not marshal alert rules: %s", err)
		return "groups: []"
	}

	return string(out)
}

func mapToDomainObject(i *v1alpha1.IndicatorDocument) indicator.Document {
	indicators := mapToDomainIndicators(i.Spec.Indicators)
	return indicator.Document{
		Product: indicator.Product{
			Name:    i.Spec.Product.Name,
			Version: i.Spec.Product.Version,
		},
		Metadata: i.Labels,
		Indicators: indicators,
	}
}

func mapToDomainIndicators(ids []v1alpha1.Indicator) []indicator.Indicator {
	indicators := make([]indicator.Indicator, 0, len(ids))
	for _, i :=range ids {
		indicators = append(indicators, indicator.Indicator{
			Name:          i.Name,
			PromQL:        i.Promql,
			Thresholds:    mapToDomainThreshold(i.Thresholds),
			Documentation: i.Documentation,
		})
	}
	return indicators
}

func mapToDomainThreshold(ths []v1alpha1.Threshold) []indicator.Threshold {
	thresholds := make([]indicator.Threshold, 0, len(ths))
	for _, t := range ths {
		op, val := resolveOperator(t)
		thresholds = append(thresholds, indicator.Threshold{
			Level:    t.Level,
			Operator: op,
			Value:    val,
		})
	}
	return thresholds
}

func resolveOperator(t v1alpha1.Threshold) (indicator.OperatorType, float64) {
	switch {
	case t.Lt != nil:
		return indicator.LessThan, *t.Lt
	case t.Lte != nil:
		return indicator.LessThanOrEqualTo, *t.Lte
	case t.Eq != nil:
		return indicator.EqualTo, *t.Eq
	case t.Neq != nil:
		return indicator.NotEqualTo, *t.Neq
	case t.Gte != nil:
		return indicator.GreaterThanOrEqualTo, *t.Gte
	case t.Gt != nil:
		return indicator.GreaterThan, *t.Gt
	}

	return indicator.LessThan, 0
}
