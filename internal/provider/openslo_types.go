package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MetadataModel struct {
	Name        string `tfsdk:"name" yaml:"name"`
	DisplayName string `tfsdk:"display_name" yaml:"displayName"`
}

type DataSourceModel struct {
	Type              string        `tfsdk:"type" yaml:"type"`
	ConnectionDetails types.Map     `tfsdk:"connection_details" yaml:"connectionDetails"`
	Metadata          MetadataModel `tfsdk:"metadata" yaml:"metadata"`
}

type ServiceModel struct {
	Description string        `tfsdk:"description" yaml:"description"`
	Metadata    MetadataModel `tfsdk:"metadata" yaml:"metadata"`
}

type AlertConditionModel struct {
	Description  string                            `tfsdk:"description" yaml:"description"`
	Severity     string                            `tfsdk:"severity" yaml:"severity"`
	Condition    AlertConditionModelConditionModel `tfsdk:"condition" yaml:"condition"`
	Metadata     MetadataModel                     `tfsdk:"metadata" yaml:"metadata"`
	ConditionRef string                            `tfsdk:"condition_ref" yaml:"conditionRef"`
}

type AlertConditionModelConditionModel struct {
	Op             string  `tfsdk:"op" yaml:"op"`
	Threshold      float64 `tfsdk:"threshold" yaml:"threshold"`
	LookbackWindow string  `tfsdk:"lookback_window" yaml:"lookbackWindow"`
	AlertAfter     string  `tfsdk:"alert_after" yaml:"alertAfter"`
}

type AlertNotificationTargetModel struct {
	TargetRef   string        `tfsdk:"target_ref" yaml:"targetRef"`
	Description string        `tfsdk:"description" yaml:"description"`
	Target      string        `tfsdk:"target" yaml:"target"`
	Metadata    MetadataModel `tfsdk:"metadata" yaml:"metadata"`
}

type AlertPolicyModel struct {
	AlertPolicyRef      string                         `tfsdk:"alert_policy_ref" yaml:"alertPolicyRef"`
	Description         string                         `tfsdk:"description" yaml:"description"`
	AlertWhenNoData     bool                           `tfsdk:"alert_when_no_data" yaml:"alertWhenNoData"`
	AlertWhenResolved   bool                           `tfsdk:"alert_when_resolved" yaml:"alertWhenResolved"`
	AlertWhenBreaching  bool                           `tfsdk:"alert_when_breaching" yaml:"alertWhenBreaching"`
	Conditions          []AlertConditionModel          `tfsdk:"condition" yaml:"condition"`
	NotificationTargets []AlertNotificationTargetModel `tfsdk:"notification_targets" yaml:"notificationTargets"`
	Metadata            MetadataModel                  `tfsdk:"metadata" yaml:"metadata"`
}

type SLIModel struct {
	Description     string           `tfsdk:"description" yaml:"description"`
	ThresholdMetric MetricModel      `tfsdk:"threshold_metric" yaml:"thresholdMetric"`
	RatioMetric     RatioMetricModel `tfsdk:"ratio_metric" yaml:"ratioMetric"`
	Metadata        MetadataModel    `tfsdk:"metadata" yaml:"metadata"`
}

type MetricModel struct {
	MetricSource MetricSourceModel `tfsdk:"metric_source" yaml:"metricSource"`
}

type MetricSourceModel struct {
	MetricSourceRef string       `tfsdk:"metric_source_ref" yaml:"metricSourceRef"`
	Type            string       `tfsdk:"type" yaml:"type"`
	Spec            types.Object `tfsdk:"spec" yaml:"spec"`
}

type RatioMetricModel struct {
	counter bool        `tfsdk:"counter" yaml:"counter"`
	good    MetricModel `tfsdk:"good" yaml:"good"`
	bad     MetricModel `tfsdk:"bad" yaml:"bad"`
	total   MetricModel `tfsdk:"total" yaml:"total"`
	rawType string      `tfsdk:"raw_type" yaml:"rawType"`
	raw     MetricModel `tfsdk:"raw" yaml:"raw"`
}

type SLOModel struct {
	Description     string             `tfsdk:"description" yaml:"description"`
	Service         ServiceModel       `tfsdk:"service" yaml:"service"`
	Indicator       SLIModel           `tfsdk:"indicator" yaml:"indicator"`
	IndicatorRef    string             `tfsdk:"indicator_ref" yaml:"indicatorRef"`
	TimeWindow      TimeWindowModel    `tfsdk:"time_window" yaml:"timeWindow"`
	BudgetingMethod string             `tfsdk:"budgeting_method" yaml:"budgetingMethod"`
	Objectives      []ObjectiveModel   `tfsdk:"objectives" yaml:"objectives"`
	AlertPolicies   []AlertPolicyModel `tfsdk:"alert_policies" yaml:"alertPolicies"`
	Metadata        MetadataModel      `tfsdk:"metadata" yaml:"metadata"`
}

type TimeWindowModel struct {
	Duration  string        `tfsdk:"duration" yaml:"duration"`
	Calendar  CalendarModel `tfsdk:"calendar" yaml:"calendar"`
	IsRolling bool          `tfsdk:"is_rolling" yaml:"isRolling"`
}

type CalendarModel struct {
	StartTime string `tfsdk:"start_time" yaml:"startTime"`
	TimeZone  string `tfsdk:"time_zone" yaml:"timeZone"`
}

type ObjectiveModel struct {
	DisplayName      string   `tfsdk:"display_name" yaml:"displayName"`
	Op               string   `tfsdk:"op" yaml:"op"`
	value            float64  `tfsdk:"value" yaml:"value"`
	Target           float64  `tfsdk:"target" yaml:"target"`
	TargetPercentage float64  `tfsdk:"target_percentage" yaml:"targetPercentage"`
	TimeSliceTarget  float64  `tfsdk:"time_slice_target" yaml:"timeSliceTarget"`
	TimeSliceWindow  float64  `tfsdk:"time_slice_window" yaml:"timeSliceWindow"`
	IndicatorRef     string   `tfsdk:"indicator_ref" yaml:"indicatorRef"`
	Indicator        SLIModel `tfsdk:"indicator" yaml:"indicator"`
	CompositeWeight  float64  `tfsdk:"composite_weight" yaml:"compositeWeight"`
}
