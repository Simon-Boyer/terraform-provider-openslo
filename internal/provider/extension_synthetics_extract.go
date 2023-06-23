package provider

import (
	"errors"
	"fmt"

	"github.com/goccy/go-yaml"
)

func (d *OpenSloDataSource) ExtractSyntheticsExtensionDocument(doc *YamlSpec, decType *yaml.Decoder) error {
	var err error
	switch doc.Kind {
	case "HTTPMonitor":
		var typedDoc YamlSpecTyped[HTTPMonitorModel]
		err = decType.Decode(&typedDoc)
		typedDoc.Spec.Metadata = doc.Metadata
		d.Extension_httpmonitor[doc.Metadata.Name] = typedDoc.Spec
	case "BrowserMonitor":
		var typedDoc YamlSpecTyped[BrowserMonitorModel]
		err = decType.Decode(&typedDoc)
		typedDoc.Spec.Metadata = doc.Metadata
		d.Extension_browsermonitor[doc.Metadata.Name] = typedDoc.Spec
	default:
		err = errors.New("Unknown kind: " + doc.Kind)
	}
	return err
}

func (d *OpenSloDataSource) SyntheticsExtensionPostExtractionLogic() error {
	for i := range d.Extension_httpmonitor {
		synthetic := d.Extension_httpmonitor[i]
		if synthetic.ServiceRef != "" {
			synthetic.Service = d.Services[d.Extension_httpmonitor[i].ServiceRef]
			if synthetic.Service.Metadata.Name == "" {
				return fmt.Errorf("bad reference: No object of kind %s with name %s", "synthetics_http", synthetic.ServiceRef)
			}
			d.Extension_httpmonitor[i] = synthetic
		}
	}
	for i := range d.Extension_browsermonitor {
		synthetic := d.Extension_browsermonitor[i]
		if synthetic.ServiceRef != "" {
			synthetic.Service = d.Services[d.Extension_browsermonitor[i].ServiceRef]
			if synthetic.Service.Metadata.Name == "" {
				return fmt.Errorf("bad reference: No object of kind %s with name %s", "synthetics_browser", synthetic.ServiceRef)
			}
			d.Extension_browsermonitor[i] = synthetic
		}
	}
	return nil
}
