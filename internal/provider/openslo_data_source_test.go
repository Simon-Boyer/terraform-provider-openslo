package provider

import (
	"reflect"
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
	openslo, err := GetOpenSloData(yamlSpec, diag.Diagnostics{})
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
	openslo, err := GetOpenSloData(yamlSpec, diag.Diagnostics{})
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
	openslo, err := GetOpenSloData(yamlSpec, diag.Diagnostics{})

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
	openslo, err := GetOpenSloData(yamlSpec, diag.Diagnostics{})
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
	openslo, err := GetOpenSloData(yamlSpec, diag.Diagnostics{})

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
	openslo, err := GetOpenSloData(yamlSpec, diag.Diagnostics{})

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
	openslo, err := GetOpenSloData(yamlSpec, diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	if !reflect.DeepEqual(openslo.Alert_notification_targets["OnCallDevopsMailNotification"], expected) {
		t.Errorf("Expected %#v, but got %#v", expected, openslo.Alert_notification_targets["OnCallDevopsMailNotification"])
	}
}