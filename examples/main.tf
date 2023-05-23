terraform {
  required_providers {
    openslo = {
      source = "registry.terraform.io/arctiq/openslo"
    }
  }
}

data "openslo_openslo" "test" {
  yaml_input = <<EOF
apiVersion: openslo/v1
kind: SLI
metadata:
  name: default-success-rate
  displayName: string
spec:
  description: string 
  ratioMetric:
    counter: true
    good:
      metricSource:
        type: datadog
        spec:
          query: sum:api.requests.status_ok{*}.as_count()
    total:
      metricSource:
        type: datadog
        spec:
          query: sum:api.requests{*}.as_count()
---
apiVersion: openslo/v1
kind: AlertPolicy
metadata:
  name: default
  displayName: Alert Policy
spec:
  conditions:
  - conditionRef: urgent-condition
  notificationTargets:
  - target: slack
---
apiVersion: openslo/v1
kind: Service
metadata:
  name: my-service
  displayName: My Service
spec:
  description: This service does blablabla
---
apiVersion: openslo/v1
kind: SLO
metadata:
  name: string
  displayName: string
spec:
  description: My service returns good responses 99.5% of the time
  service: my-service
  indicatorRef: default-success-rate
  timeWindow:
  - duration: 30d
  budgetingMethod: Occurrences
  alertPolicies:
  - kind: AlertPolicy
    metadata:
      name: success-alert-urgent
      displayName: "[URGENT] My service is in a degraded state"
    spec:
      conditions:
      - conditionRef: urgent-condition
      notificationTargets:
      - target: slack
  - kind: AlertPolicy
    metadata:
      name: success-alert-high
      displayName: "[HIGH] My service is in a degraded state"
    spec:
      conditions:
      - conditionRef: high-condition
      notificationTargets:
      - target: slack
  - kind: AlertPolicy
    metadata:
      name: success-alert-warning
      displayName: "[WARNING] My service might not respect its SLO"
    spec:
      conditions:
      - conditionRef: warn-condition
      notificationTargets:
      - target: servicenow
  objectives:
  - target: 0.995
EOF
}

output "test" {
  value = data.openslo_openslo.test.slos["string"].objectives[0]
}