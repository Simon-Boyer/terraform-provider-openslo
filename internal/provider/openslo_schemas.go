package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var MetadataSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":         types.StringType,
		"display_name": types.StringType,
		"namespace":    types.StringType,
		"labels": types.MapType{
			ElemType: types.StringType,
		},
		"annotations": types.MapType{
			ElemType: types.StringType,
		},
	},
}

var DataSourceSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description": types.StringType,
		"connection_details": types.MapType{
			ElemType: types.StringType,
		},
		"metadata": MetadataSchema,
		"type":     types.StringType,
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

var MetricSourceSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		"spec": types.MapType{
			ElemType: types.StringType,
		},
		"datasource":        DataSourceSchema,
		"metric_source_ref": types.StringType,
	},
}

var MetricSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"metric_source": MetricSourceSchema,
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
