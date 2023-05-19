package provider

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"
)

const OPENSLO_VERSION = "openslo/v1"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OpenSloDataSource{}

func NewOpenSloDataSource() datasource.DataSource {
	return &OpenSloDataSource{}
}

// OpenSloDataSource defines the data source implementation.
type OpenSloDataSource struct {
	client *http.Client
}

// OpenSloDataSourceModel describes the data source data model.
type OpenSloDataSourceModel struct {
	yaml_input                 types.String                                  `tfsdk:"yaml_input"`
	id                         types.String                                  `tfsdk:"id"`
	datasources                map[types.String]DataSourceModel              `tfsdk:"datasources"`
	services                   map[types.String]ServiceModel                 `tfsdk:"services"`
	alert_conditions           map[types.String]AlertConditionModel          `tfsdk:"alert_conditions"`
	alert_notification_targets map[types.String]AlertNotificationTargetModel `tfsdk:"alert_notification_targets"`
	alert_policies             map[types.String]AlertPolicyModel             `tfsdk:"alert_policies"`
	slis                       map[types.String]SLIModel                     `tfsdk:"slis"`
	slos                       map[types.String]SLOModel                     `tfsdk:"slos"`
}

type YamlSpec struct {
	kind       string
	apiVersion string
	metadata   MetadataModel
}

type YamlSpecTyped[T any] struct {
	kind       string
	apiVersion string
	metadata   MetadataModel
	spec       T
}

func (d *OpenSloDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_OpenSlo"
}

func (d *OpenSloDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "OpenSlo data source",

		Attributes: map[string]schema.Attribute{
			"yaml_input": schema.StringAttribute{
				MarkdownDescription: "OpenSLO yaml content input",
				Optional:            false,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "OpenSlo identifier",
				Computed:            true,
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
		},
	}
}

func (d *OpenSloDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *OpenSloDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OpenSloDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read OpenSlo, got error: %s", err))
	//     return
	// }

	// For the purposes of this OpenSlo code, hardcoding a response value to
	// save into the Terraform state.

	decKind := yaml.NewDecoder(bytes.NewReader([]byte(data.yaml_input.ValueString())))
	decType := yaml.NewDecoder(bytes.NewReader([]byte(data.yaml_input.ValueString())))

	for {
		var doc YamlSpec
		if decKind.Decode(&doc) != nil {
			if doc.apiVersion == OPENSLO_VERSION {
				switch doc.kind {
				case "DataSource":
					var typedDoc YamlSpecTyped[DataSourceModel]
					decType.Decode(&typedDoc)
					typedDoc.spec.Metadata = doc.metadata
					data.datasources[doc.metadata.Name] = typedDoc.spec
				case "Service":
					var typedDoc YamlSpecTyped[ServiceModel]
					decType.Decode(&typedDoc)
					typedDoc.spec.Metadata = doc.metadata
					data.services[doc.metadata.Name] = typedDoc.spec
				case "AlertCondition":
					var typedDoc YamlSpecTyped[AlertConditionModel]
					decType.Decode(&typedDoc)
					typedDoc.spec.Metadata = doc.metadata
					data.alert_conditions[doc.metadata.Name] = typedDoc.spec
				case "AlertNotificationTarget":
					var typedDoc YamlSpecTyped[AlertNotificationTargetModel]
					decType.Decode(&typedDoc)
					typedDoc.spec.Metadata = doc.metadata
					data.alert_notification_targets[doc.metadata.Name] = typedDoc.spec
				case "AlertPolicy":
					var typedDoc YamlSpecTyped[AlertPolicyModel]
					decType.Decode(&typedDoc)
					typedDoc.spec.Metadata = doc.metadata
					data.alert_policies[doc.metadata.Name] = typedDoc.spec
				case "SLI":
					var typedDoc YamlSpecTyped[SLIModel]
					decType.Decode(&typedDoc)
					typedDoc.spec.Metadata = doc.metadata
					data.slis[doc.metadata.Name] = typedDoc.spec
				case "SLO":
					var typedDoc YamlSpecTyped[SLOModel]
					decType.Decode(&typedDoc)
					typedDoc.spec.Metadata = doc.metadata
					data.slos[doc.metadata.Name] = typedDoc.spec
				default:
					resp.Diagnostics.AddError(
						"Unexpected Kind", "Unknown kind: "+doc.kind,
					)
				}
			}
			break
		}

		// Embed referenced objects for alert policies
		for i := range data.alert_policies {
			for j := range data.alert_policies[i].Conditions {
				condition := data.alert_policies[i].Conditions[j]
				if !condition.ConditionRef.IsNull() {
					linkedCond := data.alert_conditions[condition.ConditionRef]
					linkedCond.ConditionRef = condition.ConditionRef
					data.alert_policies[i].Conditions[j] = linkedCond
				}
			}
			for j := range data.alert_policies[i].NotificationTargets {
				condition := data.alert_policies[i].NotificationTargets[j]
				if !condition.TargetRef.IsNull() {
					linkedCond := data.alert_notification_targets[condition.TargetRef]
					linkedCond.TargetRef = condition.TargetRef
					data.alert_policies[i].NotificationTargets[j] = linkedCond
				}
			}
		}

		// Embed referenced objects for slos
		for k := range data.slos {
			slo := data.slos[k]
			if !slo.IndicatorRef.IsNull() {
				slo.Indicator = data.slis[slo.IndicatorRef]
			}
			for j := range slo.AlertPolicies {
				alertPolicy := slo.AlertPolicies[j]
				if !alertPolicy.AlertPolicyRef.IsNull() {
					linkedAlertPolicy := data.alert_policies[alertPolicy.AlertPolicyRef]
					linkedAlertPolicy.AlertPolicyRef = alertPolicy.AlertPolicyRef
					slo.AlertPolicies[j] = linkedAlertPolicy
				}
			}
			data.slos[k] = slo
		}

	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
