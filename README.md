# Indicator Protocol

This is an **observability as code** project which allows developers to define and expose performance, scaling,
and service level indicators for monitoring and alerting.
The indicator definition ideally lives in the same repository as the code and is automatically registered when the code is deployed.

There are 3 main uses cases for this project: Generating documentation, validating against an actual deployment's data,
and keeping a registry of indicators for use in monitoring tools such as prometheus alert manager and grafana.

See [the wiki](https://github.com/pivotal/indicator-protocol/wiki) for more detailed information and documentation.

## Running tests

Use the provided script to run tests: `./scripts/test.sh`

## Starting a registry and registry agent

1. A script is provided to start both components: `./scripts/run_registry_and_agent.sh`
1. To verify that this worked, open another terminal window, and use `./scripts/curl_indicators.sh` (or `./scripts/curl_indicators.sh | jq` if you have access to `jq`)
   * Your result should look something like this:

   ```json
   [
     {
       "apiVersion": "v0",
       "product": {
         "name": "my-component",
         "version": "1.2.3"
       },
       "metadata": {
         "deployment": "my-service-deployment",
         "source_id": "my-metric-source"
       },
       "indicators": [
         {
           "name": "only_in_example_yml",
           "promql": "test_query"
         },
         {
           "name": "doc_performance_indicator",
           "promql": "avg_over_time(demo_latency{source_id=\"my-metric-source\",deployment=\"my-service-deployment\"}[5m])",
           "thresholds": [
             {
               "level": "warning",
               "operator": "gte",
               "value": 50
             },
             {
               "level": "critical",
               "operator": "gt",
               "value": 100
             }
           ],
           "documentation": {
             "description": "This is a valid markdown description.\n\n**Use**: This indicates nothing. It is placeholder text.\n\n**Type**: Gauge\n**Frequency**: 60 s\n",
             "measurement": "Average latency over last 5 minutes per instance",
             "recommended_response": "Panic! Run around in circles flailing your arms.",
             "threshold_note": "These are environment specific",
             "title": "Doc Performance Indicator"
           },
           "presentation": {
             "chartType": "line",
             "currentValue": false,
             "interval": "1m0s"
           }
         },
         {
           "name": "success_percentage",
           "promql": "success_percentage_promql{source_id=\"origin\"}",
           "documentation": {
             "title": "Success Percentage"
           }
         }
       ],
       "layout": {
         "title": "Monitoring Document Product",
         "description": "Document description",
         "sections": [
           {
             "title": "Indicators",
             "description": "This section includes indicators",
             "indicators": [
               "doc_performance_indicator"
             ]
           }
         ],
         "owner": "Example Team"
       }
     }
   ]

   ```

   Specifically, you should get a single product called `my-component` with 3 indicators, some metadata,
   and a layout section.
