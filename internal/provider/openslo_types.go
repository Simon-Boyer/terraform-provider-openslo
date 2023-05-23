package provider

type MetadataModel struct {
	Name        string `tfsdk:"name" yaml:"name"`
	DisplayName string `tfsdk:"display_name" yaml:"displayName"`
}

type DataSourceModel struct {
	Type              string            `tfsdk:"type" yaml:"type"`
	ConnectionDetails map[string]string `tfsdk:"connection_details" yaml:"connectionDetails"`
	Metadata          MetadataModel     `tfsdk:"metadata" yaml:"metadata"`
	MetricSourceRef   string            `tfsdk:"metric_source_ref" yaml:"metricSourceRef"`
	Spec              map[string]string `tfsdk:"spec" yaml:"spec"`
	Description       string            `tfsdk:"description" yaml:"description"`
}

type ServiceModel struct {
	Description string        `tfsdk:"description" yaml:"description"`
	Metadata    MetadataModel `tfsdk:"metadata" yaml:"metadata"`
}

type AlertConditionModel struct {
	Description  string                            `tfsdk:"description" yaml:"description"`
	Severity     string                            `tfsdk:"severity" yaml:"severity"`
	Condition    AlertConditionModelConditionModel `tfsdk:"condition" yaml:"condition,omitempty"`
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
	Conditions          []AlertConditionModel          `tfsdk:"conditions" yaml:"conditions"`
	NotificationTargets []AlertNotificationTargetModel `tfsdk:"notification_targets" yaml:"notificationTargets"`
	Metadata            MetadataModel                  `tfsdk:"metadata" yaml:"metadata"`
}

type SLIModel struct {
	Description     string           `tfsdk:"description" yaml:"description"`
	ThresholdMetric MetricModel      `tfsdk:"threshold_metric" yaml:"thresholdMetric,omitempty"`
	RatioMetric     RatioMetricModel `tfsdk:"ratio_metric" yaml:"ratioMetric,omitempty"`
	Metadata        MetadataModel    `tfsdk:"metadata" yaml:"metadata"`
}

type MetricModel struct {
	MetricSource DataSourceModel `tfsdk:"metric_source" yaml:"metricSource,omitempty"`
}

type RatioMetricModel struct {
	Counter bool        `tfsdk:"counter" yaml:"counter"`
	Good    MetricModel `tfsdk:"good" yaml:"good,omitempty"`
	Bad     MetricModel `tfsdk:"bad" yaml:"bad,omitempty"`
	Total   MetricModel `tfsdk:"total" yaml:"total,omitempty"`
	RawType string      `tfsdk:"raw_type" yaml:"rawType"`
	Raw     MetricModel `tfsdk:"raw" yaml:"raw,omitempty"`
}

type SLOModel struct {
	Description     string             `tfsdk:"description" yaml:"description"`
	Service         ServiceModel       `tfsdk:"service" yaml:"-"`
	ServiceRef      string             `tfsdk:"service_ref" yaml:"service"`
	Indicator       SLIModel           `tfsdk:"indicator" yaml:"indicator,omitempty"`
	IndicatorRef    string             `tfsdk:"indicator_ref" yaml:"indicatorRef"`
	TimeWindow      []TimeWindowModel  `tfsdk:"time_window" yaml:"timeWindow"`
	BudgetingMethod string             `tfsdk:"budgeting_method" yaml:"budgetingMethod"`
	Objectives      []ObjectiveModel   `tfsdk:"objectives" yaml:"objectives"`
	AlertPolicies   []AlertPolicyModel `tfsdk:"alert_policies" yaml:"alertPolicies"`
	Metadata        MetadataModel      `tfsdk:"metadata" yaml:"metadata"`
}

type TimeWindowModel struct {
	Duration  string        `tfsdk:"duration" yaml:"duration"`
	Calendar  CalendarModel `tfsdk:"calendar" yaml:"calendar,omitempty"`
	IsRolling bool          `tfsdk:"is_rolling" yaml:"isRolling"`
}

type CalendarModel struct {
	StartTime string `tfsdk:"start_time" yaml:"startTime"`
	TimeZone  string `tfsdk:"time_zone" yaml:"timeZone"`
}

type ObjectiveModel struct {
	DisplayName      string   `tfsdk:"display_name" yaml:"displayName"`
	Op               string   `tfsdk:"op" yaml:"op"`
	Value            float64  `tfsdk:"value" yaml:"value"`
	Target           float64  `tfsdk:"target" yaml:"target"`
	TargetPercentage float64  `tfsdk:"target_percentage" yaml:"targetPercentage"`
	TimeSliceTarget  float64  `tfsdk:"time_slice_target" yaml:"timeSliceTarget"`
	TimeSliceWindow  float64  `tfsdk:"time_slice_window" yaml:"timeSliceWindow"`
	IndicatorRef     string   `tfsdk:"indicator_ref" yaml:"indicatorRef"`
	Indicator        SLIModel `tfsdk:"indicator" yaml:"indicator,omitempty"`
	CompositeWeight  float64  `tfsdk:"composite_weight" yaml:"compositeWeight"`
}
