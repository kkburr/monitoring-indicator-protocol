#!/usr/bin/env bash
set -efu

export CERTS=./test_fixtures
export REGISTRY_PORT=8091
export REGISTRY_HOST="https://localhost:${REGISTRY_PORT}"
export DASHBOARD_PORT=8092
export INDICATOR_DOCUMENTS='./example_indicators.yml'

echo "Running grafana agent"
go run cmd/grafana_dashboard_controller/main.go \
  -registry ${REGISTRY_HOST} \
  -output-directory /tmp/dashboards \
  -tls-pem-path ${CERTS}/client.pem \
  -tls-key-path ${CERTS}/client.key \
  -tls-root-ca-pem ${CERTS}/root.pem \
  -tls-server-cn localhost \ &

wait