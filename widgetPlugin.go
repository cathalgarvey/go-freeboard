package freeboard

import "honnef.co/go/js/dom"

// WidgetPlugin is the interface expected of a prepared widget.
type WidgetPlugin interface {
	// Called when new settings are given.
	onSettingsChanged(map[string]interface{})

	// A public function we must implement that will be called when
	// a calculated value changes. Since calculated values can change
	// at any time (like when a datasource is updated) we handle them
	// in a special callback function here.
	onCalculatedValueChanged(settingName string, newValue interface{})

	// A public function we must implement that will be called when
	// freeboard wants us to render the contents of our widget.
	// The container element is the DIV that will surround the widget.
	render(containerElement dom.BasicHTMLElement)

	// How many 45-pixel blocks does this widget expect to be when
	// render is called?
	getHeight() int

	// A public function we must implement that
	// will be called when this instance of this
	// plugin is no longer needed. Do anything
	// you need to cleanup after yourself here.
	onDispose()
}

// WtPluginDefinition is a Widget Plugin
type WtPluginDefinition struct {
	// TypeName should be a unique name for this plugin.
	// Must be a valid JS name. Avoid potential naming conflicts!
	TypeName string

	// The displayed name, need not be unique.
	DisplayName string

	// Front-facing description of this plugin.
	Description string

	// If this is set to true, the widget will fill
	// be allowed to fill the entire space given it,
	// otherwise it will contain an automatic padding
	// of around 10 pixels around it.
	FillSize bool

	// ExternalScripts are outside script URIs required for
	// this plugin. They will be loaded prior to the plugin.
	ExternalScripts []string

	// Settings are the user-facing options for this plugin.
	// They will be converted to a map[string]interface{} for
	// consumption by the plugin.
	Settings []FBSetting

	// NewInstance is called to create a new plugin. It is
	// passed the calculated settings based on the definition's
	// settings array. This should be kept by the plugin.
	// It is also passed one special function:
	// * newInstanceCallback should be called at the end of NewInstance
	//   with the new plugin.
	NewInstance func(settings map[string]interface{}, newInstanceCallback func(DsPlugin))
}

// ToFBInterface returns a map for FreeBoard's loadDatasourcePlugin func.
func (wtp WtPluginDefinition) ToFBInterface() map[string]interface{} {
	output := make(map[string]interface{})
	output["type_name"] = wtp.TypeName
	output["display_name"] = wtp.DisplayName
	output["description"] = wtp.Description
	output["fill_size"] = wtp.FillSize
	output["external_scripts"] = wtp.ExternalScripts
	output["settings"] = make([]map[string]interface{}, 0, len(wtp.Settings))
	settingSlice := output["settings"].([]map[string]interface{})
	for _, s := range wtp.Settings {
		settingSlice = append(settingSlice, s.ToFBInterface())
	}
	output["newInstance"] = wtp.NewInstance
	return output
}
