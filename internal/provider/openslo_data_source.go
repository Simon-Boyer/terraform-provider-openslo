package provider

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Yaml_input                 types.String                            `tfsdk:"yaml_input"`
	Datasources                map[string]DataSourceModel              `tfsdk:"datasources"`
	Services                   map[string]ServiceModel                 `tfsdk:"services"`
	Alert_conditions           map[string]AlertConditionModel          `tfsdk:"alert_conditions"`
	Alert_notification_targets map[string]AlertNotificationTargetModel `tfsdk:"alert_notification_targets"`
	Alert_policies             map[string]AlertPolicyModel             `tfsdk:"alert_policies"`
	Slis                       map[string]SLIModel                     `tfsdk:"slis"`
	Slos                       map[string]SLOModel                     `tfsdk:"slos"`
}

type Test struct {
	Kind       string
	ApiVersion string
	Metadata   MetadataModel
	Spec       SLIModel
}

type YamlSpec struct {
	Kind       string        `yaml:"kind"`
	ApiVersion string        `yaml:"apiVersion"`
	Metadata   MetadataModel `yaml:"metadata"`
}

type YamlSpecTyped[T any] struct {
	Kind       string
	ApiVersion string
	Metadata   MetadataModel
	Spec       T
}

func (d *OpenSloDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openslo"
}

func (d *OpenSloDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// We define schemas as variables so it can be reused as a nested schema
	// The order here is important since a schema may refer to other schemas

	var MetadataSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":         types.StringType, // test
			"display_name": types.StringType,
		},
	}

	var DataSourceSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description": types.StringType,
			"connection_details": types.MapType{
				ElemType: types.StringType,
			},
			"metadata":          MetadataSchema,
			"metric_source_ref": types.StringType,
			"spec": types.MapType{
				ElemType: types.StringType,
			},
			"type": types.StringType,
		},
	}

	var ServiceSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description": types.StringType,
			"metadata":    MetadataSchema,
		},
	}

	var AlertConditionModelConditionSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"op":              types.StringType,
			"threshold":       types.NumberType,
			"lookback_window": types.StringType,
			"alert_after":     types.StringType,
		},
	}

	var AlertConditionSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description":   types.StringType,
			"severity":      types.StringType,
			"condition":     AlertConditionModelConditionSchema,
			"metadata":      MetadataSchema,
			"condition_ref": types.StringType,
		},
	}

	var AlertNotificationTargetSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description": types.StringType,
			"target":      types.StringType,
			"metadata":    MetadataSchema,
			"target_ref":  types.StringType,
		},
	}

	var MetricSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"metric_source": DataSourceSchema,
		},
	}

	var RatioMetricSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"counter":  types.BoolType,
			"good":     MetricSchema,
			"bad":      MetricSchema,
			"total":    MetricSchema,
			"raw_type": types.StringType,
			"raw":      MetricSchema,
		},
	}

	var SLISchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description":      types.StringType,
			"threshold_metric": MetricSchema,
			"ratio_metric":     RatioMetricSchema,
			"metadata":         MetadataSchema,
		},
	}

	var AlertPolicySchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description":          types.StringType,
			"alert_when_no_data":   types.BoolType,
			"alert_when_resolved":  types.BoolType,
			"alert_when_breaching": types.BoolType,
			"conditions": types.ListType{
				ElemType: AlertConditionSchema,
			},
			"notification_targets": types.ListType{
				ElemType: AlertNotificationTargetSchema,
			},
			"metadata":         MetadataSchema,
			"alert_policy_ref": types.StringType,
		},
	}

	var CalendarSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"start_time": types.StringType,
			"time_zone":  types.StringType,
		},
	}

	var TimeWindowSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"duration":   types.StringType,
			"calendar":   CalendarSchema,
			"is_rolling": types.BoolType,
		},
	}

	var ObjectiveSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"display_name":      types.StringType,
			"op":                types.StringType,
			"value":             types.NumberType,
			"target":            types.NumberType,
			"target_percentage": types.NumberType,
			"time_slice_target": types.NumberType,
			"time_slice_window": types.NumberType,
			"indicator":         SLISchema,
			"composite_weight":  types.NumberType,
			"indicator_ref":     types.StringType,
		},
	}

	var SLOSchema = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description":   types.StringType,
			"service":       ServiceSchema,
			"service_ref":   types.StringType,
			"indicator":     SLISchema,
			"indicator_ref": types.StringType,
			"time_window": types.ListType{
				ElemType: TimeWindowSchema,
			},
			"budgeting_method": types.StringType,
			"objectives": types.ListType{
				ElemType: ObjectiveSchema,
			},
			"alert_policies": types.ListType{
				ElemType: AlertPolicySchema,
			},
			"metadata": MetadataSchema,
		},
	}

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
	var readData OpenSloDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &readData)...)
	data := OpenSloDataSourceModel{
		Yaml_input:                 readData.Yaml_input,
		Datasources:                map[string]DataSourceModel{},
		Services:                   map[string]ServiceModel{},
		Slis:                       map[string]SLIModel{},
		Slos:                       map[string]SLOModel{},
		Alert_conditions:           map[string]AlertConditionModel{},
		Alert_notification_targets: map[string]AlertNotificationTargetModel{},
		Alert_policies:             map[string]AlertPolicyModel{},
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// We decode the yaml with 2 decoder iterators, so we can get the kind then unmarshal the yaml value

	decKind := yaml.NewDecoder(bytes.NewReader([]byte(data.Yaml_input.ValueString())))
	decType := yaml.NewDecoder(bytes.NewReader([]byte(data.Yaml_input.ValueString())))

	for {
		var doc YamlSpec
		err := decKind.Decode(&doc)
		if err == nil {
			// Make sure we are dealing with an openslo document
			if doc.ApiVersion == OPENSLO_VERSION {
				var errType error
				// Then we can unmarshal based on the kind
				switch doc.Kind {
				case "DataSource":
					var typedDoc YamlSpecTyped[DataSourceModel]
					errType = decType.Decode(&typedDoc)
					typedDoc.Spec.Metadata = doc.Metadata
					data.Datasources[doc.Metadata.Name] = typedDoc.Spec
				case "Service":
					var typedDoc YamlSpecTyped[ServiceModel]
					errType = decType.Decode(&typedDoc)
					typedDoc.Spec.Metadata = doc.Metadata
					data.Services[doc.Metadata.Name] = typedDoc.Spec
				case "AlertCondition":
					var typedDoc YamlSpecTyped[AlertConditionModel]
					errType = decType.Decode(&typedDoc)
					typedDoc.Spec.Metadata = doc.Metadata
					data.Alert_conditions[doc.Metadata.Name] = typedDoc.Spec
				case "AlertNotificationTarget":
					var typedDoc YamlSpecTyped[AlertNotificationTargetModel]
					errType = decType.Decode(&typedDoc)
					typedDoc.Spec.Metadata = doc.Metadata
					data.Alert_notification_targets[doc.Metadata.Name] = typedDoc.Spec
				case "AlertPolicy":
					var typedDoc YamlSpecTyped[AlertPolicyModel]
					errType = decType.Decode(&typedDoc)
					typedDoc.Spec.Metadata = doc.Metadata
					data.Alert_policies[doc.Metadata.Name] = typedDoc.Spec
				case "SLI":
					var typedDoc YamlSpecTyped[SLIModel]
					errType = decType.Decode(&typedDoc)
					typedDoc.Spec.Metadata = doc.Metadata
					data.Slis[doc.Metadata.Name] = typedDoc.Spec
				case "SLO":
					var typedDoc YamlSpecTyped[SLOModel]
					errType = decType.Decode(&typedDoc)
					typedDoc.Spec.Metadata = doc.Metadata
					data.Slos[doc.Metadata.Name] = typedDoc.Spec
				default:
					resp.Diagnostics.AddError(
						"Unexpected Kind", "Unknown kind: "+doc.Kind,
					)
				}
				if errType != nil {
					resp.Diagnostics.AddError("Decode Error", errType.Error())
				}
			}
		} else {
			break
		}
	}

	// Embed referenced objects for alert policies
	for i := range data.Alert_policies {
		for j := range data.Alert_policies[i].Conditions {
			condition := data.Alert_policies[i].Conditions[j]
			if condition.ConditionRef != "" {
				linkedCond := data.Alert_conditions[condition.ConditionRef]
				linkedCond.ConditionRef = condition.ConditionRef
				data.Alert_policies[i].Conditions[j] = linkedCond
			}
		}
		for j := range data.Alert_policies[i].NotificationTargets {
			condition := data.Alert_policies[i].NotificationTargets[j]
			if condition.TargetRef != "" {
				linkedCond := data.Alert_notification_targets[condition.TargetRef]
				linkedCond.TargetRef = condition.TargetRef
				data.Alert_policies[i].NotificationTargets[j] = linkedCond
			}
		}
	}

	// Embed referenced objects for slis
	for k := range data.Slis {
		sli := data.Slis[k]
		if sli.ThresholdMetric.MetricSource.MetricSourceRef != "" {
			sli.ThresholdMetric.MetricSource = data.Datasources[sli.ThresholdMetric.MetricSource.MetricSourceRef]
		}
		if sli.RatioMetric.Bad.MetricSource.MetricSourceRef != "" {
			sli.RatioMetric.Bad.MetricSource = data.Datasources[sli.RatioMetric.Bad.MetricSource.MetricSourceRef]
		}
		if sli.RatioMetric.Good.MetricSource.MetricSourceRef != "" {
			sli.RatioMetric.Good.MetricSource = data.Datasources[sli.RatioMetric.Good.MetricSource.MetricSourceRef]
		}
		if sli.RatioMetric.Raw.MetricSource.MetricSourceRef != "" {
			sli.RatioMetric.Raw.MetricSource = data.Datasources[sli.RatioMetric.Raw.MetricSource.MetricSourceRef]
		}
		data.Slis[k] = sli
	}

	// Embed referenced objects for slos
	for k := range data.Slos {
		slo := data.Slos[k]
		if slo.IndicatorRef != "" {
			slo.Indicator = data.Slis[slo.IndicatorRef]
		}
		if slo.ServiceRef != "" {
			slo.Service = data.Services[slo.ServiceRef]
		}
		for j := range slo.AlertPolicies {
			alertPolicy := slo.AlertPolicies[j]
			if alertPolicy.AlertPolicyRef != "" {
				linkedAlertPolicy := data.Alert_policies[alertPolicy.AlertPolicyRef]
				linkedAlertPolicy.AlertPolicyRef = alertPolicy.AlertPolicyRef
				slo.AlertPolicies[j] = linkedAlertPolicy
			}
		}
		for j := range slo.Objectives {
			objective := slo.Objectives[j]
			if objective.IndicatorRef != "" {
				objective.Indicator = data.Slis[objective.IndicatorRef]
				slo.Objectives[j] = objective
			}
		}
		data.Slos[k] = slo
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
