package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MetadataModel struct {
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
}

var MetadataSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":         types.StringType,
		"display_name": types.StringType,
	},
}

type DataSourceModel struct {
	Type              types.String  `tfsdk:"type"`
	ConnectionDetails types.Map     `tfsdk:"connection_details"`
	Metadata          MetadataModel `tfsdk:"metadata"`
}

var DataSourceSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":        types.StringType,
		"connection_details": types.MapType{},
		"metadata":           MetadataSchema,
	},
}

type ServiceModel struct {
	description types.String  `tfsdk:"description"`
	Metadata    MetadataModel `tfsdk:"metadata"`
}

var ServiceSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description": types.StringType,
		"metadata":    MetadataSchema,
	},
}

type AlertConditionModel struct {
	Description  types.String                      `tfsdk:"description"`
	Severity     types.String                      `tfsdk:"severity"`
	Condition    AlertConditionModelConditionModel `tfsdk:"condition"`
	Metadata     MetadataModel                     `tfsdk:"metadata"`
	ConditionRef types.String                      `tfsdk:"conditionRef"`
}

var AlertConditionSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description": types.StringType,
		"severity":    types.StringType,
		"condition":   AlertConditionModelConditionSchema,
		"metadata":    MetadataSchema,
	},
}

type AlertConditionModelConditionModel struct {
	Op             types.String `tfsdk:"op"`
	Threshold      types.Number `tfsdk:"threshold"`
	LookbackWindow types.String `tfsdk:"lookback_window"`
	AlertAfter     types.String `tfsdk:"alert_after"`
}

var AlertConditionModelConditionSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"op":              types.StringType,
		"threshold":       types.NumberType,
		"lookback_window": types.StringType,
		"alert_after":     types.StringType,
		"metadata":        MetadataSchema,
	},
}

type AlertNotificationTargetModel struct {
	TargetRef   types.String  `tfsdk:"targetRef"`
	Description types.String  `tfsdk:"description"`
	Target      types.String  `tfsdk:"target"`
	Metadata    MetadataModel `tfsdk:"metadata"`
}

var AlertNotificationTargetSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description": types.StringType,
		"target":      types.StringType,
		"metadata":    MetadataSchema,
	},
}

type AlertPolicyModel struct {
	AlertPolicyRef      types.String                   `tfsdk:"alertPolicyRef"`
	Description         types.String                   `tfsdk:"description"`
	AlertWhenNoData     types.Bool                     `tfsdk:"alert_when_no_data"`
	AlertWhenResolved   types.Bool                     `tfsdk:"alert_when_resolved"`
	AlertWhenBreaching  types.Bool                     `tfsdk:"alert_when_breaching"`
	Conditions          []AlertConditionModel          `tfsdk:"condition"`
	NotificationTargets []AlertNotificationTargetModel `tfsdk:"notification_targets"`
	Metadata            MetadataModel                  `tfsdk:"metadata"`
}

var AlertPolicySchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":          types.StringType,
		"alert_when_no_data":   types.BoolType,
		"alert_when_resolved":  types.BoolType,
		"alert_when_breaching": types.BoolType,
		"condition":            AlertConditionModelConditionSchema,
		"notification_targets": AlertNotificationTargetSchema,
		"metadata":             MetadataSchema,
	},
}

type SLIModel struct {
	Description     types.String     `tfsdk:"description"`
	ThresholdMetric MetricModel      `tfsdk:"threshold_metric"`
	RatioMetric     RatioMetricModel `tfsdk:"ratio_metric"`
	Metadata        MetadataModel    `tfsdk:"metadata"`
}

var SLISchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"threshold_metric": MetricSchema,
		"ratio_metric":     RatioMetricSchema,
		"metadata":         MetadataSchema,
	},
}

type MetricModel struct {
	MetricSource MetricSourceModel `tfsdk:"metric_source"`
}

var MetricSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"metric_source": MetricSourceSchema,
		"metadata":      MetadataSchema,
	},
}

type MetricSourceModel struct {
	MetricSourceRef types.String `tfsdk:"metric_source_ref"`
	Type            types.String `tfsdk:"type"`
	Spec            types.Object `tfsdk:"spec"`
}

var MetricSourceSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"metric_source_ref": types.StringType,
		"type":              types.StringType,
		"spec":              types.ObjectType{},
	},
}

type RatioMetricModel struct {
	counter types.Bool   `tfsdk:"counter"`
	good    MetricModel  `tfsdk:"good"`
	bad     MetricModel  `tfsdk:"bad"`
	total   MetricModel  `tfsdk:"total"`
	rawType types.String `tfsdk:"raw_type"`
	raw     MetricModel  `tfsdk:"raw"`
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

type SLOModel struct {
	Description     types.String       `tfsdk:"description"`
	Service         ServiceModel       `tfsdk:"service"`
	Indicator       SLIModel           `tfsdk:"indicator"`
	IndicatorRef    types.String       `tfsdk:"indicatorRef"`
	TimeWindow      TimeWindowModel    `tfsdk:"time_window"`
	BudgetingMethod types.String       `tfsdk:"budgeting_method"`
	Objectives      []ObjectiveModel   `tfsdk:"objectives"`
	AlertPolicies   []AlertPolicyModel `tfsdk:"alert_policies"`
	Metadata        MetadataModel      `tfsdk:"metadata"`
}

var SLOSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"service":          ServiceSchema,
		"indicator":        SLISchema,
		"time_window":      TimeWindowSchema,
		"budgeting_method": types.StringType,
		"objectives":       ObjectiveSchema,
		"alert_policies":   AlertPolicySchema,
		"metadata":         MetadataSchema,
	},
}

type TimeWindowModel struct {
	Duration  types.String  `tfsdk:"duration"`
	Calendar  CalendarModel `tfsdk:"calendar"`
	IsRolling types.Bool    `tfsdk:"is_rolling"`
}

var TimeWindowSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"duration":   types.StringType,
		"calendar":   CalendarSchema,
		"is_rolling": types.BoolType,
	},
}

type CalendarModel struct {
	StartTime types.String `tfsdk:"start_time"`
	TimeZone  types.String `tfsdk:"time_zone"`
}

var CalendarSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"start_time": types.StringType,
		"time_zone":  types.StringType,
	},
}

type ObjectiveModel struct {
	DisplayName      types.String `tfsdk:"display_name"`
	Op               types.String `tfsdk:"op"`
	value            types.Number `tfsdk:"value"`
	Target           types.Number `tfsdk:"target"`
	TargetPercentage types.Number `tfsdk:"target_percentage"`
	TimeSliceTarget  types.Number `tfsdk:"time_slice_target"`
	TimeSliceWindow  types.Number `tfsdk:"time_slice_window"`
	Indicator        SLIModel     `tfsdk:"indicator"`
	CompositeWeight  types.Number `tfsdk:"composite_weight"`
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
	},
}
