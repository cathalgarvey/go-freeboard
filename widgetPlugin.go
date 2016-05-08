package freeboard

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

// WidgetPlugin is the interface expected of a prepared widget.
type WidgetPlugin interface {
	// Called when new settings are given.
	OnSettingsChanged(map[string]interface{})

	// A public function we must implement that will be called when
	// a calculated value changes. Since calculated values can change
	// at any time (like when a datasource is updated) we handle them
	// in a special callback function here.
	OnCalculatedValueChanged(settingName string, newValue interface{})

	// A public function we must implement that will be called when
	// freeboard wants us to render the contents of our widget.
	// The container element is the DIV that will surround the widget.
	// This will be automatically wrapped in jsutil.Wrap so no need
	// to worry about the js.Object->dom.HTMLElement conversion.
	Render(containerElement dom.HTMLElement)

	// How many 45-pixel blocks does this widget expect to be when
	// render is called?
	GetHeight() int

	// A public function we must implement that
	// will be called when this instance of this
	// plugin is no longer needed. Do anything
	// you need to cleanup after yourself here.
	OnDispose()
}

// WrapWidgetPlugin converts the exported function set of a DsPlugin
// into a map consisting of the required JS function set, with
// the other DsPlugin methods presented with lowercase leading
// characters in the JS style. Also present is "plugin", which
// directly references the DsPlugin object.
func WrapWidgetPlugin(wt WidgetPlugin) map[string]interface{} {
	return map[string]interface{}{
		"plugin":                   wt,
		"onSettingsChanged":        wt.OnSettingsChanged,
		"OnCalculatedValueChanged": wt.OnCalculatedValueChanged,
		"render":                   jsWrap(wt.Render),
		"getHeight":                wt.GetHeight,
		"onDispose":                wt.OnDispose,
	}
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

	// NewInstance is called to create a new plugin. This is the
	// Go wrapper. It is passed the calculated settings based on
	// the definition's settings array, as a js.Object; this should
	// be retained. It is also passed one special function, which
	// is a simple go wrapper around the updateCallback function
	// given by the FreeBoard NewInstance function.
	// Notably absent is NewInstanceCallback; this is handled under
	// the Go layer for you. All you have to do is return the
	// prepared DsPlugin-interfacing plugin object.
	NewInstance func(settings *js.Object) WidgetPlugin
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
	output["newInstance"] = func(settings, newInstanceCallback *js.Object) {
		Plugin := wtp.NewInstance(settings)
		wrapper := WrapWidgetPlugin(Plugin)
		newInstanceCallback.Invoke(wrapper)
	}
	return output
}
