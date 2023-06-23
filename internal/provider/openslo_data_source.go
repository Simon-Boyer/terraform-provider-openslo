package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const OPENSLO_VERSION = "openslo/v1"
const OPENSLO_EXTENSION_SYNTHETICS = "openslo_synthetics/v1"

func NewOpenSloDataSource() datasource.DataSource {
	return &OpenSloDataSource{}
}

// OpenSloDataSource defines the data source implementation.
type OpenSloDataSource struct {
	Yaml_input                 types.String                            `tfsdk:"yaml_input"`
	Datasources                map[string]DataSourceModel              `tfsdk:"datasources"`
	Services                   map[string]ServiceModel                 `tfsdk:"services"`
	Alert_conditions           map[string]AlertConditionModel          `tfsdk:"alert_conditions"`
	Alert_notification_targets map[string]AlertNotificationTargetModel `tfsdk:"alert_notification_targets"`
	Alert_policies             map[string]AlertPolicyModel             `tfsdk:"alert_policies"`
	Slis                       map[string]SLIModel                     `tfsdk:"slis"`
	Slos                       map[string]SLOModel                     `tfsdk:"slos"`
	Extension_browsermonitor   map[string]BrowserMonitorModel          `tfsdk:"extension_browsermonitor"`
	Extension_httpmonitor      map[string]HTTPMonitorModel             `tfsdk:"extension_httpmonitor"`
}

func (d *OpenSloDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openslo"
}

func (d *OpenSloDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// This is the actual schema definition
	resp.Schema = schema.Schema{
		MarkdownDescription: "OpenSlo data source. Please go to https://github.com/OpenSLO/OpenSLO for field definitions",
		Attributes: map[string]schema.Attribute{
			"yaml_input": schema.StringAttribute{
				MarkdownDescription: "OpenSLO yaml content input",
				Optional:            false,
				Required:            true,
			},
			"datasources": schema.MapAttribute{
				MarkdownDescription: "Datasources",
				Computed:            true,
				ElementType:         DataSourceSchema,
			},
			"services": schema.MapAttribute{
				MarkdownDescription: "Services",
				Computed:            true,
				ElementType:         ServiceSchema,
			},
			"alert_conditions": schema.MapAttribute{
				MarkdownDescription: "Alert conditions",
				Computed:            true,
				ElementType:         AlertConditionSchema,
			},
			"alert_notification_targets": schema.MapAttribute{
				MarkdownDescription: "Alert notification targets",
				Computed:            true,
				ElementType:         AlertNotificationTargetSchema,
			},
			"alert_policies": schema.MapAttribute{
				MarkdownDescription: "Alert policies",
				Computed:            true,
				ElementType:         AlertPolicySchema,
			},
			"slis": schema.MapAttribute{
				MarkdownDescription: "SLIs",
				Computed:            true,
				ElementType:         SLISchema,
			},
			"slos": schema.MapAttribute{
				MarkdownDescription: "SLOs",
				Computed:            true,
				ElementType:         SLOSchema,
			},
			"extension_httpmonitor": schema.MapAttribute{
				MarkdownDescription: "Synthetics HTTP (extension)",
				Computed:            true,
				ElementType:         HTTPMonitorSchema,
			},
			"extension_browsermonitor": schema.MapAttribute{
				MarkdownDescription: "Synthetics Browser (extension)",
				Computed:            true,
				ElementType:         BrowserMonitorSchema,
			},
		},
	}
}

func (d *OpenSloDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *OpenSloDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var readData OpenSloDataSource

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &readData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := d.GetOpenSloData(readData.Yaml_input.ValueString(), &resp.Diagnostics)
	if err != nil {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &d)...)
}

func (d *OpenSloDataSource) GetOpenSloData(yamlInput string, diagnostics *diag.Diagnostics) error {

	d.Yaml_input = types.StringValue(yamlInput)
	d.Datasources = map[string]DataSourceModel{}
	d.Services = map[string]ServiceModel{}
	d.Slis = map[string]SLIModel{}
	d.Slos = map[string]SLOModel{}
	d.Alert_conditions = map[string]AlertConditionModel{}
	d.Alert_notification_targets = map[string]AlertNotificationTargetModel{}
	d.Alert_policies = map[string]AlertPolicyModel{}
	d.Extension_browsermonitor = map[string]BrowserMonitorModel{}
	d.Extension_httpmonitor = map[string]HTTPMonitorModel{}

	// We decode the yaml with 2 decoder iterators, so we can get the kind then unmarshal the yaml value

	yamlBytes := []byte(d.Yaml_input.ValueString())
	decKind := yaml.NewDecoder(bytes.NewReader(yamlBytes))
	decType := yaml.NewDecoder(bytes.NewReader(yamlBytes))

	for {
		var doc YamlSpec
		err := decKind.Decode(&doc)

		// Break out of the loop if error or EOF
		if err != nil {
			if err != io.EOF {
				diagnostics.AddError("Failed to decode yaml", err.Error())
				return err
			}
			break
		}

		// Then we can unmarshal based on the kind
		switch doc.ApiVersion {
		case OPENSLO_VERSION:
			err = d.ExtractOpenSloDocument(&doc, decType)
		case OPENSLO_EXTENSION_SYNTHETICS:
			err = d.ExtractSyntheticsExtensionDocument(&doc, decType)
		default:
			diagnostics.AddWarning("Unsupported apiVersion, skipping", fmt.Sprintf("Expected %s, got %s", OPENSLO_VERSION, doc.ApiVersion))
		}

		if err != nil {
			diagnostics.AddError("Decode Error", err.Error())
			return err
		}
	}

	err := d.OpenSloPostExtractionLogic()
	if err != nil {
		diagnostics.AddError("OpenSLO Post Extraction Error", err.Error())
		return err
	}

	err = d.SyntheticsExtensionPostExtractionLogic()
	if err != nil {
		diagnostics.AddError("Synthetics Extension Post Extraction Error", err.Error())
		return err
	}

	return nil
}
