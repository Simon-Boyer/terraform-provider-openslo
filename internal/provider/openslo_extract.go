package provider

import (
	"errors"
	"fmt"

	"github.com/goccy/go-yaml"
)

func (d *OpenSloDataSource) ExtractOpenSloDocument(doc *YamlSpec, decType *yaml.Decoder) error {
	var err error
	switch doc.Kind {
	case "DataSource":
		var typedDoc YamlSpecTyped[DataSourceModel]
		err = decType.Decode(&typedDoc)
		typedDoc.Spec.Metadata = doc.Metadata
		d.Datasources[doc.Metadata.Name] = typedDoc.Spec
	case "Service":
		var typedDoc YamlSpecTyped[ServiceModel]
		err = decType.Decode(&typedDoc)
		typedDoc.Spec.Metadata = doc.Metadata
		d.Services[doc.Metadata.Name] = typedDoc.Spec
	case "AlertCondition":
		var typedDoc YamlSpecTyped[AlertConditionModel]
		err = decType.Decode(&typedDoc)
		typedDoc.Spec.Metadata = doc.Metadata
		d.Alert_conditions[doc.Metadata.Name] = typedDoc.Spec
	case "AlertNotificationTarget":
		var typedDoc YamlSpecTyped[AlertNotificationTargetModel]
		err = decType.Decode(&typedDoc)
		typedDoc.Spec.Metadata = doc.Metadata
		d.Alert_notification_targets[doc.Metadata.Name] = typedDoc.Spec
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
		d.Alert_policies[doc.Metadata.Name] = typedDoc.Spec
	case "SLI":
		var typedDoc YamlSpecTyped[SLIModel]
		err = decType.Decode(&typedDoc)
		typedDoc.Spec.Metadata = doc.Metadata
		d.Slis[doc.Metadata.Name] = typedDoc.Spec
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
		d.Slos[doc.Metadata.Name] = typedDoc.Spec
	default:
		return errors.New("Unknown kind: " + doc.Kind)
	}
	return err
}

func (d *OpenSloDataSource) OpenSloPostExtractionLogic() error {
	// Embed referenced objects for alert policies
	for i := range d.Alert_policies {
		for j := range d.Alert_policies[i].Conditions {
			condition := d.Alert_policies[i].Conditions[j]
			if condition.ConditionRef != "" {
				linkedCond := d.Alert_conditions[condition.ConditionRef]
				if linkedCond.Metadata.Name == "" {
					return fmt.Errorf("bad reference: No object of kind %s with name %s", "AlertCondition", condition.ConditionRef)
				}
				linkedCond.ConditionRef = condition.ConditionRef
				d.Alert_policies[i].Conditions[j] = linkedCond
			}
		}
		for j := range d.Alert_policies[i].NotificationTargets {
			condition := d.Alert_policies[i].NotificationTargets[j]
			if condition.TargetRef != "" {
				linkedCond := d.Alert_notification_targets[condition.TargetRef]
				if linkedCond.Metadata.Name == "" {
					return fmt.Errorf("bad reference: No object of kind %s with name %s", "AlertNotificationTarget", condition.TargetRef)
				}
				linkedCond.TargetRef = condition.TargetRef
				d.Alert_policies[i].NotificationTargets[j] = linkedCond
			}
		}
	}

	// Embed referenced objects for slis
	for k := range d.Slis {
		sli := d.Slis[k]
		if sli.ThresholdMetric.MetricSource.MetricSourceRef != "" {
			ref := sli.ThresholdMetric.MetricSource.MetricSourceRef
			sli.ThresholdMetric.MetricSource.DataSource = d.Datasources[ref]
			if sli.ThresholdMetric.MetricSource.DataSource.Metadata.Name == "" {
				return fmt.Errorf("bad reference: No object of kind %s with name %s", "Datasources", ref)
			}
			sli.ThresholdMetric.MetricSource.MetricSourceRef = ref
			if sli.ThresholdMetric.MetricSource.DataSource.Type != "" {
				sli.ThresholdMetric.MetricSource.Type = sli.ThresholdMetric.MetricSource.DataSource.Type
			}
		}
		if sli.RatioMetric.Bad.MetricSource.MetricSourceRef != "" {
			ref := sli.RatioMetric.Bad.MetricSource.MetricSourceRef
			sli.RatioMetric.Bad.MetricSource.DataSource = d.Datasources[ref]
			if sli.RatioMetric.Bad.MetricSource.DataSource.Metadata.Name == "" {
				return fmt.Errorf("bad reference: No object of kind %s with name %s", "Datasources", ref)
			}
			sli.RatioMetric.Bad.MetricSource.MetricSourceRef = ref
			if sli.RatioMetric.Bad.MetricSource.DataSource.Type != "" {
				sli.RatioMetric.Bad.MetricSource.Type = sli.RatioMetric.Bad.MetricSource.DataSource.Type
			}
		}
		if sli.RatioMetric.Good.MetricSource.MetricSourceRef != "" {
			ref := sli.RatioMetric.Good.MetricSource.MetricSourceRef
			sli.RatioMetric.Good.MetricSource.DataSource = d.Datasources[ref]
			if sli.RatioMetric.Good.MetricSource.DataSource.Metadata.Name == "" {
				return fmt.Errorf("bad reference: No object of kind %s with name %s", "Datasources", ref)
			}
			sli.RatioMetric.Good.MetricSource.MetricSourceRef = ref
			if sli.RatioMetric.Good.MetricSource.DataSource.Type != "" {
				sli.RatioMetric.Good.MetricSource.Type = sli.RatioMetric.Good.MetricSource.DataSource.Type
			}
		}
		if sli.RatioMetric.Raw.MetricSource.MetricSourceRef != "" {
			ref := sli.RatioMetric.Raw.MetricSource.MetricSourceRef
			sli.RatioMetric.Raw.MetricSource.DataSource = d.Datasources[ref]
			if sli.RatioMetric.Raw.MetricSource.DataSource.Metadata.Name == "" {
				return fmt.Errorf("bad reference: No object of kind %s with name %s", "Datasources", ref)
			}
			sli.RatioMetric.Raw.MetricSource.MetricSourceRef = ref
			if sli.RatioMetric.Raw.MetricSource.DataSource.Type != "" {
				sli.RatioMetric.Raw.MetricSource.Type = sli.RatioMetric.Raw.MetricSource.DataSource.Type
			}
		}
		if sli.RatioMetric.Total.MetricSource.MetricSourceRef != "" {
			ref := sli.RatioMetric.Total.MetricSource.MetricSourceRef
			sli.RatioMetric.Total.MetricSource.DataSource = d.Datasources[ref]
			if sli.RatioMetric.Total.MetricSource.DataSource.Metadata.Name == "" {
				return fmt.Errorf("bad reference: No object of kind %s with name %s", "Datasources", ref)
			}
			sli.RatioMetric.Total.MetricSource.MetricSourceRef = ref
			if sli.RatioMetric.Total.MetricSource.DataSource.Type != "" {
				sli.RatioMetric.Total.MetricSource.Type = sli.RatioMetric.Total.MetricSource.DataSource.Type
			}
		}
		d.Slis[k] = sli
	}

	// Embed referenced objects for slos
	for k := range d.Slos {
		slo := d.Slos[k]
		if slo.IndicatorRef != "" {
			slo.Indicator = d.Slis[slo.IndicatorRef]
			if slo.Indicator.Metadata.Name == "" {
				return fmt.Errorf("bad reference: No object of kind %s with name %s", "Sli", slo.IndicatorRef)
			}
		}
		if slo.ServiceRef != "" {
			slo.Service = d.Services[slo.ServiceRef]
			if slo.Service.Metadata.Name == "" {
				return fmt.Errorf("bad reference: No object of kind %s with name %s", "Service", slo.ServiceRef)
			}
		}
		for j := range slo.AlertPolicies {
			alertPolicy := slo.AlertPolicies[j]
			if alertPolicy.AlertPolicyRef != "" {
				linkedAlertPolicy := d.Alert_policies[alertPolicy.AlertPolicyRef]
				if linkedAlertPolicy.Metadata.Name == "" {
					return fmt.Errorf("bad reference: No object of kind %s with name %s", "AlertPolicy", alertPolicy.AlertPolicyRef)
				}
				linkedAlertPolicy.AlertPolicyRef = alertPolicy.AlertPolicyRef
				slo.AlertPolicies[j] = linkedAlertPolicy
			}
		}
		for j := range slo.Objectives {
			objective := slo.Objectives[j]
			if objective.IndicatorRef != "" {
				objective.Indicator = d.Slis[objective.IndicatorRef]
				if objective.Indicator.Metadata.Name == "" {
					return fmt.Errorf("bad reference: No object of kind %s with name %s", "Sli", objective.IndicatorRef)
				}
				slo.Objectives[j] = objective
			}
			if objective.CompositeWeight == 0 {
				slo.Objectives[j].CompositeWeight = 1
			}
		}
		d.Slos[k] = slo
	}

	return nil
}
