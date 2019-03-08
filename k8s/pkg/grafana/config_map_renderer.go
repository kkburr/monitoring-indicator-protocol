package grafana

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/domain"
	"github.com/pivotal/monitoring-indicator-protocol/pkg/grafana_dashboard"
	"github.com/pivotal/monitoring-indicator-protocol/pkg/indicator"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mapper func(document indicator.Document) (string, error)

func ConfigMap(doc *v1alpha1.IndicatorDocument, m mapper) (*v1.ConfigMap, error) {
	if doc == nil {
		return nil, errors.New("source indicator document was empty")
	}

	if m == nil {
		m = func(document indicator.Document) (string, error) {
			dashboard := grafana_dashboard.DocumentToDashboard(document)
			data, err := json.Marshal(dashboard)
			if err != nil {
				return "", err
			}
			return string(data), nil
		}
	}
	jsonVal, err := m(domain.Map(doc))
	if err != nil {
		return nil, err
	}

	cmName := doc.Name + "-" + fmt.Sprintf("%x", sha1.Sum([]byte(doc.Name)))[:9]
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			UID:       doc.UID,
			Labels: map[string]string{
				"grafana_dashboard": "true",
			},
		},
		Data: map[string]string{
			"dashboard.json": jsonVal,
		},
	}

	return cm, nil
}
