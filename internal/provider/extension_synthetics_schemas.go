package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var HeaderSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	},
}

var RequestSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"headers": types.ListType{
			ElemType: HeaderSchema,
		},
		"body":   types.StringType,
		"method": types.StringType,
		"path":   types.StringType,
		"expected_response": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"codes":                types.NumberType,
				"payload_contains":     types.StringType,
				"payload_not_contains": types.StringType,
				"dt_postprocessing":    types.StringType,
			},
		},
	},
}

var HTTPMonitorSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"metadata": MetadataSchema,
		"url":      types.StringType,
		"requests": types.ListType{
			ElemType: RequestSchema,
		},
		"service":     ServiceSchema,
		"service_ref": types.StringType,
	},
}

var BrowserMonitorSchema = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"metadata":    MetadataSchema,
		"url":         types.StringType,
		"script":      types.StringType,
		"service":     ServiceSchema,
		"service_ref": types.StringType,
	},
}
