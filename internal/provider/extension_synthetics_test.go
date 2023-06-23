package provider

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestSynthetics_shouldbeValid_httpMonitorYamlSpec(t *testing.T) {
	// given
	yamlSpec := `
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
`

	expected := HTTPMonitorModel{
		Metadata: MetadataModel{
			DisplayName: "My Monitor",
			Name:        "my-monitor",
		},
		Url: "https://my-host.com",
		Requests: []RequestModel{
			{
				Name:        "my-request",
				Description: "This is a request",
				Headers: []HeaderModel{
					{
						Name:  "my-header",
						Value: "my-value",
					},
				},
				Body:   "test body",
				Method: "POST",
				Path:   "/my-path",
			},
			{
				Name:        "my-other-request",
				Description: "This is a request",
				Headers: []HeaderModel{
					{
						Name:  "my-other-header",
						Value: "my-other-value",
					},
				},
				Method: "GET",
				Path:   "/my-other-path",
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
	diff := deep.Equal(openslo.Extension_httpmonitor["my-monitor"], expected)
	if diff != nil {
		t.Error(diff)
	}
}

func TestSynthetics_shouldbeValid_browserMonitorYamlSpec(t *testing.T) {
	// given
	yamlSpec := `
    apiVersion: openslo_synthetics/v1
    kind: BrowserMonitor
    metadata:
        name: my-monitor
        displayName: My Monitor
    spec:
        url: https://my-host.com
        script: |
            console.log("hello")
            abc
`

	expected := BrowserMonitorModel{
		Metadata: MetadataModel{
			DisplayName: "My Monitor",
			Name:        "my-monitor",
		},
		Url:    "https://my-host.com",
		Script: "console.log(\"hello\")\nabc\n",
	}

	// when
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	diff := deep.Equal(openslo.Extension_browsermonitor["my-monitor"], expected)
	if diff != nil {
		t.Error(diff)
	}
}

func TestSynthetics_shouldbeValid_multiYamlSpecWithRefs(t *testing.T) {
	// given
	yamlSpec := `
    apiVersion: openslo/v1
    kind: Service
    metadata:
        name: my-service
        displayName: My Service
    spec:
        description: This service does blablabla
---
    apiVersion: openslo_synthetics/v1
    kind: HTTPMonitor
    metadata:
        name: my-monitor
        displayName: My Monitor
    spec:
        serviceRef: my-service
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
---
    apiVersion: openslo_synthetics/v1
    kind: BrowserMonitor
    metadata:
        name: my-monitor
        displayName: My Monitor
    spec:
        serviceRef: my-service
        url: https://my-host.com
        script: |
            console.log("hello")
            abc
`
	expectedService := ServiceModel{
		Metadata: MetadataModel{
			Name:        "my-service",
			DisplayName: "My Service",
		},
		Description: "This service does blablabla",
	}

	expectedBrowser := BrowserMonitorModel{
		Metadata: MetadataModel{
			DisplayName: "My Monitor",
			Name:        "my-monitor",
		},
		Url:        "https://my-host.com",
		Script:     "console.log(\"hello\")\nabc\n",
		ServiceRef: "my-service",
		Service:    expectedService,
	}

	expectedHttp := HTTPMonitorModel{
		Metadata: MetadataModel{
			DisplayName: "My Monitor",
			Name:        "my-monitor",
		},
		Url: "https://my-host.com",
		Requests: []RequestModel{
			{
				Name:        "my-request",
				Description: "This is a request",
				Headers: []HeaderModel{
					{
						Name:  "my-header",
						Value: "my-value",
					},
				},
				Body:   "test body",
				Method: "POST",
				Path:   "/my-path",
			},
			{
				Name:        "my-other-request",
				Description: "This is a request",
				Headers: []HeaderModel{
					{
						Name:  "my-other-header",
						Value: "my-other-value",
					},
				},
				Method: "GET",
				Path:   "/my-other-path",
			},
		},
		ServiceRef: "my-service",
		Service:    expectedService,
	}

	// when
	openslo := OpenSloDataSource{}
	err := openslo.GetOpenSloData(yamlSpec, &diag.Diagnostics{})

	// then
	if err != nil {
		t.Error(err)
	}

	// and
	diff := deep.Equal(openslo.Extension_browsermonitor["my-monitor"], expectedBrowser)
	if diff != nil {
		t.Error(diff)
	}

	// and
	diff = deep.Equal(openslo.Extension_httpmonitor["my-monitor"], expectedHttp)
	if diff != nil {
		t.Error(diff)
	}
}
