package provider

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestOpenSLODatasource_shouldbeValid_singleYamlSpec(t *testing.T) {
	// given
	yamlSpec := `
    apiVersion: openslo/v1
    kind: DataSource
    metadata:
      name: my-datasource
      displayName: My DataSource
    spec:
      type: datasource-type
      description: Datasource description
      connectionDetails:
        host: my-host
        port: my-port
        user: my-user
        password: my-password
      spec:
        custom-parameter: my-custom-parameter
`

	expected := DataSourceModel{
		Metadata: MetadataModel{
			DisplayName: "My DataSource",
			Name:        "my-datasource",
		},
		Type:        "datasource-type",
		Description: "Datasource description",
		ConnectionDetails: map[string]string{
			"host":     "my-host",
			"port":     "my-port",
			"user":     "my-user",
			"password": "my-password",
		},
		Spec: map[string]string{
			"custom-parameter": "my-custom-parameter",
		},
	}
	// when
	openslo, err := GetOpenSloData(yamlSpec, &diag.Diagnostics{})
	// then
	if err != nil {
		t.Error(err)
	}
	// and
	if !reflect.DeepEqual(openslo.Datasources["my-datasource"], expected) {
		t.Errorf("Expected %#v, but got %#v", expected, openslo.Datasources["my-datasource"])
	}
}

func TestOpenSLOService_shouldbeValid_singleYamlSpec(t *testing.T) {
	// given
	yamlSpec := `
apiVersion: openslo/v1
kind: Service
metadata:
    name: my-service
    displayName: My Service
spec:
    description: This service does blablabla
`
	expected := ServiceModel{
		Metadata: MetadataModel{
			Name:        "my-service",
			DisplayName: "My Service",
		},
		Description: "This service does blablabla",
	}
	// when
	openslo, err := GetOpenSloData(yamlSpec, &diag.Diagnostics{})
	// then
	if err != nil {
		t.Error(err)
	}
	// and
	if !reflect.DeepEqual(openslo.Services["my-service"], expected) {
		t.Errorf("Expected '%#v', but got '%#v'", expected, openslo.Services["my-service"])
	}
}

func TestOpenSLOSLI_shouldbeValid_singleYamlSpec(t *testing.T) {
	// given
	yamlSpec := `
apiVersion: openslo/v1
kind: SLI
metadata:
  name: string
  displayName: string
spec:
  description: string
  thresholdMetric:
    metricSource:
      type: string
      spec:
        query: a_query
  ratioMetric:
    counter: true
    good:
      metricSource:
        type: string
        spec:
            query: a_query
    bad:
      metricSource:
        type: string
        spec:
            query: a_query
    total:
      metricSource:
        type: string
        spec:
          query: a_query
    rawType: success
    raw:
      metricSource:
        type: string
        spec:
          query: a_query
`

	expected := SLIModel{
		Metadata: MetadataModel{
			Name:        "string",
			DisplayName: "string",
		},
		Description: "string",
		ThresholdMetric: MetricModel{
			MetricSource: DataSourceModel{
				Type: "string",
				Spec: map[string]string{
					"query": "a_query",
				},
			},
		},
		RatioMetric: RatioMetricModel{
			Counter: true,
			Good: MetricModel{
				MetricSource: DataSourceModel{
					Type: "string",
					Spec: map[string]string{
						"query": "a_query",
					},
				},
			},
			Bad: MetricModel{
				MetricSource: DataSourceModel{
					Type: "string",
					Spec: map[string]string{
						"query": "a_query",
					},
				},
			},
			Total: MetricModel{
				MetricSource: DataSourceModel{
					Type: "string",
					Spec: map[string]string{
						"query": "a_query",
					},
				},
			},
			RawType: "success",
			Raw: MetricModel{
				MetricSource: DataSourceModel{
					Type: "string",
					Spec: map[string]string{
						"query": "a_query",
					},
				},
			},
		},
	}

	// when
	openslo, err := GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	if !reflect.DeepEqual(openslo.Slis["string"], expected) {
		t.Errorf("Expected '%#v', but got '%#v'", expected, openslo.Slis["string"])
	}
}

func TestOpenSLOAlertPolicy_shouldbeValid_singleYamlSpec(t *testing.T) {
	// given
	yamlSpec := `
apiVersion: openslo/v1
kind: AlertPolicy
metadata:
    name: default
    displayName: Alert Policy
spec:
    conditions:
    - kind: AlertCondition
      metadata:
        name: cpu-usage-breach
        displayName: CPU Usage breaching
      spec:
        description: SLO burn rate for cpu-usage-breach exceeds 2
        severity: page
        condition:
          kind: burnrate
          op: lte
          threshold: 2
          lookbackWindow: 1h
          alertAfter: 5m
    notificationTargets:
    - target: slack
      description: Notify on slack
`
	expected := AlertPolicyModel{
		Metadata: MetadataModel{
			Name:        "default",
			DisplayName: "Alert Policy",
		},
		Conditions: []AlertConditionModel{
			{
				Metadata: MetadataModel{
					Name:        "cpu-usage-breach",
					DisplayName: "CPU Usage breaching",
				},
				Description: "SLO burn rate for cpu-usage-breach exceeds 2",
				Severity:    "page",
				Condition: AlertConditionModelCondition{
					Kind:           "burnrate",
					Op:             "lte",
					Threshold:      2,
					LookbackWindow: "1h",
					AlertAfter:     "5m",
				},
			},
		},
		NotificationTargets: []AlertNotificationTargetModel{
			{
				Target:      "slack",
				Description: "Notify on slack",
			},
		},
	}
	// when
	openslo, err := GetOpenSloData(yamlSpec, &diag.Diagnostics{})
	// then
	if err != nil {
		t.Error(err)
	}
	// and
	if !reflect.DeepEqual(openslo.Alert_policies["default"], expected) {
		t.Errorf("Expected %#v, but got %#v", expected, openslo.Alert_policies["default"])
	}
}

func TestOpenSLOSLO_shouldbeValid_singleYamlSpec(t *testing.T) {
	// given
	yamlSpec := `
    apiVersion: openslo/v1
    kind: SLO
    metadata:
      name: string
      displayName: string 
    spec:
      description: string 
      indicator:
        apiVersion: openslo/v1
        kind: SLI
        metadata:
            name: string
            displayName: string
        spec:
            description: string
            thresholdMetric:
                metricSource:
                    type: string
                    spec:
                        query: a_query
      timeWindow:
        - duration: 1d
          isRolling: true
      budgetingMethod: Occurrences
      objectives:
      - displayName: string
        op: lte 
        value: 0.5
        target: 0.99
        targetPercent: 99
        timeSliceTarget: 0.5
        timeSliceWindow: 1h
        indicator:
            apiVersion: openslo/v1
            kind: SLI
            metadata:
                name: string
                displayName: string
            spec:
                description: string
                thresholdMetric:
                    metricSource:
                        type: string
                        spec:
                            query: a_query
        compositeWeight: 0.8
      alertPolicies:
      - apiVersion: openslo/v1
        kind: AlertPolicy
        metadata:
            name: default
            displayName: Alert Policy
        spec:
            conditions:
            - kind: AlertCondition
              metadata:
                name: cpu-usage-breach
                displayName: CPU Usage breaching
              spec:
                description: SLO burn rate for cpu-usage-breach exceeds 2
                severity: page
                condition:
                  kind: burnrate
                  op: lte
                  threshold: 2
                  lookbackWindow: 1h
                  alertAfter: 5m
            notificationTargets:
            - target: slack
              description: Notify on slack
`

	expected := SLOModel{
		Metadata: MetadataModel{
			Name:        "string",
			DisplayName: "string",
		},
		Description: "string",
		Indicator: SLIModel{
			Metadata: MetadataModel{
				Name:        "string",
				DisplayName: "string",
			},
			Description: "string",
			ThresholdMetric: MetricModel{
				MetricSource: DataSourceModel{
					Type: "string",
					Spec: map[string]string{
						"query": "a_query",
					},
				},
			},
		},
		TimeWindow: []TimeWindowModel{
			{
				Duration:  "1d",
				IsRolling: true,
			},
		},
		BudgetingMethod: "Occurrences",
		Objectives: []ObjectiveModel{
			{
				DisplayName:     "string",
				Op:              "lte",
				Value:           0.5,
				Target:          0.99,
				TargetPercent:   99,
				TimeSliceTarget: 0.5,
				TimeSliceWindow: "1h",
				Indicator: SLIModel{
					Metadata: MetadataModel{
						Name:        "string",
						DisplayName: "string",
					},
					Description: "string",
					ThresholdMetric: MetricModel{
						MetricSource: DataSourceModel{
							Type: "string",
							Spec: map[string]string{
								"query": "a_query",
							},
						},
					},
				},
				CompositeWeight: 0.8,
			},
		},
		AlertPolicies: []AlertPolicyModel{
			{
				Metadata: MetadataModel{
					Name:        "default",
					DisplayName: "Alert Policy",
				},
				Conditions: []AlertConditionModel{
					{
						Metadata: MetadataModel{
							Name:        "cpu-usage-breach",
							DisplayName: "CPU Usage breaching",
						},
						Description: "SLO burn rate for cpu-usage-breach exceeds 2",
						Severity:    "page",
						Condition: AlertConditionModelCondition{
							Kind:           "burnrate",
							Op:             "lte",
							Threshold:      2,
							LookbackWindow: "1h",
							AlertAfter:     "5m",
						},
					},
				},
				NotificationTargets: []AlertNotificationTargetModel{
					{
						Target:      "slack",
						Description: "Notify on slack",
					},
				},
			},
		},
	}

	// when
	openslo, err := GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	if !reflect.DeepEqual(openslo.Slos["string"], expected) {
		t.Errorf("Expected %#v, but got %#v", expected, openslo.Slos["string"])
	}
}

func TestOpenSLOAlertConditions_shouldbeValid_singleYamlSpec(t *testing.T) {

	// given
	yamlSpec := `
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
    alertAfter: 2h
`

	expected := AlertConditionModel{
		Metadata: MetadataModel{
			Name:        "string",
			DisplayName: "string",
		},
		Description: "string",
		Severity:    "string",
		Condition: AlertConditionModelCondition{
			Kind:           "string",
			Op:             "enum",
			Threshold:      1,
			LookbackWindow: "1h",
			AlertAfter:     "2h",
		},
	}

	// when
	openslo, err := GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	if !reflect.DeepEqual(openslo.Alert_conditions["string"], expected) {
		t.Errorf("Expected %#v, but got %#v", expected, openslo.Alert_conditions["string"])
	}
}

func TestOpenSLOAlertNotificationTargets_shouldbeValid_singleYamlSpec(t *testing.T) {

	// given
	yamlSpec := `
apiVersion: openslo/v1
kind: AlertNotificationTarget
metadata:
  name: OnCallDevopsMailNotification
  displayName: Display name
spec:
  description: Notifies by a mail message to the on-call devops mailing group
  target: email
`

	expected := AlertNotificationTargetModel{
		Metadata: MetadataModel{
			Name:        "OnCallDevopsMailNotification",
			DisplayName: "Display name",
		},
		Target:      "email",
		Description: "Notifies by a mail message to the on-call devops mailing group",
	}

	// when
	openslo, err := GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	if !reflect.DeepEqual(openslo.Alert_notification_targets["OnCallDevopsMailNotification"], expected) {
		t.Errorf("Expected %#v, but got %#v", expected, openslo.Alert_notification_targets["OnCallDevopsMailNotification"])
	}
}

func TestOpenSLOAll_shouldbeValid_withRefLinks(t *testing.T) {

	// given
	yamlSpec := `
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
`

	dataSource := DataSourceModel{
		Metadata: MetadataModel{
			Name:        "default",
			DisplayName: "Default",
		},
		Type: "datadog",
		Spec: map[string]string{
			"query": "sum:api.requests.status_ok{*}.as_count()",
		},
	}
	dataSourceWithRef := dataSource
	dataSourceWithRef.MetricSourceRef = "default"

	service := ServiceModel{
		Metadata: MetadataModel{
			Name:        "my-service",
			DisplayName: "My Service",
		},
		Description: "This service does blablabla",
	}

	alertCondition := AlertConditionModel{
		Metadata: MetadataModel{
			Name:        "string",
			DisplayName: "string",
		},
		Description: "string",
		Severity:    "string",
		Condition: AlertConditionModelCondition{
			Kind:           "string",
			Op:             "enum",
			Threshold:      1,
			LookbackWindow: "1h",
			AlertAfter:     "5m",
		},
	}
	alertConditionWithRef := alertCondition
	alertConditionWithRef.ConditionRef = "string"

	alertNotificationTarget := AlertNotificationTargetModel{
		Metadata: MetadataModel{
			Name: "OnCallDevopsMailNotification",
		},
		Target:      "email",
		Description: "Notifies by a mail message to the on-call devops mailing group",
	}
	alertNotificationTargetWithRef := alertNotificationTarget
	alertNotificationTargetWithRef.TargetRef = "OnCallDevopsMailNotification"

	alertPolicy := AlertPolicyModel{
		Metadata: MetadataModel{
			Name:        "default",
			DisplayName: "Alert Policy",
		},
		Conditions: []AlertConditionModel{
			alertConditionWithRef,
		},
		NotificationTargets: []AlertNotificationTargetModel{
			alertNotificationTargetWithRef,
		},
	}
	alertPolicyWithRef := alertPolicy
	alertPolicyWithRef.AlertPolicyRef = "default"

	sli := SLIModel{
		Metadata: MetadataModel{
			Name:        "default-success-rate",
			DisplayName: "string",
		},
		Description: "string",
		RatioMetric: RatioMetricModel{
			Counter: true,
			Good: MetricModel{
				MetricSource: dataSourceWithRef,
			},
			Total: MetricModel{
				MetricSource: DataSourceModel{
					Type: "datadog",
					Spec: map[string]string{
						"query": "sum:api.requests{*}.as_count()",
					},
				},
			},
		},
	}

	slo := SLOModel{
		Metadata: MetadataModel{
			Name:        "string",
			DisplayName: "string",
		},
		Description:  "My service returns good responses 99.5 of the time",
		Service:      service,
		ServiceRef:   "my-service",
		Indicator:    sli,
		IndicatorRef: "default-success-rate",
		TimeWindow: []TimeWindowModel{
			{
				Duration: "30d",
			},
		},
		BudgetingMethod: "Occurrences",
		AlertPolicies: []AlertPolicyModel{
			alertPolicyWithRef,
		},
		Objectives: []ObjectiveModel{
			{
				Target:          0.995,
				CompositeWeight: 1,
			},
		},
	}

	// when
	openslo, err := GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	if !reflect.DeepEqual(openslo.Services["my-service"], service) {
		t.Errorf("Expected %#v, but got %#v", service, openslo.Alert_conditions["my-service"])
	}

	// and
	if !reflect.DeepEqual(openslo.Datasources["default"], dataSource) {
		t.Errorf("Expected %#v, but got %#v", dataSource, openslo.Datasources["default"])
	}

	// and
	if !reflect.DeepEqual(openslo.Alert_conditions["string"], alertCondition) {
		t.Errorf("Expected %#v, but got %#v", alertCondition, openslo.Alert_conditions["string"])
	}

	// and
	if !reflect.DeepEqual(openslo.Alert_notification_targets["OnCallDevopsMailNotification"], alertNotificationTarget) {
		t.Errorf("Expected %#v, but got %#v", alertNotificationTarget, openslo.Alert_notification_targets["OnCallDevopsMailNotification"])
	}

	// and
	if !reflect.DeepEqual(openslo.Alert_policies["default"], alertPolicy) {
		t.Errorf("Expected %#v, but got %#v", alertPolicy, openslo.Alert_policies["default"])
	}

	// and
	if !reflect.DeepEqual(openslo.Slis["default-success-rate"], sli) {
		t.Errorf("Expected %#v, but got %#v", sli, openslo.Slis["default-success-rate"])
	}

	// and
	if !reflect.DeepEqual(openslo.Slos["string"], slo) {
		t.Errorf("Expected %#v, but got %#v", slo, openslo.Slos["string"])
	}

}

func TestOpenSLO_shouldbeWarning_badApiVersion(t *testing.T) {
	// given
	yamlSpec := `
apiVersion: other/v1
kind: Service
metadata:
  name: my-service
  displayName: My Service
spec:
  description: This service does blablabla
`

	// when
	diagnostic := diag.Diagnostics{}
	data, err := GetOpenSloData(yamlSpec, &diagnostic)

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	if len(data.Services) != 0 {
		t.Errorf("Expected 0 services, but got %d", len(data.Services))
	}

	// and
	if len(diagnostic.Warnings()) != 1 {
		t.Errorf("Expected 1 warning, but got %d", len(diagnostic.Warnings()))
	}

	// and
	if !strings.Contains(diagnostic.Warnings()[0].Summary(), "Unsupported apiVersion") {
		t.Errorf("Expected 'Unsupported apiVersion', but got %s", diagnostic.Warnings()[0].Summary())
	}
}

func TestOpenSLO_shouldbeError_badKind(t *testing.T) {
	// given
	yamlSpec := `
apiVersion: openslo/v1
kind: Unsupported
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
`

	// when
	diagnostics := diag.Diagnostics{}
	_, err := GetOpenSloData(yamlSpec, &diagnostics)

	// then
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	// and
	if !strings.Contains(err.Error(), "Unknown kind") {
		t.Errorf("Expected 'Unknown kind', but got %s", err.Error())
	}

	// and
	if len(diagnostics.Errors()) != 1 {
		t.Errorf("Expected 1 error, but got %d", len(diagnostics.Errors()))
	}

	// and
	if !strings.Contains(diagnostics.Errors()[0].Summary(), "Unsupported kind") {
		t.Errorf("Expected 'Unsupported kind', but got %s", diagnostics.Errors()[0].Summary())
	}
}

func TestOpenSLOAlertPolicy_shouldbeError_badRef(t *testing.T) {
	// given
	yamlSpec1 := `
apiVersion: openslo/v1
kind: AlertPolicy
metadata:
  name: default
  displayName: Alert Policy
spec:
  conditions:
  - conditionRef: string
`
	yamlSpec2 := `
apiVersion: openslo/v1
kind: AlertPolicy
metadata:
    name: default
    displayName: Alert Policy
spec:
    notificationTargets:
    - targetRef: OnCallDevopsMailNotification
`

	// when
	diagnostics := diag.Diagnostics{}
	_, err1 := GetOpenSloData(yamlSpec1, &diagnostics)
	_, err2 := GetOpenSloData(yamlSpec2, &diagnostics)

	// then
	if err1 == nil {
		t.Error("Expected error for bad conditionRef, but got nil")
	}

	// and
	if err2 == nil {
		t.Error("Expected error for bad notificationTarget.targetRef, but got nil")
	}

	// and
	if len(diagnostics.Errors()) != 2 {
		t.Errorf("Expected 2 errors, but got %d", len(diagnostics.Errors()))
	}

	// and
	if !strings.Contains(diagnostics.Errors()[0].Summary(), "Bad reference") || !strings.Contains(diagnostics.Errors()[1].Summary(), "Bad reference") {
		t.Errorf("Expected 'Bad reference', but got %s", diagnostics.Errors()[0].Summary())
	}
}

func TestOpenSLOSLI_shouldbeError_badRef(t *testing.T) {

	yamlSpec_threshold := `
apiVersion: openslo/v1
kind: SLI
metadata:
  name: default-success-rate
  displayName: string
spec:
  description: string
  thresholdMetric:
    metricSource:
        metricSourceRef: default1
`

	yamlSpec_bad := `
apiVersion: openslo/v1
kind: SLI
metadata:
  name: default-success-rate
  displayName: string
spec:
  description: string 
  ratioMetric:
    bad:
      metricSource:
        metricSourceRef: default2
`

	yamlSpec_good := `
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
        metricSourceRef: default3
`

	yamlSpec_total := `
apiVersion: openslo/v1
kind: SLI
metadata:
  name: default-success-rate
  displayName: string
spec:
  description: string 
  ratioMetric:
    total:
      metricSource:
        metricSourceRef: default4
`

	yamlSpec_raw := `
apiVersion: openslo/v1
kind: SLI
metadata:
  name: default-success-rate
  displayName: string
spec:
  description: string 
  ratioMetric:
    raw:
      metricSource:
        metricSourceRef: default5
`

	// when
	diagnostics := diag.Diagnostics{}
	_, err_threshold := GetOpenSloData(yamlSpec_threshold, &diagnostics)
	_, err_bad := GetOpenSloData(yamlSpec_bad, &diagnostics)
	_, err_good := GetOpenSloData(yamlSpec_good, &diagnostics)
	_, err_total := GetOpenSloData(yamlSpec_total, &diagnostics)
	_, err_raw := GetOpenSloData(yamlSpec_raw, &diagnostics)

	// then
	if err_threshold == nil {
		t.Error("Expected error for bad thresholdMetric.metricSourceRef, but got nil")
	}

	if err_bad == nil {
		t.Error("Expected error for bad ratioMetric.bad.metricSourceRef, but got nil")
	}

	if err_good == nil {
		t.Error("Expected error for bad ratioMetric.good.metricSourceRef, but got nil")
	}

	if err_total == nil {
		t.Error("Expected error for bad ratioMetric.total.metricSourceRef, but got nil")
	}

	if err_raw == nil {
		t.Error("Expected error for bad ratioMetric.raw.metricSourceRef, but got nil")
	}

	// and

	if len(diagnostics.Errors()) != 5 {
		t.Errorf("Expected 5 errors, but got %d", len(diagnostics.Errors()))
	}

	// and

	if !strings.Contains(diagnostics.Errors()[0].Summary(), "Bad reference") ||
		!strings.Contains(diagnostics.Errors()[1].Summary(), "Bad reference") ||
		!strings.Contains(diagnostics.Errors()[2].Summary(), "Bad reference") ||
		!strings.Contains(diagnostics.Errors()[3].Summary(), "Bad reference") ||
		!strings.Contains(diagnostics.Errors()[4].Summary(), "Bad reference") {

		t.Errorf("Expected 'Bad reference', but got %s", diagnostics.Errors()[0].Summary())
	}
}
