package provider

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
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
	}
	// when
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})
	// then
	if err != nil {
		t.Error(err)
	}
	// and
	diff := deep.Equal(openslo.Datasources["my-datasource"], expected)
	if diff != nil {
		t.Error(diff)
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
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})
	// then
	if err != nil {
		t.Error(err)
	}
	// and
	diff := deep.Equal(openslo.Services["my-service"], expected)
	if diff != nil {
		t.Error(diff)
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
			MetricSource: MetricSource{
				Type: "string",
				Spec: map[string]interface{}{
					"query": "a_query",
				},
			},
		},
		RatioMetric: RatioMetricModel{
			Counter: true,
			Good: MetricModel{
				MetricSource: MetricSource{
					Type: "string",
					Spec: map[string]interface{}{
						"query": "a_query",
					},
				},
			},
			Bad: MetricModel{
				MetricSource: MetricSource{
					Type: "string",
					Spec: map[string]interface{}{
						"query": "a_query",
					},
				},
			},
			Total: MetricModel{
				MetricSource: MetricSource{
					Type: "string",
					Spec: map[string]interface{}{
						"query": "a_query",
					},
				},
			},
			RawType: "success",
			Raw: MetricModel{
				MetricSource: MetricSource{
					Type: "string",
					Spec: map[string]interface{}{
						"query": "a_query",
					},
				},
			},
		},
	}

	// when
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	diff := deep.Equal(openslo.Slis["string"], expected)
	if diff != nil {
		t.Error(diff)
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
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})
	// then
	if err != nil {
		t.Error(err)
	}
	// and
	diff := deep.Equal(openslo.Alert_policies["default"], expected)
	if diff != nil {
		t.Error(diff)
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
				MetricSource: MetricSource{
					Type: "string",
					Spec: map[string]interface{}{
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
						MetricSource: MetricSource{
							Type: "string",
							Spec: map[string]interface{}{
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
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	diff := deep.Equal(openslo.Slos["string"], expected)
	if diff != nil {
		t.Error(diff)
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
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	diff := deep.Equal(openslo.Alert_conditions["string"], expected)
	if diff != nil {
		t.Error(diff)
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
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	diff := deep.Equal(openslo.Alert_notification_targets["OnCallDevopsMailNotification"], expected)
	if diff != nil {
		t.Error(diff)
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
        spec:
          query: "sum:api.requests{200}.as_count()"
    total:
      metricSource:
        type: datadog
        spec:
          query: "sum:api.requests{*}.as_count()"
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
	}

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
				MetricSource: MetricSource{
					MetricSourceRef: "default",
					Type:            "datadog",
					Spec: map[string]interface{}{
						"query": "sum:api.requests{200}.as_count()",
					},
					DataSource: dataSource,
				},
			},
			Total: MetricModel{
				MetricSource: MetricSource{
					Type: "datadog",
					Spec: map[string]interface{}{
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
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	diff := deep.Equal(openslo.Services["my-service"], service)
	if diff != nil {
		t.Error(diff)
	}

	// and
	diff = deep.Equal(openslo.Datasources["default"], dataSource)
	if diff != nil {
		t.Error(diff)
	}

	// and
	diff = deep.Equal(openslo.Alert_conditions["string"], alertCondition)
	if diff != nil {
		t.Error(diff)
	}

	// and
	diff = deep.Equal(openslo.Alert_notification_targets["OnCallDevopsMailNotification"], alertNotificationTarget)
	if diff != nil {
		t.Error(diff)
	}

	// and
	diff = deep.Equal(openslo.Alert_policies["default"], alertPolicy)
	if diff != nil {
		t.Error(diff)
	}

	// and
	diff = deep.Equal(openslo.Slis["default-success-rate"], sli)
	if diff != nil {
		t.Error(diff)
	}

	// and
	diff = deep.Equal(openslo.Slos["string"], slo)
	if diff != nil {
		t.Error(diff)
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
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diagnostic)

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	if len(openslo.Services) != 0 {
		t.Errorf("Expected 0 services, but got %d", len(openslo.Services))
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
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diagnostics)

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
	if !strings.Contains(diagnostics.Errors()[0].Summary(), "Decode Error") {
		t.Errorf("Expected 'Decode Error', but got %s", diagnostics.Errors()[0].Summary())
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
	openslo1 := OpenSloDataSource{}
	openslo2 := OpenSloDataSource{}
	err1 := openslo1.GetOpenSloData(yamlSpec1, &diagnostics)
	err2 := openslo2.GetOpenSloData(yamlSpec2, &diagnostics)

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
	if !strings.Contains(diagnostics.Errors()[0].Summary(), "OpenSLO Post Extraction Error") {
		t.Errorf("Expected 'OpenSLO Post Extraction Error', but got %s", diagnostics.Errors()[0].Summary())
	}

	if !strings.Contains(diagnostics.Errors()[1].Detail(), "bad reference") {
		t.Errorf("Expected 'Bad reference', but got %s", diagnostics.Errors()[1].Summary())
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
	openslo_threshold := OpenSloDataSource{}
	openslo_bad := OpenSloDataSource{}
	openslo_good := OpenSloDataSource{}
	openslo_total := OpenSloDataSource{}
	openslo_raw := OpenSloDataSource{}
	err_threshold := openslo_threshold.GetOpenSloData(yamlSpec_threshold, &diagnostics)
	err_bad := openslo_bad.GetOpenSloData(yamlSpec_bad, &diagnostics)
	err_good := openslo_good.GetOpenSloData(yamlSpec_good, &diagnostics)
	err_total := openslo_total.GetOpenSloData(yamlSpec_total, &diagnostics)
	err_raw := openslo_raw.GetOpenSloData(yamlSpec_raw, &diagnostics)

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

	if !strings.Contains(diagnostics.Errors()[0].Summary(), "OpenSLO Post Extraction Error") ||
		!strings.Contains(diagnostics.Errors()[1].Summary(), "OpenSLO Post Extraction Error") ||
		!strings.Contains(diagnostics.Errors()[2].Summary(), "OpenSLO Post Extraction Error") ||
		!strings.Contains(diagnostics.Errors()[3].Summary(), "OpenSLO Post Extraction Error") ||
		!strings.Contains(diagnostics.Errors()[4].Summary(), "OpenSLO Post Extraction Error") {

		t.Errorf("Expected 'OpenSLO Post Extraction Error', but got %s", diagnostics.Errors()[0].Summary())
	}
}
