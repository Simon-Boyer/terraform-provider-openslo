terraform {
  required_providers {
    openslo = {
      source = "registry.terraform.io/Simon-Boyer/openslo"
    }
  }
}

data "openslo_openslo" "test" {
  yaml_input = <<EOF
apiVersion: openslo/v1
kind: DataSource
metadata:
  name: default
  displayName: Default
spec:
  type: datadog
  spec:
    query: sum:api.requests.status_ok{*}.as_count()
---
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
        metricSourceRef: default
    total:
      metricSource:
        type: datadog
        spec:
          query: sum:api.requests{*}.as_count()
---
apiVersion: openslo/v1
kind: AlertCondition
metadata:
  name: string
  displayName: string
spec:
  description: string
  severity: string 
  condition:
    kind: string
    op: enum
    threshold: 1
    lookbackWindow: 1h
    alertAfter: 5m
---
apiVersion: openslo/v1
kind: AlertNotificationTarget
metadata:
  name: OnCallDevopsMailNotification
spec:
  description: Notifies by a mail message to the on-call devops mailing group
  target: email
---
apiVersion: openslo/v1
kind: AlertPolicy
metadata:
  name: default
  displayName: Alert Policy
spec:
  conditions:
  - conditionRef: string
  notificationTargets:
  - targetRef: OnCallDevopsMailNotification
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
  description: My service returns good responses 99.5 of the time
  service: my-service
  indicatorRef: default-success-rate
  timeWindow:
  - duration: 30d
  budgetingMethod: Occurrences
  alertPolicies:
  - alertPolicyRef: default
  objectives:
  - target: 0.995
    apiVersion: openslo_synthetics/v1
    kind: HTTPMonitor
    metadata:
        name: my-monitor
        displayName: My Monitor
    spec:
        url: https://my-host.com
        requests:
        - name: my-request
          description: This is a request
          headers:
            - name: my-header
              value: my-value
          body: test body
          method: POST
          path: /my-path
        - name: my-other-request
          description: This is a request
          headers:
          - name: my-other-header
            value: my-other-value
          method: GET
          path: /my-other-path
EOF
}

output "test" {
  value = data.openslo_openslo.test
}