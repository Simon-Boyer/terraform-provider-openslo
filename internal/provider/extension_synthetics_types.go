package provider

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
	PATCH   HTTPMethod = "PATCH"
)

type HTTPMonitorModel struct {
	Metadata   MetadataModel  `tfsdk:"metadata" yaml:"metadata"`
	Url        string         `tfsdk:"url" yaml:"url"`
	Requests   []RequestModel `tfsdk:"requests" yaml:"requests"`
	Service    ServiceModel   `tfsdk:"service" yaml:"-"`
	ServiceRef string         `tfsdk:"service_ref" yaml:"serviceRef"`
}

type RequestModel struct {
	Name             string        `tfsdk:"name" yaml:"name"`
	Description      string        `tfsdk:"description" yaml:"description"`
	Headers          []HeaderModel `tfsdk:"headers" yaml:"headers"`
	Body             string        `tfsdk:"body" yaml:"body"`
	Method           HTTPMethod    `tfsdk:"method" yaml:"method"`
	Path             string        `tfsdk:"path" yaml:"path"`
	ExpectedResponse ResponseModel `tfsdk:"expected_response" yaml:"expectedResponse"`
}

type ResponseModel struct {
	Codes                   []int  `tfsdk:"codes" yaml:"code"`
	PayloadContains         string `tfsdk:"payload_contains" yaml:"payloadContains"`
	PayloadNotContains      string `tfsdk:"payload_not_contains" yaml:"payloadNotContains"`
	DynatracePostProcessing string `tfdsk:"dt_postprocessing" yaml:"dynatrace_postprocessing"`
}

type HeaderModel struct {
	Name  string `tfsdk:"name" yaml:"name"`
	Value string `tfsdk:"value" yaml:"value"`
}

type BrowserMonitorModel struct {
	Metadata   MetadataModel `tfsdk:"metadata" yaml:"metadata"`
	Url        string        `tfsdk:"url" yaml:"url"`
	Script     string        `tfsdk:"script" yaml:"script"`
	Service    ServiceModel  `tfsdk:"service" yaml:"-"`
	ServiceRef string        `tfsdk:"service_ref" yaml:"serviceRef"`
}
