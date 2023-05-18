package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Metadata struct {
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
}

type DataSource struct {
	Type              types.String `tfsdk:"type"`
	ConnectionDetails types.Object `tfsdk:"connection_details"`
	Metadata
}

type ServiceModel struct {
	description types.String `tfsdk:"description"`
	Metadata
}

type AlertConditionModel struct {
	Description types.String                 `tfsdk:"description"`
	Severity    types.String                 `tfsdk:"severity"`
	Condition   AlertConditionModelCondition `tfsdk:"condition"`
	Metadata
}

type AlertConditionModelCondition struct {
	Op             types.String `tfsdk:"op"`
	Threshold      types.Number `tfsdk:"threshold"`
	LookbackWindow types.String `tfsdk:"lookback_window"`
	AlertAfter     types.String `tfsdk:"alert_after"`
}

type AlertNotificationTargetModel struct {
	Description types.String `tfsdk:"description"`
	Target      types.String `tfsdk:"target"`
	Metadata
}

type AlertPolicy struct {
	Description         types.String                   `tfsdk:"description"`
	AlertWhenNoData     types.Bool                     `tfsdk:"alert_when_no_data"`
	AlertWhenResolved   types.Bool                     `tfsdk:"alert_when_resolved"`
	AlertWhenBreaching  types.Bool                     `tfsdk:"alert_when_breaching"`
	Condition           AlertConditionModel            `tfsdk:"condition"`
	NotificationTargets []AlertNotificationTargetModel `tfsdk:"notification_targets"`
	Metadata
}

type SLI struct {
	Description     types.String `tfsdk:"description"`
	ThresholdMetric Metric       `tfsdk:"threshold_metric"`
	RatioMetric     RatioMetric  `tfsdk:"ratio_metric"`
	Metadata
}

type Metric struct {
	MetricSource MetricSource `tfsdk:"metric_source"`
}

type MetricSource struct {
	MetricSourceRef types.String `tfsdk:"metric_source_ref"`
	Type            types.String `tfsdk:"type"`
	Spec            types.Object `tfsdk:"spec"`
}

type RatioMetric struct {
	counter types.Bool   `tfsdk:"counter"`
	good    Metric       `tfsdk:"good"`
	bad     Metric       `tfsdk:"bad"`
	total   Metric       `tfsdk:"total"`
	rawType types.String `tfsdk:"raw_type"`
	raw     Metric       `tfsdk:"raw"`
}

type SLO struct {
	Description     types.String  `tfsdk:"description"`
	Service         ServiceModel  `tfsdk:"service"`
	Indicator       SLI           `tfsdk:"indicator"`
	TimeWindow      TimeWindow    `tfsdk:"time_window"`
	BudgetingMethod types.String  `tfsdk:"budgeting_method"`
	Objectives      []Objective   `tfsdk:"objectives"`
	AlertPolicies   []AlertPolicy `tfsdk:"alert_policies"`
	Metadata
}

type TimeWindow struct {
	Duration  types.String `tfsdk:"duration"`
	Calendar  Calendar     `tfsdk:"calendar"`
	IsRolling types.Bool   `tfsdk:"is_rolling"`
}

type Calendar struct {
	StartTime types.String `tfsdk:"start_time"`
	TimeZone  types.String `tfsdk:"time_zone"`
}

type Objective struct {
	DisplayName      types.String `tfsdk:"display_name"`
	Op               types.String `tfsdk:"op"`
	value            types.Number `tfsdk:"value"`
	Target           types.Number `tfsdk:"target"`
	TargetPercentage types.Number `tfsdk:"target_percentage"`
	TimeSliceTarget  types.Number `tfsdk:"time_slice_target"`
	TimeSliceWindow  types.Number `tfsdk:"time_slice_window"`
	Indicator        SLI          `tfsdk:"indicator"`
	CompositeWeight  types.Number `tfsdk:"composite_weight"`
}
