package provider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const OPENSLO_VERSION = "openslo/v1"

func NewOpenSloDataSource() datasource.DataSource {
	return &OpenSloDataSource{}
}

// OpenSloDataSource defines the data source implementation.
type OpenSloDataSource struct {
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
			"kind":            types.StringType,
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
			"time_slice_window": types.StringType,
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
}

func (d *OpenSloDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var readData OpenSloDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &readData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := GetOpenSloData(readData.Yaml_input.String(), &resp.Diagnostics)
	if err != nil {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func GetOpenSloData(yamlInput string, diagnostics *diag.Diagnostics) (OpenSloDataSourceModel, error) {
	data := OpenSloDataSourceModel{
		Yaml_input:                 types.StringValue(yamlInput),
		Datasources:                map[string]DataSourceModel{},
		Services:                   map[string]ServiceModel{},
		Slis:                       map[string]SLIModel{},
		Slos:                       map[string]SLOModel{},
		Alert_conditions:           map[string]AlertConditionModel{},
		Alert_notification_targets: map[string]AlertNotificationTargetModel{},
		Alert_policies:             map[string]AlertPolicyModel{},
	}

	// We decode the yaml with 2 decoder iterators, so we can get the kind then unmarshal the yaml value

	yamlBytes := []byte(data.Yaml_input.ValueString())
	decKind := yaml.NewDecoder(bytes.NewReader(yamlBytes))
	decType := yaml.NewDecoder(bytes.NewReader(yamlBytes))

	for {
		var doc YamlSpec
		err := decKind.Decode(&doc)

		// Break out of the loop if error or EOF
		if err != nil {
			if err != io.EOF {
				diagnostics.AddError("Failed to decode yaml", err.Error())
				return data, err
			}
			break
		}

		// Make sure we are dealing with an openslo document
		if doc.ApiVersion != OPENSLO_VERSION {
			diagnostics.AddWarning("Unsupported apiVersion", fmt.Sprintf("Expected %s, got %s", OPENSLO_VERSION, doc.ApiVersion))
			continue
		}

		// Then we can unmarshal based on the kind
		switch doc.Kind {
		case "DataSource":
			var typedDoc YamlSpecTyped[DataSourceModel]
			err = decType.Decode(&typedDoc)
			typedDoc.Spec.Metadata = doc.Metadata
			data.Datasources[doc.Metadata.Name] = typedDoc.Spec
		case "Service":
			var typedDoc YamlSpecTyped[ServiceModel]
			err = decType.Decode(&typedDoc)
			typedDoc.Spec.Metadata = doc.Metadata
			data.Services[doc.Metadata.Name] = typedDoc.Spec
		case "AlertCondition":
			var typedDoc YamlSpecTyped[AlertConditionModel]
			err = decType.Decode(&typedDoc)
			typedDoc.Spec.Metadata = doc.Metadata
			data.Alert_conditions[doc.Metadata.Name] = typedDoc.Spec
		case "AlertNotificationTarget":
			var typedDoc YamlSpecTyped[AlertNotificationTargetModel]
			err = decType.Decode(&typedDoc)
			typedDoc.Spec.Metadata = doc.Metadata
			data.Alert_notification_targets[doc.Metadata.Name] = typedDoc.Spec
		case "AlertPolicy":
			var typedDoc YamlSpecTyped[AlertPolicyModel]
			err = decType.Decode(&typedDoc)
			typedDoc.Spec.Metadata = doc.Metadata
			for _, cond := range typedDoc.Spec.ConditionsInternal {
				if typedDoc.Spec.ConditionsInternal[0].Kind != "" {
					cond.Spec.Metadata = cond.Metadata
					typedDoc.Spec.Conditions = append(typedDoc.Spec.Conditions, cond.Spec)
				} else {
					typedDoc.Spec.Conditions = append(typedDoc.Spec.Conditions, AlertConditionModel{
						ConditionRef: cond.ConditionRef,
					})
				}
			}
			typedDoc.Spec.ConditionsInternal = nil
			data.Alert_policies[doc.Metadata.Name] = typedDoc.Spec
		case "SLI":
			var typedDoc YamlSpecTyped[SLIModel]
			err = decType.Decode(&typedDoc)
			typedDoc.Spec.Metadata = doc.Metadata
			data.Slis[doc.Metadata.Name] = typedDoc.Spec
		case "SLO":
			var typedDoc YamlSpecTyped[SLOModel]
			err = decType.Decode(&typedDoc)
			typedDoc.Spec.Metadata = doc.Metadata
			if typedDoc.Spec.IndicatorInternal.Kind != "" {
				typedDoc.Spec.Indicator = typedDoc.Spec.IndicatorInternal.Spec
				typedDoc.Spec.Indicator.Metadata = typedDoc.Spec.IndicatorInternal.Metadata
			}
			for _, alertPolicy := range typedDoc.Spec.AlertPoliciesInternal {
				if typedDoc.Spec.AlertPoliciesInternal[0].Kind != "" {
					alertPolicy.Spec.Metadata = alertPolicy.Metadata
					for _, cond := range alertPolicy.Spec.ConditionsInternal {
						if alertPolicy.Spec.ConditionsInternal[0].Kind != "" {
							cond.Spec.Metadata = cond.Metadata
							alertPolicy.Spec.Conditions = append(alertPolicy.Spec.Conditions, cond.Spec)
						} else {
							alertPolicy.Spec.Conditions = append(alertPolicy.Spec.Conditions, AlertConditionModel{
								ConditionRef: cond.ConditionRef,
							})
						}
					}
					alertPolicy.Spec.ConditionsInternal = nil
					typedDoc.Spec.AlertPolicies = append(typedDoc.Spec.AlertPolicies, alertPolicy.Spec)
				} else {
					typedDoc.Spec.AlertPolicies = append(typedDoc.Spec.AlertPolicies, AlertPolicyModel{
						AlertPolicyRef: alertPolicy.AlertPolicyRef,
					})
				}
			}
			for i, objective := range typedDoc.Spec.Objectives {
				if objective.IndicatorInternal.Kind != "" {
					objective.Indicator = objective.IndicatorInternal.Spec
					objective.Indicator.Metadata = objective.IndicatorInternal.Metadata
				}
				objective.IndicatorInternal = YamlSpecTyped[SLIModel]{}
				typedDoc.Spec.Objectives[i] = objective
			}
			typedDoc.Spec.AlertPoliciesInternal = nil
			typedDoc.Spec.IndicatorInternal = YamlSpecTyped[SLIModel]{}
			data.Slos[doc.Metadata.Name] = typedDoc.Spec
		default:
			diagnostics.AddError(
				"Unsupported kind", "Unknown kind: "+doc.Kind,
			)
			return data, errors.New("Unknown kind: " + doc.Kind)
		}

		if err != nil {
			diagnostics.AddError("Decode Error", err.Error())
			return data, err
		}
	}

	// Embed referenced objects for alert policies
	for i := range data.Alert_policies {
		for j := range data.Alert_policies[i].Conditions {
			condition := data.Alert_policies[i].Conditions[j]
			if condition.ConditionRef != "" {
				linkedCond := data.Alert_conditions[condition.ConditionRef]
				if linkedCond.Metadata.Name == "" {
					diagnostics.AddError(
						"Bad reference", fmt.Sprintf("No object of kind %s with name %s", "AlertCondition", condition.ConditionRef),
					)
					return data, errors.New(fmt.Sprintf("Bad reference: No object of kind %s with name %s", "AlertCondition", condition.ConditionRef))
				}
				linkedCond.ConditionRef = condition.ConditionRef
				data.Alert_policies[i].Conditions[j] = linkedCond
			}
		}
		for j := range data.Alert_policies[i].NotificationTargets {
			condition := data.Alert_policies[i].NotificationTargets[j]
			if condition.TargetRef != "" {
				linkedCond := data.Alert_notification_targets[condition.TargetRef]
				if linkedCond.Metadata.Name == "" {
					diagnostics.AddError(
						"Bad reference", fmt.Sprintf("No object of kind %s with name %s", "AlertNotificationTarget", condition.TargetRef),
					)
					return data, errors.New(fmt.Sprintf("Bad reference: No object of kind %s with name %s", "AlertNotificationTarget", condition.TargetRef))
				}
				linkedCond.TargetRef = condition.TargetRef
				data.Alert_policies[i].NotificationTargets[j] = linkedCond
			}
		}
	}

	// Embed referenced objects for slis
	for k := range data.Slis {
		sli := data.Slis[k]
		if sli.ThresholdMetric.MetricSource.MetricSourceRef != "" {
			ref := sli.ThresholdMetric.MetricSource.MetricSourceRef
			sli.ThresholdMetric.MetricSource = data.Datasources[ref]
			if sli.ThresholdMetric.MetricSource.Metadata.Name == "" {
				diagnostics.AddError(
					"Bad reference", fmt.Sprintf("No object of kind %s with name %s", "DatasSource", ref),
				)
				return data, errors.New(fmt.Sprintf("Bad reference: No object of kind %s with name %s", "AlertNotificationTarget", ref))
			}
			sli.ThresholdMetric.MetricSource.MetricSourceRef = ref
		}
		if sli.RatioMetric.Bad.MetricSource.MetricSourceRef != "" {
			ref := sli.RatioMetric.Bad.MetricSource.MetricSourceRef
			sli.RatioMetric.Bad.MetricSource = data.Datasources[ref]
			if sli.RatioMetric.Bad.MetricSource.Metadata.Name == "" {
				diagnostics.AddError(
					"Bad reference", fmt.Sprintf("No object of kind %s with name %s", "DatasSource", ref),
				)
				return data, errors.New(fmt.Sprintf("Bad reference: No object of kind %s with name %s", "AlertNotificationTarget", ref))
			}
			sli.RatioMetric.Bad.MetricSource.MetricSourceRef = ref
		}
		if sli.RatioMetric.Good.MetricSource.MetricSourceRef != "" {
			ref := sli.RatioMetric.Good.MetricSource.MetricSourceRef
			sli.RatioMetric.Good.MetricSource = data.Datasources[ref]
			if sli.RatioMetric.Good.MetricSource.Metadata.Name == "" {
				diagnostics.AddError(
					"Bad reference", fmt.Sprintf("No object of kind %s with name %s", "DatasSource", ref),
				)
				return data, errors.New(fmt.Sprintf("Bad reference: No object of kind %s with name %s", "AlertNotificationTarget", ref))
			}
			sli.RatioMetric.Good.MetricSource.MetricSourceRef = ref
		}
		if sli.RatioMetric.Raw.MetricSource.MetricSourceRef != "" {
			ref := sli.RatioMetric.Raw.MetricSource.MetricSourceRef
			sli.RatioMetric.Raw.MetricSource = data.Datasources[ref]
			if sli.RatioMetric.Raw.MetricSource.Metadata.Name == "" {
				diagnostics.AddError(
					"Bad reference", fmt.Sprintf("No object of kind %s with name %s", "DatasSource", ref),
				)
				return data, errors.New(fmt.Sprintf("Bad reference: No object of kind %s with name %s", "AlertNotificationTarget", ref))
			}
			sli.RatioMetric.Raw.MetricSource.MetricSourceRef = ref
		}
		if sli.RatioMetric.Total.MetricSource.MetricSourceRef != "" {
			ref := sli.RatioMetric.Total.MetricSource.MetricSourceRef
			sli.RatioMetric.Total.MetricSource = data.Datasources[ref]
			if sli.RatioMetric.Total.MetricSource.Metadata.Name == "" {
				diagnostics.AddError(
					"Bad reference", fmt.Sprintf("No object of kind %s with name %s", "DatasSource", ref),
				)
				return data, errors.New(fmt.Sprintf("Bad reference: No object of kind %s with name %s", "AlertNotificationTarget", ref))
			}
			sli.RatioMetric.Total.MetricSource.MetricSourceRef = ref
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
			if objective.CompositeWeight == 0 {
				slo.Objectives[j].CompositeWeight = 1
			}
		}
		data.Slos[k] = slo
	}

	return data, nil
}
